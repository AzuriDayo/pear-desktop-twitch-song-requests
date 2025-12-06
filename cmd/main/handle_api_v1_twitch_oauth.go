package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/labstack/echo/v4"
)

func (a *App) processTwitchOAuth(c echo.Context) error {
	// auth data in url hash string params as get request
	body := c.Request().Body
	rawBodyData, err := io.ReadAll(body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "read request body",
		})
	}
	defer body.Close()

	authData := struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		State       string `json:"state,omitempty"`
		TokenType   string `json:"token_type"`
	}{}
	err = json.Unmarshal(rawBodyData, &authData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "parse request body",
		})
	}

	if authData.TokenType != "bearer" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unexpected token type",
		})
	}

	isValid, response, err := a.helix.ValidateToken(authData.AccessToken)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusServiceUnavailable, echo.Map{
			"error": "Twitch authentication validation failed, see console.",
		})
	}
	if response.StatusCode == http.StatusOK && isValid {
		expiresIn := response.Data.ExpiresIn
		strDate := response.Header.Get("Date")
		t, err := time.Parse(data.TWITCH_SERVER_DATE_LAYOUT, strDate)
		if err != nil {
			c.Logger().Error(errors.New("Failed to validate server date time expiry, original error:\n" + err.Error()))
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"error": "incorrect expected date data from Twitch",
			})
		}
		t = t.Add(time.Duration(expiresIn) * time.Second)
		a.helix.SetUserAccessToken(authData.AccessToken)
		a.twitchDataStruct.expiresDate = t
		a.twitchDataStruct.isAuthenticated = true
		a.twitchDataStruct.userID = response.Data.UserID
		a.twitchDataStruct.login = response.Data.Login

		return c.NoContent(http.StatusOK)
	} else {
		return c.NoContent(http.StatusUnauthorized)
	}
}
