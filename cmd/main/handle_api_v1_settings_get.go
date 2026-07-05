package main

import (
	"net/http"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/labstack/echo/v4"
)

// handleApiV1SettingsGET returns the current settings so the control panel
// can pre-populate forms without relying solely on WebSocket push.
func (a *App) handleApiV1SettingsGET(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"reward_id": a.songRequestRewardID,
		"cmd_permissions": echo.Map{
			data.DB_KEY_CMD_PERMISSION_SR:      a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_SR],
			data.DB_KEY_CMD_PERMISSION_QUEUE:   a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_QUEUE],
			data.DB_KEY_CMD_PERMISSION_SONG:    a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_SONG],
			data.DB_KEY_CMD_PERMISSION_DELSONG: a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_DELSONG],
		},
	})
}
