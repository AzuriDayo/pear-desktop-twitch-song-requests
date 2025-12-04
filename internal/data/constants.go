package data

import "os"

var twitchClientID = "7k7nl6w8e0owouonj7nb9g3k5s6gs5"
var pearDesktopHostname = "127.0.0.1"
var pearDesktopPort = "26538"

func GetTwitchClientID() string {
	return twitchClientID
}

func GetTwitchClientSecret() string {
	return os.Getenv("TWITCH_CLIENT_SECRET")
}

func GetPearDesktopHost() string {
	return pearDesktopHostname + ":" + pearDesktopPort
}
