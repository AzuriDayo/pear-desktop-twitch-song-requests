package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"slices"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/nicklaw5/helix/v2"
)

func (a *App) songRequestSubmit(useProperHelix *helix.Client, properUserID string, event twitch.EventChannelChatMessage) {
	s := songrequests.ParseSearchQuery(event.Message.Text)
	song, err := songrequests.SearchSong(s, 60, 600)
	if err != nil {
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
		compareVideoIDs := []string{}
		if v.PlaylistPanelVideoWrapperRenderer != nil {
			compareVideoIDs = append(compareVideoIDs, v.PlaylistPanelVideoWrapperRenderer.PrimaryRenderer.PlaylistPanelVideoRenderer.VideoId)
			if v.PlaylistPanelVideoWrapperRenderer.PrimaryRenderer.PlaylistPanelVideoRenderer.Selected {
				selected = true
			}
			for _, v2 := range v.PlaylistPanelVideoWrapperRenderer.Counterpart {
				compareVideoIDs = append(compareVideoIDs, v2.CounterpartRenderer.PlaylistPanelVideoRenderer.VideoId)
			}
		}
		if v.PlaylistPanelVideoRenderer != nil {
			compareVideoIDs = append(compareVideoIDs, v.PlaylistPanelVideoRenderer.VideoId)
			if v.PlaylistPanelVideoRenderer.Selected {
				selected = true
			}
		}
		if selected {
			afterSelected = true
		}
		if afterSelected && slices.Contains(compareVideoIDs, song.VideoID) {
			songExistsInQueue = true
			break
		}
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
	useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
		BroadcasterID:        event.BroadcasterUserId,
		SenderID:             properUserID,
		Message:              "Added song: " + song.Title + " - " + song.Artist + " " + "https://youtu.be/" + song.VideoID,
		ReplyParentMessageID: event.MessageId,
	})
	srChan <- struct {
		song  *songrequests.SongResult
		event twitch.EventChannelChatMessage
	}{
		song:  song,
		event: event,
	}
}
