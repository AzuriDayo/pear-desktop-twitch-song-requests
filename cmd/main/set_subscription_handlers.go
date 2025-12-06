package main

import (
	"encoding/json"
	"log"

	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/labstack/echo/v4"
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
		log.Printf("Chat message from %s: %s %s\n", event.ChatterUserLogin, event.Message.Text, event.ChannelPointsCustomRewardId)
	})
}
