package main

import (
	"sync"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
)

var songQueueMutex = sync.RWMutex{}

type SongQueueItem struct {
	RequestedBy       string                  `json:"requested_by"`
	RequestedByUserID string                  `json:"-"`
	Song              songrequests.SongResult `json:"song"`
	IsNinja           bool                    `json:"is_ninja"`
}

var songQueue = []SongQueueItem{}

// QueueInfoPayload is the payload for the QUEUE_INFO event (full queue snapshot).
type QueueInfoPayload struct {
	SongQueue []SongQueueItem `json:"song_queue"`
}

type playerSonginfo struct {
	VideoId          string `json:"videoId"`
	ImageSrc         string `json:"imageSrc"`
	Artist           string `json:"artist"`
	SongDuration     int    `json:"songDuration"`
	AlternativeTitle string `json:"alternativeTitle"`
}

func (s playerSonginfo) GetUrl() string {
	return "https://youtu.be/" + s.VideoId
}

var playerInfo = struct {
	Position  int
	IsPlaying bool
	Song      playerSonginfo `json:"song"`
}{
	Song: playerSonginfo{},
}
