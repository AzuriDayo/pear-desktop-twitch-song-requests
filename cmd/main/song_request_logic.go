package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/labstack/echo/v4"
	"github.com/nicklaw5/helix/v2"
	"golang.org/x/net/websocket"
)

func (a *App) songRequestLogic(song *songrequests.SongResult, requestedStringIsSameVideoID bool, event twitch.EventChannelChatMessage, properUserID string, useProperHelix *helix.Client) {
	songQueueMutex.Lock()
	defer songQueueMutex.Unlock()

	songQueueItem := SongQueueItem{
		RequestedBy:       event.ChatterUserLogin,
		RequestedByUserID: event.ChatterUserId,
		Song:              *song,
		IsNinja:           strings.EqualFold(event.BroadcasterUserLogin, a.twitchDataStructBot.login),
	}
	songQueue = append(songQueue, songQueueItem)

	if len(songQueue) == 1 {
		// 1st new queue item must immediately add to next in raw player
		b := echo.Map{
			"videoId":        song.VideoID,
			"insertPosition": "INSERT_AFTER_CURRENT_VIDEO",
		}
		bb, _ := json.Marshal(b)
		resp, err := http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue", "application/json", bytes.NewBuffer(bb))
		if err != nil || resp.StatusCode != http.StatusNoContent {
			emsg := "Internal error when adding song to queue"
			log.Println(emsg, err)
			useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
				BroadcasterID:        event.BroadcasterUserId,
				SenderID:             properUserID,
				Message:              emsg,
				ReplyParentMessageID: event.MessageId,
			})
			return
		}

		// validate if song was really added
		nextInQueueIsAdded := false
		intervalDelay := time.Second
		maxRetries := 5
		for range maxRetries {
			time.Sleep(intervalDelay)
			found, _ := helpers.FindAllVideoIDCounterparts(songQueue[0].Song.VideoID)
			if found {
				nextInQueueIsAdded = true
				break
			}
		}

		// Notify failed
		if !nextInQueueIsAdded {
			senderID := a.twitchDataStruct.userID
			if a.twitchDataStructBot.isAuthenticated {
				senderID = a.twitchDataStructBot.userID
			}
			a.helix.SendChatMessage(&helix.SendChatMessageParams{
				BroadcasterID: a.twitchDataStruct.userID,
				SenderID:      senderID,
				Message:       "Sorry " + songQueue[0].RequestedBy + " , I failed to add https://youtu.be/" + songQueue[0].Song.VideoID + " next and will be removed from queue.",
			})

			songQueue = songQueue[1:]
			return
		}
	}

	// Broadcast song added to queue to browser control panel
	queueInfoOnAdd, _ := json.Marshal(echo.Map{
		"type": "QUEUE_ADD",
		"song": songQueueItem,
	})
	a.clientsMu.Lock()
	for ws := range a.clients {
		websocket.Message.Send(ws, string(queueInfoOnAdd))
	}
	a.clientsMu.Unlock()

	// save to history
	if !song.IsUnknown {
		isNinja := strings.EqualFold(event.BroadcasterUserLogin, a.twitchDataStructBot.login)
		go saveSongHistory(song, event.ChatterUserLogin, event.ChatterUserId, isNinja)
	}
	// if it is unknown, it will be saved later

	if strings.EqualFold(event.BroadcasterUserLogin, a.twitchDataStructBot.login) {
		log.Println("hehe chatter " + event.ChatterUserLogin + ": Queued song " + song.Title + " - " + song.Artist)
	} else {
		log.Println(event.ChatterUserLogin + ": Queued song " + song.Title + " - " + song.Artist)
	}

	if requestedStringIsSameVideoID {
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              "Added song: " + song.Title + " - " + song.Artist + " " + "https://youtu.be/" + song.VideoID,
			ReplyParentMessageID: event.MessageId,
		})
	} else {
		useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
			BroadcasterID:        event.BroadcasterUserId,
			SenderID:             properUserID,
			Message:              "Added song: https://youtu.be/" + song.VideoID,
			ReplyParentMessageID: event.MessageId,
		})
	}
}
