package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"strconv"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/gen/model"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"

	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	. "github.com/go-jet/jet/v2/sqlite"
)

func (a *App) loadSqliteSettings() error {
	db, err := databaseconn.NewDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	results := []model.Settings{}
	stmt := SELECT(Settings.Key, Settings.Value).FROM(Settings)
	err = stmt.QueryContext(a.ctx, db, &results)
	if err != nil {
		return err
	}

	permissionKeys := map[string]bool{
		data.DB_KEY_CMD_PERMISSION_SR:      true,
		data.DB_KEY_CMD_PERMISSION_QUEUE:   true,
		data.DB_KEY_CMD_PERMISSION_SONG:    true,
		data.DB_KEY_CMD_PERMISSION_DELSONG: true,
	}

	for _, result := range results {
		if result.Key == data.DB_KEY_TWITCH_ACCESS_TOKEN {
			a.twitchDataStruct.accessToken = result.Value
		}
		if result.Key == data.DB_KEY_TWITCH_REFRESH_TOKEN {
			a.twitchDataStruct.refreshToken = result.Value
		}
		if result.Key == data.DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED {
			if t, err := time.Parse(time.RFC3339, result.Value); err == nil {
				a.twitchDataStruct.refreshTokenLastUsedAt = t
			}
		}
		if result.Key == data.DB_KEY_TWITCH_SONG_REQUEST_REWARD_ID {
			a.songRequestRewardID = result.Value
		}
		if result.Key == data.DB_KEY_TWITCH_ACCESS_TOKEN_BOT {
			a.twitchDataStructBot.accessToken = result.Value
		}
		if result.Key == data.DB_KEY_TWITCH_REFRESH_TOKEN_BOT {
			a.twitchDataStructBot.refreshToken = result.Value
		}
		if result.Key == data.DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED_BOT {
			if t, err := time.Parse(time.RFC3339, result.Value); err == nil {
				a.twitchDataStructBot.refreshTokenLastUsedAt = t
			}
		}
		if permissionKeys[result.Key] {
			if v, err := strconv.Atoi(result.Value); err == nil && v >= data.PermissionLevelBroadcaster && v <= data.PermissionLevelViewer {
				a.cmdPermissions[result.Key] = v
			}
		}
	}

	if err := a.validateLoadedTwitchToken(false); err != nil {
		return err
	}
	if a.twitchDataStructBot.accessToken != "" || a.twitchDataStructBot.refreshToken != "" {
		if err := a.validateLoadedTwitchToken(true); err != nil {
			return err
		}
	}

	a.migrateRefreshTokenLastUsed(false)
	a.migrateRefreshTokenLastUsed(true)

	return nil
}

func (a *App) migrateRefreshTokenLastUsed(forBot bool) {
	td := twitchDataForBot(a, forBot)
	if td.refreshToken == "" || !td.refreshTokenLastUsedAt.IsZero() {
		return
	}
	_ = a.markRefreshTokenUsed(forBot)
}
