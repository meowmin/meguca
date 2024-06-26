package feeds

import (
	"errors"
	"github.com/bakape/meguca/common"
	"github.com/bakape/meguca/db"
	"github.com/bakape/meguca/pb"
	"github.com/bakape/meguca/websockets/feeds/nekotv"
	"github.com/go-playground/log"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
	"time"
)

type NekoTVFeed struct {
	baseFeed
	videoTimer *nekotv.VideoTimer
	videoList  *nekotv.VideoList
	thread     uint64
	isRunning  bool
	actions    chan func()
	ticker     *time.Ticker
	isPaused   bool
}

func NewNekoTVFeed() *NekoTVFeed {
	nf := NekoTVFeed{
		videoTimer: nekotv.NewVideoTimer(),
		videoList:  nekotv.NewVideoList(),
	}
	nf.baseFeed.init()
	nf.actions = make(chan func(), 10)
	return &nf
}

func (f *NekoTVFeed) start(thread uint64) (err error) {
	log.Info("Starting NekoTV feed for thread ", thread)
	f.thread = thread
	f.isRunning = true
	state, dbErr := db.GetNekoTVState(thread)
	if dbErr == nil {
		if state.Timer != nil {
			f.videoList.SetItems(state.VideoList)
			f.videoList.SetPos(int(state.ItemPos))
			f.videoTimer.FromProto(state.Timer)
		}
	}
	f.ticker = time.NewTicker(time.Second)
	go f.mainLoop()
	return
}

func (f *NekoTVFeed) mainLoop() {
	for {
		select {
		case c := <-f.add:
			f.addClient(c)
			f.sendConnectedMessage(c)
		case c := <-f.remove:
			if f.removeClient(c) {
				f.isRunning = false
				log.Info("shutting down feed for thread ", f.thread)
				db.DeleteNekoTVValue(f.thread)
				return
			}
		case action := <-f.actions:
			action()
		case <-f.ticker.C:
			f.syncVideoState()
		}
	}
}

func (f *NekoTVFeed) syncVideoState() {
	item, err := f.videoList.CurrentItem()
	if err != nil {
		return
	}
	maxTime := item.Duration - 0.01
	if f.videoTimer.GetTime() > maxTime {
		f.videoTimer.Pause()
		f.videoTimer.SetTime(maxTime)
		skipUrl := item.Url
		time.AfterFunc(time.Second, func() {
			f.actions <- func() {
				if f.videoList.Length() == 0 {
					return
				}
				currentItem, err := f.videoList.CurrentItem()
				if err != nil || currentItem.Url != skipUrl {
					return
				}
				f.SkipVideo()
				f.Play()
			}
		})
		return
	}
	if f.videoList.Length() != 0 {
		f.SendTimeSyncMessage()
	}
}

func (e *NekoTVFeed) GetCurrentState() *pb.ServerState {
	return &pb.ServerState{
		VideoList: e.videoList.GetItems(),
		ItemPos:   int32(e.videoList.Pos),
		Timer:     e.videoTimer.ToProto(),
	}
}

func (f *NekoTVFeed) WriteStateToDb() {
	if f.videoList.Length() == 0 {
		db.DeleteNekoTVValue(f.thread)
	} else {
		db.SetNekoTVState(f.thread, f.GetCurrentState())
	}
}

func (f *NekoTVFeed) sendConnectedMessage(c common.Client) {
	conMessage := pb.ConnectedEvent{
		VideoList:      f.videoList.GetItems(),
		ItemPos:        int32(f.videoList.Pos),
		IsPlaylistOpen: true,
		GetTime:        f.videoTimer.GetTimeData(),
	}
	wsMessage := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_ConnectedEvent{ConnectedEvent: &conMessage}}
	data, err := proto.Marshal(&wsMessage)
	data = append(data, uint8(common.MessageNekoTV))
	if err != nil {
		return
	}
	c.SendBinary(data)
}

func (f *NekoTVFeed) AddVideo(v *pb.VideoItem, atEnd bool) {

	if f.videoList.Exists(func(item *pb.VideoItem) bool {
		return item.Url == v.Url
	}) {
		return
	}
	f.videoList.AddItem(v, atEnd)
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_AddVideoEvent{AddVideoEvent: &pb.AddVideoEvent{
		Item:  v,
		AtEnd: atEnd,
	}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	if f.videoList.Length() == 1 {
		f.videoTimer.Start()
	}
	f.WriteStateToDb()
}

// RemoveVideo removes a video from the playlist
func (f *NekoTVFeed) RemoveVideo(url string) {
	if !f.videoList.IsOpen {
		return
	}

	index := f.videoList.FindIndex(func(item *pb.VideoItem) bool {
		return item.Url == url
	})
	if index == -1 {
		return
	}

	f.videoList.RemoveItem(index)
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_RemoveVideoEvent{RemoveVideoEvent: &pb.RemoveVideoEvent{
		Url: url,
	}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	f.WriteStateToDb()
}

// SkipVideo skips to the next video in the playlist
func (f *NekoTVFeed) SkipVideo() {

	if !f.videoList.IsOpen {
		return
	}
	if f.videoList.Length() == 0 {
		return
	}

	currentItem, err := f.videoList.CurrentItem()
	if err != nil {
		return
	}

	isEmpty := f.videoList.SkipItem()
	if isEmpty {
		f.videoTimer.Stop()
	} else {
		f.videoTimer.SetTime(0)
	}
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_SkipVideoEvent{SkipVideoEvent: &pb.SkipVideoEvent{
		Url: currentItem.Url,
	}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	f.WriteStateToDb()
}

// Pause pauses the current video
func (f *NekoTVFeed) Pause() {

	if !f.videoList.IsOpen {
		return
	}
	if f.videoList.Length() == 0 {
		return
	}

	f.videoTimer.Pause()
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_PauseEvent{PauseEvent: &pb.PauseEvent{
		Time: f.videoTimer.GetTime(),
	}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	f.WriteStateToDb()
}

// Play plays the current video or resumes if paused
func (f *NekoTVFeed) Play() {

	if f.videoList.Length() == 0 {
		return
	}

	time := f.videoTimer.GetTime()
	f.videoTimer.Play()
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_PlayEvent{PlayEvent: &pb.PlayEvent{
		Time: time,
	}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	f.WriteStateToDb()
}

// SetTime sets the current playback time
func (f *NekoTVFeed) SetTime(time float32) {

	if !f.videoList.IsOpen {
		return
	}
	if f.videoList.Length() == 0 {
		return
	}

	f.videoTimer.SetTime(time)
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_SetTimeEvent{SetTimeEvent: &pb.SetTimeEvent{
		Time: time,
	}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	f.WriteStateToDb()
}

// UpdatePlaylist updates the playlist
func (f *NekoTVFeed) UpdatePlaylist() {
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_UpdatePlaylistEvent{UpdatePlaylistEvent: &pb.UpdatePlaylistEvent{
		VideoList: &pb.VideoItemList{
			Items: f.videoList.GetItems(),
		},
	}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	f.WriteStateToDb()
}

// ClearPlaylist clears the playlist
func (f *NekoTVFeed) ClearPlaylist() {

	if !f.videoList.IsOpen {
		return
	}
	f.videoList.Clear()
	f.videoTimer.Stop()
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_ClearPlaylistEvent{ClearPlaylistEvent: &pb.ClearPlaylistEvent{}}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
	db.DeleteNekoTVValue(f.thread)
}

func (f *NekoTVFeed) SendTimeSyncMessage() {
	msg := pb.WebSocketMessage{MessageType: &pb.WebSocketMessage_GetTimeEvent{GetTimeEvent: f.videoTimer.GetTimeData()}}
	data, _ := proto.Marshal(&msg)
	data = append(data, uint8(common.MessageNekoTV))
	f.sendToAllBinary(data)
}
func (f *NekoTVFeed) GetIsOpen() bool {
	return f.videoList.IsOpen
}
func (f *NekoTVFeed) SetIsOpen(b bool) {
	f.videoList.IsOpen = b
}
func parseTimestamp(timestamp string) (time float32, err error) {
	if strings.Contains(timestamp, ":") {
		parts := strings.Split(timestamp, ":")
		if len(parts) == 2 {
			var minutes int
			minutes, err = strconv.Atoi(parts[0])
			if err != nil {
				return
			}
			time = float32(minutes * 60)
			var seconds float64
			seconds, err = strconv.ParseFloat(parts[1], 32)
			time += float32(seconds)
			return
		} else if len(parts) == 3 {
			var hours int
			hours, err = strconv.Atoi(parts[0])
			if err != nil {
				return
			}
			time = float32(hours * 60 * 60)
			var minutes int
			minutes, err = strconv.Atoi(parts[0])
			if err != nil {
				return
			}
			time += float32(minutes * 60)
			var seconds float64
			seconds, err = strconv.ParseFloat(parts[1], 32)
			time += float32(seconds)
			return
		} else {
			err = errors.New("invalid timestamp")
			return
		}
	}
	var seconds float64
	seconds, err = strconv.ParseFloat(timestamp, 32)
	time = float32(seconds)
	return
}

func HandleMediaCommand(thread uint64, c *common.MediaCommand) {
	feeds.mu.RLock()
	ntv, ok := feeds.nekotvFeeds[thread]
	feeds.mu.RUnlock()
	if !ok {
		return
	}
	switch c.Type {
	case common.AddVideo:
		ntv.actions <- func() {
			videoData, err := nekotv.GetVideoData(c.Args)
			if err == nil {
				log.Infof("Video data retrieved: %v", videoData)
				ntv.AddVideo(&videoData, true)
			} else {
				log.Errorf("Failed to get video data: %v", err)
			}
		}
		break
	case common.RemoveVideo:
		ntv.actions <- func() {
			ntv.RemoveVideo(c.Args)
		}
	case common.SkipVideo:
		ntv.actions <- func() {
			ntv.SkipVideo()
		}
	case common.Pause:
		ntv.actions <- func() {
			ntv.Pause()
		}
	case common.Play:
		ntv.actions <- func() {
			ntv.Play()
		}
	case common.SetTime:
		time, err := parseTimestamp(c.Args)
		if err != nil {
			log.Errorf("Failed to parse timestamp: %v", err)
		} else {
			ntv.actions <- func() {
				ntv.SetTime(time)
			}
		}
	case common.ClearPlaylist:
		ntv.actions <- func() {
			ntv.ClearPlaylist()
		}
	default:
		log.Warnf("Unknown media command type: %v", c.Type)
	}
}
