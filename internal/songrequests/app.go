package songrequests

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/labstack/echo/v4"
)

var queueNextSong chan struct {
	VideoID string
}

var RcvMsgChan chan twitch.EventChannelChatMessage

type SongQueueAction = string

const (
	SongQueueActionAddSong    SongQueueAction = "ADD_SONG"
	SongQueueActionRemoveSong SongQueueAction = "REMOVE_SONG"
	SongQueueActionMoveSong   SongQueueAction = "MOVE_SONG"
)

type SongQueueInteraction struct {
	Action SongQueueAction
}

var SongQueueInteractionChan chan SongQueueInteraction

type SongRequestQueueItem struct {
	SongTitle string
	Artist    string
	VideoID   string
}

var SongRequestQueue []SongRequestQueueItem

func RunGoroutineAddNextSong(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-queueNextSong:
			b := echo.Map{
				"videoId":        msg.VideoID,
				"insertPosition": "INSERT_AFTER_CURRENT_VIDEO",
			}
			bb, _ := json.Marshal(b)
			http.Post("http://"+GetPearDesktopHost()+"/api/v1/queue", "application/json", bytes.NewBuffer(bb))
			log.Println("added video id " + msg.VideoID)
		}
	}
}

func QueueNextSong(videoId string) {
	queueNextSong <- struct{ VideoID string }{VideoID: videoId}
}

func init() {
	queueNextSong = make(chan struct {
		VideoID string
	})

	RcvMsgChan = make(chan twitch.EventChannelChatMessage)

	SongQueueInteractionChan = make(chan SongQueueInteraction)

	SongRequestQueue = []SongRequestQueueItem{}
}
