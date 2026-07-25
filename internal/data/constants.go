package data

import "time"

var twitchClientID = "7k7nl6w8e0owouonj7nb9g3k5s6gs5"

func GetTwitchClientID() string {
	return twitchClientID
}

const (
	DB_KEY_TWITCH_ACCESS_TOKEN           = "twitch_access_token"
	DB_KEY_TWITCH_ACCESS_TOKEN_BOT       = "twitch_access_token_bot"
	DB_KEY_TWITCH_SONG_REQUEST_REWARD_ID = "twitch_song_request_reward_id"
	TWITCH_SERVER_DATE_LAYOUT            = time.RFC1123

	// Command permission DB keys
	DB_KEY_CMD_PERMISSION_SR      = "cmd_permission_sr"
	DB_KEY_CMD_PERMISSION_QUEUE   = "cmd_permission_queue"
	DB_KEY_CMD_PERMISSION_SONG    = "cmd_permission_song"
	DB_KEY_CMD_PERMISSION_DELSONG = "cmd_permission_delsong"

	// Command aliases: JSON object mapping alias (e.g. "!q") -> canonical short name (e.g. "queue")
	DB_KEY_CMD_ALIASES = "cmd_aliases"

	// Disabled built-in command names: JSON array of short names (e.g. ["queue","song"]).
	// Aliases targeting a disabled built-in still work; only the original !name is blocked.
	DB_KEY_CMD_DISABLED_BUILTINS = "cmd_disabled_builtins"

	// Permission levels (0 = most restrictive, 4 = everyone)
	PermissionLevelBroadcaster = 0
	PermissionLevelModerator   = 1
	PermissionLevelVIP         = 2
	PermissionLevelSubscriber  = 3
	PermissionLevelViewer      = 4
)
