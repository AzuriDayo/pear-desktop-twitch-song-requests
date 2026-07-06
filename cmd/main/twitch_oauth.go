package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/gen/model"
	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/nicklaw5/helix/v2"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

var oauthLoginMu sync.Mutex

const (
	twitchOAuthAuthorizeURL        = "https://id.twitch.tv/oauth2/authorize"
	twitchOAuthTokenURL            = "https://id.twitch.tv/oauth2/token"
	twitchOAuthLoginTimeout        = 5 * time.Minute
	twitchRefreshTokenInactivity   = 30 * 24 * time.Hour
	twitchRefreshTokenWarnBefore   = 5 * 24 * time.Hour
	twitchAccessTokenRefreshBefore = 30 * time.Minute
	twitchTokenMaintenanceInterval = 15 * time.Minute
)

type twitchOAuthErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type twitchTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// TwitchAuthResult is emitted to the frontend when OAuth login finishes.
type TwitchAuthResult struct {
	Success bool   `json:"success"`
	ForBot  bool   `json:"for_bot"`
	Error   string `json:"error,omitempty"`
}

type oauthCallbackResult struct {
	accessToken string
	expiresIn   int
	err         error
}

const oauthCaptureHTML = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>Completing Twitch login</title></head>
<body>
<p>Completing login…</p>
<script>
(function () {
  const query = new URLSearchParams(window.location.search);
  const hash = new URLSearchParams(window.location.hash.slice(1));
  const err = hash.get("error") || query.get("error");
  if (err) {
    const desc = hash.get("error_description") || query.get("error_description") || err;
    window.location.replace("/oauth/twitch/done?error=" + encodeURIComponent(desc));
    return;
  }
  const accessToken = hash.get("access_token");
  const state = hash.get("state") || "";
  const expiresIn = hash.get("expires_in") || "0";
  if (!accessToken) {
    window.location.replace("/oauth/twitch/done?error=" + encodeURIComponent("missing access token"));
    return;
  }
  const params = new URLSearchParams({
    access_token: accessToken,
    state: state,
    expires_in: expiresIn,
  });
  window.location.replace("/oauth/twitch/done?" + params.toString());
})();
</script>
</body>
</html>`

const oauthSuccessHTML = `<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>Twitch authorization complete</title></head>
<body>
<h1>Authorization complete</h1>
<p>You can close this window and return to the app.</p>
</body>
</html>`

func twitchOAuthErrorFromBody(body []byte, statusCode int) error {
	var errResp twitchOAuthErrorResponse
	if json.Unmarshal(body, &errResp) == nil && errResp.Message != "" {
		return errors.New(errResp.Message)
	}
	return fmt.Errorf("Twitch OAuth request failed (HTTP %d)", statusCode)
}

func generateOAuthState() (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func buildTwitchImplicitAuthorizeURL(redirectURI, state string, scopes []string) string {
	params := url.Values{}
	params.Set("response_type", "token")
	params.Set("client_id", data.GetTwitchClientID())
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", strings.Join(scopes, " "))
	params.Set("state", state)
	return twitchOAuthAuthorizeURL + "?" + params.Encode()
}

func refreshTokenExpiresAt(lastUsed time.Time) time.Time {
	return lastUsed.Add(twitchRefreshTokenInactivity)
}

func formatRefreshTokenExpiryDate(td *twitchData) string {
	if td.refreshToken == "" || td.refreshTokenLastUsedAt.IsZero() {
		return ""
	}
	return refreshTokenExpiresAt(td.refreshTokenLastUsedAt).Local().Format(data.TWITCH_SERVER_DATE_LAYOUT)
}

func (a *App) markRefreshTokenUsed(forBot bool) error {
	td := twitchDataForBot(a, forBot)
	now := time.Now()
	td.refreshTokenLastUsedAt = now

	lastUsedKey := data.DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED
	if forBot {
		lastUsedKey = data.DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED_BOT
	}
	return a.saveSetting(lastUsedKey, now.UTC().Format(time.RFC3339))
}

func twitchAuthScopes(forBot bool) []string {
	scopes := []string{"user:read:chat", "user:write:chat", "user:bot", "channel:bot"}
	if !forBot {
		scopes = append(scopes,
			"channel:read:redemptions",
			"channel:read:vips",
			"moderation:read",
			"channel:read:subscriptions",
		)
	}
	return scopes
}

func helixForBot(a *App, forBot bool) *helix.Client {
	if forBot {
		return a.helixBot
	}
	return a.helix
}

func twitchDataForBot(a *App, forBot bool) *twitchData {
	if forBot {
		return a.twitchDataStructBot
	}
	return a.twitchDataStruct
}

// Login opens the system browser for Twitch implicit-grant OAuth and captures the
// access token on http://localhost:3999/oauth/twitch (no client secret required).
func (a *App) Login(forBot bool) error {
	if forBot && a.twitchDataStruct.login == "" {
		return errors.New("connect the main Twitch account first")
	}

	oauthLoginMu.Lock()
	defer oauthLoginMu.Unlock()

	err := a.runTwitchImplicitLogin(forBot)
	if err != nil {
		a.emit("TWITCH_AUTH_ERROR", TwitchAuthResult{
			Success: false,
			ForBot:  forBot,
			Error:   err.Error(),
		})
		return err
	}

	a.emit("TWITCH_AUTH_SUCCESS", TwitchAuthResult{Success: true, ForBot: forBot})
	return nil
}

func (a *App) runTwitchImplicitLogin(forBot bool) error {
	state, err := generateOAuthState()
	if err != nil {
		return fmt.Errorf("generate OAuth state: %w", err)
	}

	listener, err := net.Listen("tcp", data.TWITCH_OAUTH_LISTEN_ADDR)
	if err != nil {
		return fmt.Errorf("start OAuth listener on %s: %w (is another app using port 3999?)", data.TWITCH_OAUTH_LISTEN_ADDR, err)
	}
	defer listener.Close()

	redirectURI := data.TWITCH_OAUTH_REDIRECT_URI
	scopes := twitchAuthScopes(forBot)

	callbackCh := make(chan oauthCallbackResult, 1)
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/twitch", func(w http.ResponseWriter, r *http.Request) {
		if oauthErr := r.URL.Query().Get("error"); oauthErr != "" {
			description := r.URL.Query().Get("error_description")
			if description == "" {
				description = oauthErr
			}
			callbackCh <- oauthCallbackResult{err: errors.New(description)}
			http.Error(w, description, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(oauthCaptureHTML))
	})
	mux.HandleFunc("/oauth/twitch/done", func(w http.ResponseWriter, r *http.Request) {
		if oauthErr := r.URL.Query().Get("error"); oauthErr != "" {
			callbackCh <- oauthCallbackResult{err: errors.New(oauthErr)}
			http.Error(w, oauthErr, http.StatusBadRequest)
			return
		}

		if r.URL.Query().Get("state") != state {
			callbackCh <- oauthCallbackResult{err: errors.New("invalid OAuth state")}
			http.Error(w, "invalid OAuth state", http.StatusBadRequest)
			return
		}

		accessToken := r.URL.Query().Get("access_token")
		if accessToken == "" {
			callbackCh <- oauthCallbackResult{err: errors.New("missing access token")}
			http.Error(w, "missing access token", http.StatusBadRequest)
			return
		}

		expiresIn, _ := strconv.Atoi(r.URL.Query().Get("expires_in"))
		callbackCh <- oauthCallbackResult{accessToken: accessToken, expiresIn: expiresIn}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(oauthSuccessHTML))
	})

	srv := &http.Server{Handler: mux}
	go func() {
		if serveErr := srv.Serve(listener); serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			log.Println("OAuth loopback server error:", serveErr)
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	authURL := buildTwitchImplicitAuthorizeURL(redirectURI, state, scopes)
	if a.wailsCtx != nil {
		wailsruntime.BrowserOpenURL(a.wailsCtx, authURL)
	} else {
		return errors.New("application is not ready to open the system browser")
	}

	timeoutCtx, cancel := context.WithTimeout(a.ctx, twitchOAuthLoginTimeout)
	defer cancel()

	var callback oauthCallbackResult
	select {
	case <-timeoutCtx.Done():
		return errors.New("Twitch authorization timed out")
	case callback = <-callbackCh:
	}

	if callback.err != nil {
		return callback.err
	}

	return a.applyTwitchTokens(forBot, callback.accessToken, "", callback.expiresIn)
}

func (a *App) saveSetting(key, value string) error {
	db, err := databaseconn.NewDBConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	newSetting := model.Settings{Key: key, Value: value}
	stmt := Settings.INSERT(Settings.AllColumns).MODEL(newSetting).ON_CONFLICT(Settings.Key).DO_UPDATE(SET(
		Settings.Value.SET(String(value)),
	))
	_, err = stmt.ExecContext(a.ctx, db)
	return err
}

func (a *App) applyTwitchTokens(forBot bool, accessToken, refreshToken string, expiresInSec int) error {
	client := helixForBot(a, forBot)
	td := twitchDataForBot(a, forBot)

	if forBot {
		if a.twitchDataStruct.login == "" {
			return errors.New("cannot save bot token before main token")
		}
	}

	client.SetUserAccessToken(accessToken)
	isValid, validateResp, err := client.ValidateToken(accessToken)
	if err != nil {
		return fmt.Errorf("validate token: %w", err)
	}
	if validateResp.StatusCode != http.StatusOK || !isValid {
		return errors.New("Twitch rejected the access token")
	}
	if forBot && strings.EqualFold(a.twitchDataStruct.login, validateResp.Data.Login) {
		return errors.New("bot token is same as main token")
	}

	expiresDate := time.Now().Add(time.Duration(expiresInSec) * time.Second)
	if expiresInSec <= 0 {
		strDate := validateResp.Header.Get("Date")
		t, parseErr := time.Parse(data.TWITCH_SERVER_DATE_LAYOUT, strDate)
		if parseErr != nil {
			return errors.New("incorrect expected date data from Twitch")
		}
		expiresDate = t.Add(time.Duration(validateResp.Data.ExpiresIn) * time.Second)
	}

	accessKey := data.DB_KEY_TWITCH_ACCESS_TOKEN
	refreshKey := data.DB_KEY_TWITCH_REFRESH_TOKEN
	if forBot {
		accessKey = data.DB_KEY_TWITCH_ACCESS_TOKEN_BOT
		refreshKey = data.DB_KEY_TWITCH_REFRESH_TOKEN_BOT
	}

	if err := a.saveSetting(accessKey, accessToken); err != nil {
		return errors.New("failed to save access token")
	}
	if refreshToken != "" {
		if err := a.saveSetting(refreshKey, refreshToken); err != nil {
			return errors.New("failed to save refresh token")
		}
		if err := a.markRefreshTokenUsed(forBot); err != nil {
			return errors.New("failed to save refresh token last used time")
		}
		td.refreshToken = refreshToken
	} else {
		td.refreshToken = ""
		td.refreshTokenLastUsedAt = time.Time{}
		if err := a.saveSetting(refreshKey, ""); err != nil {
			return errors.New("failed to clear refresh token")
		}
		lastUsedKey := data.DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED
		if forBot {
			lastUsedKey = data.DB_KEY_TWITCH_REFRESH_TOKEN_LAST_USED_BOT
		}
		if err := a.saveSetting(lastUsedKey, ""); err != nil {
			return errors.New("failed to clear refresh token last used time")
		}
	}

	client.SetUserAccessToken(accessToken)
	td.accessToken = accessToken
	td.expiresDate = expiresDate
	td.isAuthenticated = true
	td.userID = validateResp.Data.UserID
	td.login = validateResp.Data.Login

	if !forBot {
		streamResp, err := client.GetStreams(&helix.StreamsParams{
			UserLogins: []string{a.twitchDataStruct.login},
		})
		if err == nil && len(streamResp.Data.Streams) > 0 && streamResp.Data.Streams[0].ID != "" {
			a.streamOnline = true
		}
	}

	a.emit("TWITCH_INFO", a.twitchInfoPayload())
	return nil
}

func (a *App) refreshTwitchAccessToken(refreshToken string) (accessToken, newRefreshToken string, expiresIn int, err error) {
	form := url.Values{}
	form.Set("client_id", data.GetTwitchClientID())
	form.Set("grant_type", "refresh_token")
	form.Set("refresh_token", refreshToken)

	resp, err := http.PostForm(twitchOAuthTokenURL, form)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", 0, twitchOAuthErrorFromBody(body, resp.StatusCode)
	}

	var token twitchTokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return "", "", 0, err
	}
	return token.AccessToken, token.RefreshToken, token.ExpiresIn, nil
}

func (a *App) tryRefreshTwitchToken(forBot bool) error {
	td := twitchDataForBot(a, forBot)
	if td.refreshToken == "" {
		return errors.New("no refresh token stored")
	}

	accessToken, refreshToken, expiresIn, err := a.refreshTwitchAccessToken(td.refreshToken)
	if err != nil {
		return err
	}
	return a.applyTwitchTokens(forBot, accessToken, refreshToken, expiresIn)
}

func (a *App) validateLoadedTwitchToken(forBot bool) error {
	client := helixForBot(a, forBot)
	td := twitchDataForBot(a, forBot)
	if td.accessToken == "" {
		if td.refreshToken != "" {
			log.Printf("Twitch %s access token missing, refreshing...", map[bool]string{true: "bot", false: "main"}[forBot])
			return a.tryRefreshTwitchToken(forBot)
		}
		return nil
	}

	isValid, response, err := client.ValidateToken(td.accessToken)
	if err != nil {
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
		client.SetUserAccessToken(td.accessToken)
		td.expiresDate = t
		td.isAuthenticated = true
		td.userID = response.Data.UserID
		td.login = response.Data.Login

		if !forBot {
			resp, err := client.GetStreams(&helix.StreamsParams{
				UserLogins: []string{a.twitchDataStruct.login},
			})
			if err == nil && len(resp.Data.Streams) > 0 && resp.Data.Streams[0].ID != "" {
				a.streamOnline = true
			}
		}
		return nil
	}

	if td.refreshToken != "" {
		log.Printf("Twitch %s access token expired, refreshing...", map[bool]string{true: "bot", false: "main"}[forBot])
		return a.tryRefreshTwitchToken(forBot)
	}

	td.isAuthenticated = false
	return nil
}

func (a *App) runTwitchTokenMaintenance() {
	ticker := time.NewTicker(twitchTokenMaintenanceInterval)
	defer ticker.Stop()

	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			a.maintainTwitchTokens(false)
			a.maintainTwitchTokens(true)
		}
	}
}

func (a *App) notifyImplicitGrantAccessTokenExpiresSoon(forBot bool) {
	td := twitchDataForBot(a, forBot)
	if td.refreshToken != "" {
		return
	}
	account := "main"
	if forBot {
		account = "bot"
	}
	notifyTokenExpiresSoon(account, td.expiresDate, td.isAuthenticated)
}

func (a *App) maintainTwitchTokens(forBot bool) {
	td := twitchDataForBot(a, forBot)
	if td.refreshToken == "" {
		return
	}

	account := "main"
	if forBot {
		account = "bot"
	}

	if td.isAuthenticated && time.Until(td.expiresDate) < twitchAccessTokenRefreshBefore {
		if err := a.tryRefreshTwitchToken(forBot); err != nil {
			log.Printf("Twitch %s proactive token refresh failed: %v", account, err)
		}
	}

	notifyRefreshTokenExpiresSoon(account, td.refreshTokenLastUsedAt, td.refreshToken)
}
