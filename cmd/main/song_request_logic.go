package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/gen/model"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/labstack/echo/v4"
	"github.com/nicklaw5/helix/v2"
	"golang.org/x/net/websocket"
)

func (a *App) songRequestLogic(song *songrequests.SongResult, event twitch.EventChannelChatMessage, properUserID string, useProperHelix *helix.Client) {
	songQueueMutex.Lock()
	defer songQueueMutex.Unlock()

	songQueueItem := SongQueueItem{
		RequestedBy: event.ChatterUserLogin,
		Song:        *song,
		IsNinja:     strings.EqualFold(event.BroadcasterUserLogin, a.twitchDataStructBot.login),
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
		if strings.EqualFold(event.BroadcasterUserLogin, a.twitchDataStructBot.login) {
			log.Println("hehe chatter " + event.ChatterUserLogin + ": Queued song " + song.Title + " - " + song.Artist)
		} else {
			log.Println(event.ChatterUserLogin + ": Queued song " + song.Title + " - " + song.Artist)
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
	go func() {
		db, err := databaseconn.NewDBConnection()
		if err != nil {
			log.Println("Somehow failed to add !sr history to database")
			return
		}
		srData := model.SongRequests{
			VideoID:    song.VideoID,
			SongTitle:  song.Title,
			ArtistName: song.Artist,
			ImageURL:   song.ImageUrl,
		}
		stmt := SongRequests.INSERT(SongRequests.AllColumns).MODEL(srData).ON_CONFLICT(SongRequests.VideoID).DO_NOTHING()
		_, err = stmt.Exec(db)
		if err != nil {
			log.Println("Somehow failed to save !sr song to database")
			return
		}
		srrData := model.SongRequestRequesters{
			VideoID:        song.VideoID,
			TwitchUsername: event.ChatterUserLogin,
			RequestedAt:    time.Now().Local().Format(data.TWITCH_SERVER_DATE_LAYOUT),
			UserID:         event.ChatterUserId,
			IsNinja:        strings.EqualFold(event.BroadcasterUserLogin, a.twitchDataStructBot.login),
		}
		stmt = SongRequestRequesters.INSERT(SongRequestRequesters.AllColumns).MODEL(srrData)
		_, err = stmt.Exec(db)
		if err != nil {
			log.Println("Somehow failed to save !sr requester name to database")
			return
		}
	}()

	useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
		BroadcasterID:        event.BroadcasterUserId,
		SenderID:             properUserID,
		Message:              "Added song: " + song.Title + " - " + song.Artist + " " + "https://youtu.be/" + song.VideoID,
		ReplyParentMessageID: event.MessageId,
	})
}
