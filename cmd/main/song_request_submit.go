package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/nicklaw5/helix/v2"
)

func (a *App) songRequestSubmit(useProperHelix *helix.Client, properUserID string, event twitch.EventChannelChatMessage) {
	s := songrequests.ParseSearchQuery(event.Message.Text)
	song, err := songrequests.SearchSong(strings.TrimPrefix(s, "-"), 60, 600)
	if err != nil {
		log.Println("Searching for song failed: "+event.Message.Text+"\n", err)
		emsg := "Failed to search for your song"
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              emsg,
			ReplyParentMessageID: event.MessageId,
		})
		return
	}

	// Loop through queue state to check if song is queued already
	queue := songrequests.QueueResponse{}

	preResponse, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/queue")
	if err != nil || preResponse.StatusCode != http.StatusOK {
		emsg := "Internal error when checking if song is already in queue"
		log.Println(emsg, err)
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              emsg,
			ReplyParentMessageID: event.MessageId,
		})
		return
	}
	qb, err := io.ReadAll(preResponse.Body)
	if err != nil {
		emsg := "Internal error processing data to check if song is already in queue"
		log.Println(emsg, err)
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              emsg,
			ReplyParentMessageID: event.MessageId,
		})
		return
	}
	err = json.Unmarshal(qb, &queue)
	preResponse.Body.Close()
	if err != nil {
		emsg := "Failed to check if song exists in queue."
		log.Println(emsg, err)
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              emsg,
			ReplyParentMessageID: event.MessageId,
		})
		return
	}

	afterSelected := false
	songExistsInQueue := false
	for _, v := range queue.Items {
		selected := false
		compareVideoIDs := map[string]struct{}{}
		if v.PlaylistPanelVideoWrapperRenderer != nil {
			compareVideoIDs[v.PlaylistPanelVideoWrapperRenderer.PrimaryRenderer.PlaylistPanelVideoRenderer.VideoId] = struct{}{}
			if v.PlaylistPanelVideoWrapperRenderer.PrimaryRenderer.PlaylistPanelVideoRenderer.Selected {
				selected = true
			}
			for _, v2 := range v.PlaylistPanelVideoWrapperRenderer.Counterpart {
				compareVideoIDs[v2.CounterpartRenderer.PlaylistPanelVideoRenderer.VideoId] = struct{}{}
			}
		}
		if v.PlaylistPanelVideoRenderer != nil {
			compareVideoIDs[v.PlaylistPanelVideoRenderer.VideoId] = struct{}{}
			if v.PlaylistPanelVideoRenderer.Selected {
				selected = true
			}
		}
		if selected {
			afterSelected = true
		}
		if _, ok := compareVideoIDs[song.VideoID]; afterSelected && ok {
			songExistsInQueue = true
			break
		}
	}

	// done check raw client queue, now check internal queue
	songQueueMutex.Lock()
	for _, v := range songQueue {
		if song.VideoID == v.Song.VideoID {
			songExistsInQueue = true
			break
		}
	}
	songQueueMutex.Unlock()

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
	a.songRequestLogic(song, event, properUserID, useProperHelix)
}
