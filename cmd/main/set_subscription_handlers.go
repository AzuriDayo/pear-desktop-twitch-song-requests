package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/labstack/echo/v4"
	"github.com/nicklaw5/helix/v2"
)

func (a *App) SetSubscriptionHandlers() {
	a.twitchWSService.Client().OnEventStreamOnline(func(event twitch.EventStreamOnline) {
		a.streamOnline = true

		j, _ := json.Marshal(echo.Map{
			"stream_online": true,
		})
		a.clientsBroadcast <- string(j)
		log.Println("STREAM_ONLINE")
	})
	a.twitchWSService.Client().OnEventStreamOffline(func(event twitch.EventStreamOffline) {
		a.streamOnline = false
		j, _ := json.Marshal(echo.Map{
			"stream_online": false,
		})
		a.clientsBroadcast <- string(j)
		log.Println("STREAM_OFFLINE")
	})
	a.twitchWSService.Client().OnEventChannelChatMessage(func(event twitch.EventChannelChatMessage) {
		// Song request
		isSub := false
		for _, v := range event.Badges {
			if v.SetId == "subscriber" {
				isSub = true
			}
		}
		log.Printf("Chat message from %s: %s %s\n", event.ChatterUserLogin, event.Message.Text, event.ChannelPointsCustomRewardId)
		if a.streamOnline {
			if a.songRequestRewardID == event.ChannelPointsCustomRewardId || (strings.HasPrefix(event.Message.Text, "!sr ") && isSub) {
				s := songrequests.ParseSearchQuery(event.Message.Text)
				song, err := songrequests.SearchSong(s, 600)
				if err != nil {
					a.helix.SendChatMessage(&helix.SendChatMessageParams{
						BroadcasterID:        event.BroadcasterUserId,
						SenderID:             a.twitchDataStruct.userID,
						Message:              err.Error(),
						ReplyParentMessageID: event.MessageId,
					})
					return
				}
				log.Println("Searched song: " + song.Title + " - " + song.Artist + ": " + song.VideoID)
				songrequests.QueueNextSong(song.VideoID)
			}
		}
	})
}
