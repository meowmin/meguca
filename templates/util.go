package templates

import (
	"github.com/bakape/meguca/common"
	"github.com/xeonx/timeago"
	"html"
	"strconv"
	"time"
)

// CalculateOmit returns the omitted post and image counts for a thread
func CalculateOmit(t common.Thread) (int, int) {
	// There might still be posts missing due to deletions even in complete
	// thread queries. Ensure we are actually retrieving an abbreviated thread
	// before calculating.
	if !t.Abbrev {
		return 0, 0
	}

	var (
		omit    = int(t.PostCount) - (len(t.Posts) + 1)
		imgOmit uint32
	)
	if omit != 0 {
		imgOmit = t.ImageCount
		if t.Image != nil {
			imgOmit--
		}
		for _, p := range t.Posts {
			if p.Image != nil {
				imgOmit--
			}
		}
	}
	return omit, int(imgOmit)
}

func bold(s string) string {
	s = html.EscapeString(s)
	b := make([]byte, 3, len(s)+7)
	copy(b, "<b>")
	b = append(b, s...)
	b = append(b, "</b>"...)
	return string(b)
}

func getTokID(filename string) *string {
	digits := ""
	for _, c := range filename {
		if c >= '0' && c <= '9' {
			digits += string(c)
		} else {
			if len(digits) == 19 {
				return &digits
			}
			digits = ""
		}
	}
	if len(digits) == 19 {
		return &digits
	}
	return nil
}

const scale int64 = 4294967296

func relativeTime(tokID string) string {
	then, err := strconv.ParseInt(tokID, 10, 64)
	if err != nil {
		// handle error
	}
	then = then / scale
	thenTime := time.Unix(then, 0)
	config := timeago.English
	config.Max = 1<<63 - 1
	config.Periods = config.Periods[:len(config.Periods)-1]
	return config.Format(thenTime)
}
