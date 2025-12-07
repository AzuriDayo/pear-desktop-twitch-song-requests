package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"errors"
	"net/http"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/gen/model"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	"github.com/nicklaw5/helix/v2"

	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	. "github.com/go-jet/jet/v2/sqlite"
)

func (a *App) loadSqliteSettings() error {
	db, err := databaseconn.NewDBConnection()
	if err != nil {
		return err
	}
	defer func() {
		db.Close()
	}()

	results := []model.Settings{}
	stmt := SELECT(Settings.Key, Settings.Value).FROM(Settings)
	err = stmt.QueryContext(a.ctx, db, &results)
	if err != nil {
		return err
	}

	for _, result := range results {
		if result.Key == data.DB_KEY_TWITCH_ACCESS_TOKEN {
			a.twitchDataStruct.accessToken = result.Value
		}
		if result.Key == data.DB_KEY_TWITCH_SONG_REQUEST_REWARD_ID {
			a.songRequestRewardID = result.Value
		}
	}

	if a.twitchDataStruct.accessToken != "" {
		isValid, response, err := a.helix.ValidateToken(a.twitchDataStruct.accessToken)
		if err != nil {
			// req error
			return err
		}
		if response.StatusCode == http.StatusOK && isValid {
			expiresIn := response.Data.ExpiresIn
			strDate := response.Header.Get("Date")
			t, err := time.Parse(data.TWITCH_SERVER_DATE_LAYOUT, strDate)
			if err != nil {
				return errors.New("Failed to validate server date time expiry, original error:\n" + err.Error())
			}
			t = t.Add(time.Duration(expiresIn) * time.Second)
			a.helix.SetUserAccessToken(a.twitchDataStruct.accessToken)
			a.twitchDataStruct.expiresDate = t
			a.twitchDataStruct.isAuthenticated = true
			a.twitchDataStruct.userID = response.Data.UserID
			a.twitchDataStruct.login = response.Data.Login

			resp, err := a.helix.GetStreams(&helix.StreamsParams{
				UserLogins: []string{a.twitchDataStruct.login},
			})
			if err == nil && len(resp.Data.Streams) > 0 && resp.Data.Streams[0].ID != "" {
				a.streamOnline = true
			}
		}
	}

	return nil
}
