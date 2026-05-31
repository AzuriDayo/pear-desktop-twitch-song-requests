package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/nicklaw5/helix/v2"
)

func (a *App) songRequestSubmit(useProperHelix *helix.Client, properUserID string, event twitch.EventChannelChatMessage) {
	s, isNativeVideoID := songrequests.ParseSearchQuery(event.Message.Text)
	minTimeS := 30
	maxTimeS := 600
	song, err := songrequests.SearchSong(strings.TrimPrefix(s, "-"), minTimeS, maxTimeS)
	if err != nil {
		log.Println("Searching for song failed: "+event.Message.Text+"\n", err)
		emsg := "Failed to search for your song"
		switch err {
		case songrequests.ErrNoResults:
			emsg = "Error: no songs found"
		case songrequests.ErrSongLength:
			emsg = fmt.Sprintf("Error: the song must be between %ds and %ds", minTimeS, maxTimeS)
		}
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              emsg,
			ReplyParentMessageID: event.MessageId,
		})
		return
	}

	requestedStringIsSameVideoID := true
	if isNativeVideoID && !strings.Contains(s, song.VideoID) {
		song.VideoID = s
		song.Artist = "unknown"
		song.ImageUrl = ""
		song.RawTimeData = ""
		song.Title = "unknown"
		song.IsUnknown = true
		requestedStringIsSameVideoID = false

		// TODO: cleanup isUnknown code
		// For now simply return early and deny song request
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              "Sorry, this video was probably available in YouTube, but not in YouTube Music.",
			ReplyParentMessageID: event.MessageId,
		})
		return
	}

	// Loop through queue state to check if song is queued already
	nowPlayingSong := songrequests.SongResult{}

	resp, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/song")
	if err != nil {
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              "Failed to get details for currently playing song.",
			ReplyParentMessageID: event.MessageId,
		})
		return
	}
	bb, err := io.ReadAll(resp.Body)
	if err != nil {
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              "Failed to read details for currently playing song.",
			ReplyParentMessageID: event.MessageId,
		})
		return
	}
	err = json.Unmarshal(bb, &nowPlayingSong)
	if err != nil {
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              "Failed to check data for currently playing song.",
			ReplyParentMessageID: event.MessageId,
		})
		return
	}

	// nowPlayingSong now is correct
	songExistsInQueue := false
	if nowPlayingSong.VideoID == song.VideoID {
		songExistsInQueue = true
	}

	// done check raw client queue, now check internal queue
	if !songExistsInQueue {
		songQueueMutex.Lock()
		for _, v := range songQueue {
			if song.VideoID == v.Song.VideoID {
				songExistsInQueue = true
				break
			}
		}
		songQueueMutex.Unlock()
	}

	if songExistsInQueue {
		msg := "Song is already in queue!"
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              msg,
			ReplyParentMessageID: event.MessageId,
		})
		return
	}

	// Committing to adding song to q
	a.songRequestLogic(song, requestedStringIsSameVideoID, event, properUserID, useProperHelix)
}
