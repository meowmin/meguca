// Package feeds manages client synchronization to update feeds and provides a
// thread-safe interface for propagating messages to them and reassigning feeds
// to and from clients.
package feeds

import (
	"errors"
	"github.com/bakape/meguca/common"
	"github.com/bakape/meguca/db"
	"github.com/bakape/meguca/pb"
	"sync"
)

// Contains and manages all active update feeds
var feeds = feedMap{
	// 64 len map to avoid some possible reallocation as the server starts
	feeds:       make(map[uint64]*Feed, 64),
	tvFeeds:     make(map[string]*tvFeed, 64),
	nekotvFeeds: make(map[uint64]*NekoTVFeed, 64),
}

// Export to avoid circular dependency
func init() {
	common.SendTo = SendTo
	common.ClosePost = ClosePost
}

// Container for managing client<->update-feed assignment and interaction
type feedMap struct {
	feeds       map[uint64]*Feed
	tvFeeds     map[string]*tvFeed
	mu          sync.RWMutex
	nekotvFeeds map[uint64]*NekoTVFeed
}

// Add client to feed and send it the current status of the feed for
// synchronization to the feed's internal state
func addToFeed(id uint64, board string, c common.Client) (
	feed *Feed, err error,
) {
	feeds.mu.Lock()
	defer feeds.mu.Unlock()

	var ok bool

	if id != 0 {
		feed, ok = feeds.feeds[id]
		if !ok {
			feed = &Feed{
				id:                       id,
				send:                     make(chan []byte),
				insertPost:               make(chan postCreationMessage),
				closePost:                make(chan message),
				spoilerImage:             make(chan message),
				moderatePost:             make(chan moderationMessage),
				setOpenBody:              make(chan postBodyModMessage),
				insertImage:              make(chan imageInsertionMessage),
				messageBuffer:            make([]string, 0, 64),
				binaryMessages:           make(chan []byte),
				updatePendingTiktokState: make(chan pendingTiktokState),
			}

			feed.baseFeed.init()
			feeds.feeds[id] = feed
			err = feed.Start()
			if err != nil {
				return
			}
		}
		feed.add <- c
	}

	return
}

// SubscribeToMeguTV subscribes to random video stream.
// Clients are automatically unsubscribed, when leaving their current sync feed.
func SubscribeToMeguTV(c common.Client) (err error) {
	feeds.mu.Lock()
	defer feeds.mu.Unlock()

	sync, _, board := GetSync(c)
	if !sync {
		return errors.New("meguTV: not synced")
	}

	tvf, ok := feeds.tvFeeds[board]
	if !ok {
		tvf = &tvFeed{}
		tvf.init()
		feeds.tvFeeds[board] = tvf
		err = tvf.start(board)
		if err != nil {
			return
		}
	}
	tvf.add <- c
	return
}

func HandleNekoTV(c common.Client, data []byte) (err error) {
	//Add user to nekotv
	feeds.mu.Lock()
	defer feeds.mu.Unlock()
	synced, thread, _ := GetSync(c)
	if !synced {
		return
	}
	nekoTVFeed, ok := feeds.nekotvFeeds[thread]
	if data[0] == 1 {
		if ok {
			nekoTVFeed.add <- c
		} else {
			var locked bool
			locked, err = db.CheckThreadLocked(thread)
			if err != nil {
				return
			}
			if !locked {
				newFeed := NewNekoTVFeed()
				feeds.nekotvFeeds[thread] = newFeed
				newFeed.start(thread)
				newFeed.add <- c
			}
		}
	} else if data[0] == 0 {
		if ok {
			nekoTVFeed.remove <- c
			rem := <-nekoTVFeed.remove
			if rem != nil {
				delete(feeds.nekotvFeeds, thread)
			}
		}
	}
	return
}

func ToggleNekoTVLock(lock *pb.SetPlaylistLock) {
	feeds.mu.RLock()
	defer feeds.mu.RUnlock()
	feed, ok := feeds.nekotvFeeds[uint64(lock.Post)]
	if !ok {
		return
	}
	feed.SetIsOpen(lock.IsOpen)
}

// Remove client from a subscribed feed
func removeFromFeed(id uint64, board string, c common.Client) {
	feeds.mu.Lock()
	defer feeds.mu.Unlock()

	if feed := feeds.feeds[id]; feed != nil {
		feed.remove <- c
		if nil != <-feed.remove {
			delete(feeds.feeds, feed.id)
		}
	}

	if feed := feeds.tvFeeds[board]; feed != nil {
		feed.remove <- c
		if nil != <-feed.remove {
			delete(feeds.tvFeeds, feed.board)
		}
	}
}

// SendTo sends a message to a feed, if it exists
func SendTo(id uint64, msg []byte) {
	SendIfExists(id, func(f *Feed) error {
		f.Send(msg)
		return nil
	})
}

// Run a send function of a feed, if it exists
func SendIfExists(id uint64, fn func(*Feed) error) error {
	feeds.mu.RLock()
	defer feeds.mu.RUnlock()

	if feed := feeds.feeds[id]; feed != nil {
		fn(feed)
	}
	return nil
}

// InsertPostInto inserts a post into a tread feed, if it exists. Only use for
// already closed posts.
func InsertPostInto(post common.StandalonePost, msg []byte) {
	SendIfExists(post.OP, func(f *Feed) error {
		f.InsertPost(post.Post, msg)
		return nil
	})
}

// ClosePost closes a post in a feed, if it exists
func ClosePost(id, op uint64, links []common.Link, commands []common.Command, claude *common.ClaudeState) (err error) {
	msg, err := common.EncodeMessage(common.MessageClosePost, struct {
		ID       uint64              `json:"id"`
		Links    []common.Link       `json:"links"`
		Commands []common.Command    `json:"commands"`
		Claude   *common.ClaudeState `json:"claude"`
	}{
		ID:       id,
		Links:    links,
		Commands: commands,
		Claude:   claude,
	})
	if err != nil {
		return
	}

	SendIfExists(op, func(f *Feed) error {
		f.ClosePost(id, msg)
		return nil
	})

	return
}

// Initialize internal runtime
func Init() (err error) {
	return db.Listen("post_moderated", func(msg string) (err error) {
		return handlePostModeration(msg)
	})
}

func GetNekoTVFeed(id uint64) (f *NekoTVFeed) {
	feeds.mu.RLock()
	feed, ok := feeds.nekotvFeeds[id]
	feeds.mu.RUnlock()
	if ok {
		return feed
	}
	newFeed := NewNekoTVFeed()
	newFeed.start(id)
	feeds.mu.Lock()
	defer feeds.mu.Unlock()
	feeds.nekotvFeeds[id] = newFeed
	return newFeed
}

// Separate function for testing
func handlePostModeration(msg string) (err error) {
	arr, err := db.SplitUint64s(msg, 2)
	if err != nil {
		return
	}
	op, logID := arr[0], arr[1]
	return SendIfExists(op, func(f *Feed) (err error) {
		e, err := db.GetModLogEntry(logID)
		if err != nil {
			return
		}

		msg, err := common.EncodeMessage(common.MessageModeratePost, struct {
			ID uint64 `json:"id"`
			common.ModerationEntry
		}{e.ID, e.ModerationEntry})
		if err != nil {
			return
		}

		f._moderatePost(e.ID, msg, e.ModerationEntry)
		return
	})
}

// Clear removes all existing feeds and clients. Used only in tests.
func Clear() {
	feeds.mu.Lock()
	defer feeds.mu.Unlock()
	feeds.feeds = make(map[uint64]*Feed, 32)
}
