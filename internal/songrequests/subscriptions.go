package songrequests

import (
	"log"

	"github.com/joeyak/go-twitch-eventsub/v3"
)

func GetSubscriptions() []twitch.EventSubscription {
	events := []twitch.EventSubscription{
		twitch.SubStreamOnline,
		twitch.SubStreamOffline,
		twitch.SubChannelChatMessage,
		twitch.SubChannelChannelPointsCustomRewardRedemptionAdd, // claim reward points
	}

	return events
}

func SetSubscriptionHandlers(client *twitch.Client) {
	client.OnEventStreamOnline(func(event twitch.EventStreamOnline) {
		log.Printf("Stream online: %s\n", event.BroadcasterUserName)
	})
	client.OnEventStreamOffline(func(event twitch.EventStreamOffline) {
		log.Printf("Stream offline: %s\n", event.BroadcasterUserName)
	})
	client.OnEventChannelChatMessage(func(event twitch.EventChannelChatMessage) {
		log.Printf("Chat message from %s: %s %s\n", event.ChatterUserLogin, event.Message.Text, event.ChannelPointsCustomRewardId)
	})
}
