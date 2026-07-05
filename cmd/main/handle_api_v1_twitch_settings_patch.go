package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/gen/model"
	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/labstack/echo/v4"
)

// parsePermissionLevel parses a string to a valid permission level integer (0–4).
// Returns an error if the string is not a valid integer in that range.
func parsePermissionLevel(s string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("permission level must be a number, got %q", s)
	}
	if v < data.PermissionLevelBroadcaster || v > data.PermissionLevelViewer {
		return 0, fmt.Errorf("permission level must be between %d and %d, got %d", data.PermissionLevelBroadcaster, data.PermissionLevelViewer, v)
	}
	return v, nil
}

func (a *App) handleApiV1SettingsPATCH(c echo.Context) error {
	// auth data in url hash string params as get request
	body := c.Request().Body
	rawBodyData, err := io.ReadAll(body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "read request body",
		})
	}
	defer body.Close()

	settings := map[string]string{}
	err = json.Unmarshal(rawBodyData, &settings)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "parse request body",
		})
	}
	db, err := databaseconn.NewDBConnection()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "save data failed",
		})
	}
	defer db.Close()

	permissionKeys := map[string]bool{
		data.DB_KEY_CMD_PERMISSION_SR:      true,
		data.DB_KEY_CMD_PERMISSION_QUEUE:   true,
		data.DB_KEY_CMD_PERMISSION_SONG:    true,
		data.DB_KEY_CMD_PERMISSION_DELSONG: true,
	}

	for k, v := range settings {
		if k == data.DB_KEY_TWITCH_SONG_REQUEST_REWARD_ID {
			newSetting := model.Settings{
				Key:   data.DB_KEY_TWITCH_SONG_REQUEST_REWARD_ID,
				Value: v,
			}
			stmt := Settings.INSERT(Settings.AllColumns).MODEL(newSetting).ON_CONFLICT(Settings.Key).DO_UPDATE(SET(
				Settings.Value.SET(String(v)),
			))
			_, err = stmt.ExecContext(c.Request().Context(), db)
			if err != nil {
				log.Println("handleApiV1SettingsPATCH: failed to save setting", k, err)
			}
			a.songRequestRewardID = v
		}

		if permissionKeys[k] {
			level, err := parsePermissionLevel(v)
			if err != nil {
				log.Println("handleApiV1SettingsPATCH: invalid permission level for", k, err)
				continue
			}
			newSetting := model.Settings{Key: k, Value: v}
			stmt := Settings.INSERT(Settings.AllColumns).MODEL(newSetting).ON_CONFLICT(Settings.Key).DO_UPDATE(SET(
				Settings.Value.SET(String(v)),
			))
			_, err = stmt.ExecContext(c.Request().Context(), db)
			if err != nil {
				log.Println("handleApiV1SettingsPATCH: failed to save permission setting", k, err)
				continue
			}
			a.cmdPermissions[k] = level
		}
	}

	expiryDate := ""
	if a.twitchDataStruct.isAuthenticated {
		expiryDate = a.twitchDataStruct.expiresDate.Local().Format(data.TWITCH_SERVER_DATE_LAYOUT)
	}
	expiryDateBot := ""
	if a.twitchDataStructBot.isAuthenticated {
		expiryDateBot = a.twitchDataStructBot.expiresDate.Local().Format(data.TWITCH_SERVER_DATE_LAYOUT)
	}
	b := echo.Map{
		"type":            "TWITCH_INFO",
		"stream_online":   a.streamOnline,
		"reward_id":       a.songRequestRewardID,
		"login":           a.twitchDataStruct.login,
		"login_bot":       a.twitchDataStructBot.login,
		"expiry_date":     expiryDate,
		"expiry_date_bot": expiryDateBot,
	}
	bb, _ := json.Marshal(b)
	a.clientsBroadcast <- string(bb)
	return c.NoContent(http.StatusOK)

}
