package data

import "time"

var twitchClientID = "7k7nl6w8e0owouonj7nb9g3k5s6gs5"

func GetTwitchClientID() string {
	return twitchClientID
}

const (
	DB_KEY_TWITCH_ACCESS_TOKEN           = "twitch_access_token"
	DB_KEY_TWITCH_ACCESS_TOKEN_BOT       = "twitch_access_token_bot"
	DB_KEY_TWITCH_REFRESH_TOKEN               = "twitch_refresh_token"
	DB_KEY_TWITCH_REFRESH_TOKEN_BOT           = "twitch_refresh_token_bot"
	DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED     = "twitch_refresh_token_last_used"
	DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED_BOT = "twitch_refresh_token_last_used_bot"
	DB_KEY_TWITCH_SONG_REQUEST_REWARD_ID = "twitch_song_request_reward_id"
	TWITCH_SERVER_DATE_LAYOUT            = time.RFC1123
	TWITCH_OAUTH_REDIRECT_URI            = "http://localhost:3999/oauth/twitch"
	TWITCH_OAUTH_LISTEN_ADDR             = "127.0.0.1:3999"

	// Command permission DB keys
	DB_KEY_CMD_PERMISSION_SR      = "cmd_permission_sr"
	DB_KEY_CMD_PERMISSION_QUEUE   = "cmd_permission_queue"
	DB_KEY_CMD_PERMISSION_SONG    = "cmd_permission_song"
	DB_KEY_CMD_PERMISSION_DELSONG = "cmd_permission_delsong"

	// Permission levels (0 = most restrictive, 4 = everyone)
	PermissionLevelBroadcaster = 0
	PermissionLevelModerator   = 1
	PermissionLevelVIP         = 2
	PermissionLevelSubscriber  = 3
	PermissionLevelViewer      = 4
)
