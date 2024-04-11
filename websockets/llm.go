package websockets

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/bakape/meguca/common"
	"github.com/bakape/meguca/config"
	"io"
	"net/http"
	"strings"
)

const (
	Claude3Opus         = "claude-3-opus-20240229"
	Claude3Sonnet       = "claude-3-sonnet-20240229"
	Claude3Haiku        = "claude-3-haiku-20240307"
	DefaultSystemPrompt = `You are an AI assistant designed to provide concise and helpful responses to user questions on an online chatroom. 
You will assist users by answering their queries directly and succinctly. This is a fast-paced platform so you must keep your responses extremely concise.`
)

func encodeMessages(prompt string, img *[]byte) []byte {
	buf := bytes.Buffer{}
	buf.WriteRune('[')
	if img != nil {
		buf.Write(*img)
		buf.WriteRune(',')
	}
	buf.WriteString(`{"type":"text","text":`)
	promptStr, _ := json.Marshal(prompt)
	buf.Write(promptStr)
	buf.WriteString(`}`)
	buf.WriteRune(']')
	return buf.Bytes()
}

func StreamMessages(model string, systemPrompt string, maxTokens int, claudeState *common.ClaudeState, img *[]byte, start func(), token func(string), done func()) error {
	apiKey := config.Server.AnthropicApiKey

	url := "https://api.anthropic.com/v1/messages"

	body := requestData{
		model,
		[]messageParam{
			{
				"user", encodeMessages(claudeState.Prompt, img),
			},
		},
		maxTokens,
		true,
		systemPrompt,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "messages-2023-12-15")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		line, err := reader.ReadBytes('\n')
		if strings.HasPrefix(string(line), `{"type":"error"`) {
			var errData errorResponse
			err = json.Unmarshal(line, &errData)
			if err != nil {
				return err
			}
			claudeState.Status = common.Error
			claudeState.Response.Reset()
			claudeState.Response.WriteString(errData.Error.Message)
			done()
			return nil
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := bytes.SplitN(line, []byte(": "), 2)
		if len(parts) != 2 {
			continue
		}

		eventType := parts[0]
		if bytes.Equal(eventType, []byte("event")) {
			eventVal := parts[1]
			if bytes.Equal(eventVal, []byte("content_block_start")) {
				claudeState.Status = common.Generating
				start()
			} else if bytes.Equal(eventVal, []byte("content_block_stop")) {
				claudeState.Status = common.Done
				done()
			}
			continue
		}
		if bytes.Equal(eventType, []byte("data")) {
			var event contentBlockDeltaEvent
			err = json.Unmarshal(parts[1], &event)
			if event.Type == "content_block_delta" {
				claudeState.Response.WriteString(event.Delta.Text)
				token(event.Delta.Text)
			}
			if err != nil {
				return err
			}
		}

	}

	return nil
}

type errorResponse struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}
type requestData struct {
	Model        string         `json:"model"`
	Messages     []messageParam `json:"messages"`
	MaxTokens    int            `json:"max_tokens"`
	Stream       bool           `json:"stream"`
	SystemPrompt string         `json:"system,omitempty"`
}

type contentBlock struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type contentBlockDeltaEvent struct {
	Delta textDelta `json:"delta"`
	Index int       `json:"index"`
	Type  string    `json:"type"`
}

type contentBlockStartEvent struct {
	ContentBlock contentBlock `json:"content_block"`
	Index        int          `json:"index"`
	Type         string       `json:"type"`
}

type contentBlockStopEvent struct {
	Index int    `json:"index"`
	Type  string `json:"type"`
}

type imageBlockParam struct {
	Source imageBlockParamSource `json:"source"`
	Type   string                `json:"type,omitempty"`
}

type imageBlockParamSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type,omitempty"`
}

type message struct {
	ID           string         `json:"id"`
	Content      []contentBlock `json:"content"`
	Model        string         `json:"model"`
	Role         string         `json:"role"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Type         string         `json:"type"`
	Usage        usage          `json:"usage"`
}

type messageDeltaEvent struct {
	Delta messageDeltaEventDelta `json:"delta"`
	Type  string                 `json:"type"`
	Usage messageDeltaUsage      `json:"usage"`
}

type messageDeltaEventDelta struct {
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
}

type messageDeltaUsage struct {
	OutputTokens int `json:"output_tokens"`
}

type messageParam struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type messageStartEvent struct {
	Message message `json:"message"`
	Type    string  `json:"type"`
}

type messageStopEvent struct {
	Type string `json:"type"`
}

type textBlock struct {
	Text string `json:"text"`
	Type string `json:"type,omitempty"`
}

type textDelta struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
