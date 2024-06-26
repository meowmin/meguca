package feeds

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"

	"github.com/bakape/meguca/common"
	"github.com/go-playground/log"
)

type message struct {
	id  uint64
	msg []byte
}

type postCreationMessage struct {
	post common.Post
	message
}

type imageInsertionMessage struct {
	spoilered bool
	message
}

type postBodyModMessage struct {
	id   uint64
	body string
}

type moderationMessage struct {
	message
	entry common.ModerationEntry
}

type syncCount struct {
	Active int `json:"active"`
	Total  int `json:"total"`
}

type pendingTiktokState struct {
	message message
	state   PendingTikToks
}

// Feed is a feed with synchronization logic of a certain thread
type Feed struct {
	// Thread ID
	id uint64
	// Message flushing ticker
	ticker
	// Common functionality
	baseFeed
	// Buffer of unsent messages
	messageBuffer
	// Entire thread cached into memory
	cache threadCache
	// Propagates mesages to all listeners
	send chan []byte
	// Insert a new post into the thread and propagate to listeners
	insertPost chan postCreationMessage
	// Insert an image into an already allocated post
	insertImage chan imageInsertionMessage
	// Send message to close a post along with parsed post data
	closePost chan message
	// Send message to spoiler image of a specific post
	spoilerImage chan message
	// Set body of an open post
	setOpenBody chan postBodyModMessage
	// Send message about post moderation
	moderatePost chan moderationMessage
	// Let sent sync counter
	lastSyncCount syncCount
	// Let sent sync counter
	binaryMessages chan []byte
	// Tiktok state
	updatePendingTiktokState chan pendingTiktokState
}

// Start read existing posts into cache and start main loop
func (f *Feed) Start() (err error) {
	f.cache, err = newThreadCache(f.id)
	if err != nil {
		return
	}

	go func() {
		// Stop the timer, if there are no messages and resume on new ones.
		// Keeping the goroutine asleep reduces CPU usage.
		f.start()
		defer f.pause()

		evictionTimer := time.NewTicker(time.Minute)
		defer evictionTimer.Stop()

		for {
			select {

			case <-evictionTimer.C:
				f.cache.evict()

			// Add client
			case c := <-f.add:
				f.addClient(c)

				msg, err := f.cache.getSyncMessage()
				if err != nil {
					log.Errorf("sync message: %s", err)
				}
				c.Send(msg)

				f.sendIPCount()

			// Remove client and close feed, if no clients left
			case c := <-f.remove:
				if f.removeClient(c) {
					return
				}

				f.sendIPCount()

			// Buffer external message and prepare for sending to all clients
			case msg := <-f.send:
				f.bufferMessage(msg)

			// Send any buffered messages to any listening clients
			case <-f.C:
				if buf := f.flush(); buf == nil {
					f.pause()
				} else {
					f.sendToAll(buf)
				}

			// Insert a new post, cache and propagate
			case msg := <-f.insertPost:
				src := msg.post
				f.modifyPostImmediate(msg.message, func(p *cachedPost) {
					*p = cachedPost{
						HasImage:  src.Image != nil,
						Spoilered: src.Image != nil && src.Image.Spoiler,
						Time:      src.Time,
						Closed:    !src.Editing,
						Body:      src.Body,
					}
				})
				// Post can be automatically deleted on insertion
				if src.IsDeleted() {
					f.cache.Moderation[msg.id] = src.Moderation
				}
				f.sendIPCount()

			// Set the body of an open post and propagate
			case msg := <-f.setOpenBody:
				f.updateCachedPost(msg.id, msg.body)

			case msg := <-f.insertImage:
				f.modifyPostImmediate(msg.message, func(p *cachedPost) {
					p.HasImage = true
					p.Spoilered = msg.spoilered
				})

			case msg := <-f.spoilerImage:
				f.modifyPostImmediate(msg, func(p *cachedPost) {
					p.Spoilered = true
				})

			case msg := <-f.closePost:
				f.modifyPostImmediate(msg, func(p *cachedPost) {
					p.Closed = true
				})
			case msg := <-f.updatePendingTiktokState:
				f.modifyPost(msg.message, func(p *cachedPost) {
					p.PendingTikToks = msg.state
				})
			case msg := <-f.binaryMessages:
				f.sendToAllBinary(msg)

			// Posts being moderated
			case msg := <-f.moderatePost:
				f.modifyPost(msg.message, func(p *cachedPost) {
					switch msg.entry.Type {
					case common.PurgePost:
						p.Body = ""
						fallthrough
					case common.DeleteImage:
						p.HasImage = false
						p.Spoilered = false
					case common.SpoilerImage:
						p.Spoilered = true
					}
				})
				f.cache.Moderation[msg.id] = append(f.cache.Moderation[msg.id],
					msg.entry)
			}
		}
	}()

	return
}

func (f *Feed) modifyPost(msg message, fn func(*cachedPost)) {
	f.startIfPaused()

	p := f.cache.Recent[msg.id]
	fn(&p)
	f.cache.Recent[msg.id] = p

	if msg.msg != nil {
		f.write(msg.msg)
	}
	f.cache.clearMemoized()
}

func (f *Feed) modifyPostImmediate(msg message, fn func(*cachedPost)) {
	f.startIfPaused()

	p := f.cache.Recent[msg.id]
	fn(&p)
	f.cache.Recent[msg.id] = p

	if msg.msg != nil {
		f.sendToAll(msg.msg)
	}
	f.cache.clearMemoized()
}

// Simply updates the body of a cached post
func (f *Feed) updateCachedPost(id uint64, newBody string) {
	f.startIfPaused()
	p := f.cache.Recent[id]
	p.Body = newBody
	f.cache.Recent[id] = p
	f.cache.clearMemoized()
}

// Send a message to all listening clients
func (f *Feed) Send(msg []byte) {
	f.send <- msg
}

// Buffer a message to be sent on the next tick
func (f *Feed) bufferMessage(msg []byte) {
	f.startIfPaused()
	f.write(msg)
}

// Send unique IP count to all connected clients
func (f *Feed) sendIPCount() {
	var active int
	ips := make(map[string]struct{}, len(f.clients))
	pastHour := time.Now().Add(-time.Hour).Unix()

	for c := range f.clients {
		ip := c.IP()
		if _, ok := ips[ip]; !ok && c.LastTime() >= pastHour {
			active++
		}
		ips[ip] = struct{}{}
	}

	new := syncCount{
		Active: active,
		Total:  len(ips),
	}
	if new != f.lastSyncCount {
		f.lastSyncCount = new
		msg, _ := common.EncodeMessage(common.MessageSyncCount, new)
		f.bufferMessage(msg)
	}
}

// InsertPost inserts a new post into the thread or reclaim an open post after disconnect
// and propagate to listeners
func (f *Feed) InsertPost(p common.Post, msg []byte) {
	f.insertPost <- postCreationMessage{
		message: message{
			id:  p.ID,
			msg: msg,
		},
		post: p,
	}
}

// InsertImage inserts an image into an already allocated post
func (f *Feed) InsertImage(id uint64, spoilered bool, msg []byte) {
	f.insertImage <- imageInsertionMessage{
		message: message{
			id:  id,
			msg: msg,
		},
		spoilered: spoilered,
	}
}

// ClosePost closes a feed's post
func (f *Feed) ClosePost(id uint64, msg []byte) {
	f.closePost <- message{
		id:  id,
		msg: msg,
	}
}

// SpoilerImage spoilers a feed's image
func (f *Feed) SpoilerImage(id uint64, msg []byte) {
	f.spoilerImage <- message{id, msg}
}

func (f *Feed) _moderatePost(id uint64, msg []byte,
	entry common.ModerationEntry,
) {
	f.moderatePost <- moderationMessage{
		message: message{
			id:  id,
			msg: msg,
		},
		entry: entry,
	}
}

// UpdateBody sets the body of an open post and send update message to clients
func (f *Feed) UpdateBody(id uint64, body string, msg []byte) {
	f.binaryMessages <- msg
	f.setOpenBody <- postBodyModMessage{
		id:   id,
		body: body,
	}
}

func (f *Feed) SendClaudeToken(id uint64, token string) {
	messageSize := 9 + len(token)
	message := make([]byte, messageSize)
	binary.LittleEndian.PutUint64(message, math.Float64bits(float64(id)))
	copy(message[8:], token)
	message[messageSize-1] = uint8(common.MessageClaudeAppend)
	f.binaryMessages <- message
}
func (f *Feed) SendClaudeComplete(id uint64, isError bool, response *bytes.Buffer) {
	messageSize := 9 + response.Len()
	message := make([]byte, messageSize)
	binary.LittleEndian.PutUint64(message, math.Float64bits(float64(id)))
	copy(message[8:], response.Bytes())
	if isError {
		message[messageSize-1] = uint8(common.MessageClaudeError)
	} else {
		message[messageSize-1] = uint8(common.MessageClaudeDone)
	}
	f.binaryMessages <- message
}
func (f *Feed) GetPendingTiktokState(id uint64) (p PendingTikToks, ok bool) {
	//return f.cache.Recent[id].PendingTikToks
	var post cachedPost
	post, ok = f.cache.Recent[id]
	if ok {
		p = post.PendingTikToks
	}
	return
}

func (f *Feed) UpdatePendingTiktokState(id uint64, state PendingTikToks) {
	stateUpdate := struct {
		ID    uint64         `json:"id"`
		State PendingTikToks `json:"state"`
	}{
		id,
		state,
	}
	msg, _ := common.EncodeMessage(common.MessageTiktokState, stateUpdate)
	f.updatePendingTiktokState <- pendingTiktokState{
		message: message{id, msg},
		state:   state,
	}
}

func (f *Feed) SendBinaryMessage(msg []byte) {
	f.binaryMessages <- msg
}
