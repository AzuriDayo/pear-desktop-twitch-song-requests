package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/gen/model"
	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/go-jet/jet/v2/sqlite"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/nicklaw5/helix/v2"
)

// TwitchInfoPayload mirrors the old "TWITCH_INFO" websocket message and is used
// both as the return type of GetTwitchInfo and the payload of the TWITCH_INFO event.
type TwitchInfoPayload struct {
	Type               string `json:"type"`
	StreamOnline       bool   `json:"stream_online"`
	RewardID           string `json:"reward_id"`
	Login              string `json:"login"`
	LoginBot           string `json:"login_bot"`
	ExpiryDate         string `json:"expiry_date"`
	ExpiryDateBot      string `json:"expiry_date_bot"`
	RefreshExpiryDate    string `json:"refresh_expiry_date"`
	RefreshExpiryDateBot string `json:"refresh_expiry_date_bot"`
}

// twitchInfoPayload builds the current Twitch info snapshot.
func (a *App) twitchInfoPayload() TwitchInfoPayload {
	// Access-token expiry is only relevant for legacy implicit-grant tokens (no refresh token).
	expiryDate := ""
	if a.twitchDataStruct.isAuthenticated && a.twitchDataStruct.refreshToken == "" {
		expiryDate = a.twitchDataStruct.expiresDate.Local().Format(data.TWITCH_SERVER_DATE_LAYOUT)
	}
	expiryDateBot := ""
	if a.twitchDataStructBot.isAuthenticated && a.twitchDataStructBot.refreshToken == "" {
		expiryDateBot = a.twitchDataStructBot.expiresDate.Local().Format(data.TWITCH_SERVER_DATE_LAYOUT)
	}
	return TwitchInfoPayload{
		Type:                 "TWITCH_INFO",
		StreamOnline:         a.streamOnline,
		RewardID:             a.songRequestRewardID,
		Login:                a.twitchDataStruct.login,
		LoginBot:             a.twitchDataStructBot.login,
		ExpiryDate:           expiryDate,
		ExpiryDateBot:        expiryDateBot,
		RefreshExpiryDate:    formatRefreshTokenExpiryDate(a.twitchDataStruct),
		RefreshExpiryDateBot: formatRefreshTokenExpiryDate(a.twitchDataStructBot),
	}
}

// GetTwitchInfo returns the current Twitch info snapshot so the frontend can seed
// its state on mount (replacing the initial TWITCH_INFO websocket push).
func (a *App) GetTwitchInfo() TwitchInfoPayload {
	return a.twitchInfoPayload()
}

// GetQueue returns the current internal song queue so the frontend can seed its
// state on mount (replacing the initial QUEUE_INFO websocket push).
func (a *App) GetQueue() []SongQueueItem {
	songQueueMutex.RLock()
	defer songQueueMutex.RUnlock()
	out := make([]SongQueueItem, len(songQueue))
	copy(out, songQueue)
	return out
}

// SettingsResponse mirrors the old GET /api/v1/settings response body.
type SettingsResponse struct {
	RewardID       string         `json:"reward_id"`
	CmdPermissions map[string]int `json:"cmd_permissions"`
}

// GetSettings returns the current reward id and command permission levels.
func (a *App) GetSettings() SettingsResponse {
	return SettingsResponse{
		RewardID: a.songRequestRewardID,
		CmdPermissions: map[string]int{
			data.DB_KEY_CMD_PERMISSION_SR:      a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_SR],
			data.DB_KEY_CMD_PERMISSION_QUEUE:   a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_QUEUE],
			data.DB_KEY_CMD_PERMISSION_SONG:    a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_SONG],
			data.DB_KEY_CMD_PERMISSION_DELSONG: a.cmdPermissions[data.DB_KEY_CMD_PERMISSION_DELSONG],
		},
	}
}

// parsePermissionLevel parses a string to a valid permission level integer (0-4).
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

// SaveSettings persists the given settings (reward id and/or command permissions)
// and emits an updated TWITCH_INFO event. It replaces PATCH /api/v1/settings.
func (a *App) SaveSettings(settings map[string]string) error {
	db, err := databaseconn.NewDBConnection()
	if err != nil {
		return errors.New("save data failed")
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
			_, err = stmt.ExecContext(a.ctx, db)
			if err != nil {
				log.Println("SaveSettings: failed to save setting", k, err)
			}
			a.songRequestRewardID = v
		}

		if permissionKeys[k] {
			level, err := parsePermissionLevel(v)
			if err != nil {
				log.Println("SaveSettings: invalid permission level for", k, err)
				continue
			}
			newSetting := model.Settings{Key: k, Value: v}
			stmt := Settings.INSERT(Settings.AllColumns).MODEL(newSetting).ON_CONFLICT(Settings.Key).DO_UPDATE(SET(
				Settings.Value.SET(String(v)),
			))
			_, err = stmt.ExecContext(a.ctx, db)
			if err != nil {
				log.Println("SaveSettings: failed to save permission setting", k, err)
				continue
			}
			a.cmdPermissions[k] = level
		}
	}

	a.emit("TWITCH_INFO", a.twitchInfoPayload())
	return nil
}

// SelectableReward is a channel-points reward that can be selected as the song request trigger.
type SelectableReward struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Cost int    `json:"cost"`
}

// GetTwitchCustomRewards lists the channel's custom rewards that require user input.
func (a *App) GetTwitchCustomRewards() ([]SelectableReward, error) {
	if a.helix == nil {
		return nil, errors.New("not authenticated")
	}
	rewards, err := a.helix.GetCustomRewards(&helix.GetCustomRewardsParams{
		BroadcasterID: a.twitchDataStruct.userID,
	})
	if err != nil {
		e := "Failed to list custom rewards from channel " + a.twitchDataStruct.login
		log.Println(e)
		return nil, errors.New(e)
	}
	selectableRewards := []SelectableReward{}
	for _, v := range rewards.Data.ChannelCustomRewards {
		if v.IsUserInputRequired {
			selectableRewards = append(selectableRewards, SelectableReward{
				ID:   v.ID,
				Name: v.Title,
				Cost: v.Cost,
			})
		}
	}
	return selectableRewards, nil
}

// HistoryItem is a single requester history row.
type HistoryItem struct {
	ID             int64  `json:"id"`
	VideoID        string `json:"video_id"`
	TwitchUsername string `json:"twitch_username"`
	RequestedAt    string `json:"requested_at"`
	IsNinja        bool   `json:"is_ninja"`
	SongTitle      string `json:"song_title"`
	ArtistName     string `json:"artist_name"`
	ImageURL       string `json:"image_url"`
}

// HistoryResponse is the paginated requester history response.
type HistoryResponse struct {
	MaxResults int64         `json:"max_results"`
	Items      []HistoryItem `json:"items"`
}

// GetRequestersHistory returns a page of the requester history (page is 0-based,
// perPage must be between 10 and 100). It replaces GET /api/v1/requesters/history.
func (a *App) GetRequestersHistory(page int, perPage int) (HistoryResponse, error) {
	if page < 0 {
		return HistoryResponse{}, errors.New("page must be >= 0")
	}
	if perPage < 10 || perPage > 100 {
		return HistoryResponse{}, errors.New("perPage must be between 10 and 100")
	}

	db, err := databaseconn.NewDBConnection()
	if err != nil {
		log.Println("GetRequestersHistory: failed to open database connection")
		return HistoryResponse{}, errors.New("internal error")
	}
	defer db.Close()

	maxResults := struct {
		MaxResults int64
	}{}
	queryStmt := SELECT(MAX(sqlite.IntegerColumn("rowid")).AS("max_results")).FROM(SongRequestRequesters)
	err = queryStmt.QueryContext(a.ctx, db, &maxResults)
	if err != nil {
		log.Println("GetRequestersHistory: failed to query max page")
		return HistoryResponse{}, errors.New("internal error")
	}
	maxPage := maxResults.MaxResults / int64(perPage)
	if maxResults.MaxResults%int64(perPage) > 0 {
		maxPage++
	}

	if int64(page) > maxPage {
		log.Println("GetRequestersHistory: page > max_page")
		return HistoryResponse{}, errors.New("page out of range")
	}

	results := []HistoryItem{}
	queryStmt = SELECT(sqlite.IntegerColumn("rowid").AS("id"), SongRequestRequesters.VideoID.AS("video_id"), SongRequestRequesters.TwitchUsername.AS("twitch_username"), SongRequestRequesters.RequestedAt.AS("requested_at"), SongRequestRequesters.IsNinja.AS("is_ninja"), SongRequests.SongTitle.AS("song_title"), SongRequests.ArtistName.AS("artist_name"), SongRequests.ImageURL.AS("image_url")).FROM(SongRequestRequesters.LEFT_JOIN(SongRequests, SongRequests.VideoID.EQ(SongRequestRequesters.VideoID))).ORDER_BY(sqlite.IntegerColumn("rowid").DESC()).LIMIT(int64(perPage)).OFFSET(int64(page) * int64(perPage))
	err = queryStmt.QueryContext(a.ctx, db, &results)
	if err != nil {
		log.Println("GetRequestersHistory: failed to query data", err)
		return HistoryResponse{}, errors.New("internal error")
	}

	for i := range results {
		t, _ := time.Parse(data.TWITCH_SERVER_DATE_LAYOUT, results[i].RequestedAt)
		results[i].RequestedAt = t.UTC().Format(time.RFC3339)
	}

	return HistoryResponse{
		MaxResults: maxResults.MaxResults,
		Items:      results,
	}, nil
}

// DeleteQueueItem removes a song from the internal queue by 1-based index (matching
// !delsong # semantics). If the head is removed it also cleans up the pear-desktop
// player queue. It emits QUEUE_INFO afterwards. Replaces DELETE /api/v1/queue/:idx.
func (a *App) DeleteQueueItem(idx int) error {
	if idx < 1 {
		return errors.New("idx must be a positive integer")
	}
	idx-- // convert to 0-based

	songQueueMutex.Lock()
	defer songQueueMutex.Unlock()

	if idx >= len(songQueue) {
		return errors.New("idx out of range")
	}

	song := songQueue[idx]

	// Remove from internal queue
	songQueue = append(songQueue[:idx], songQueue[idx+1:]...)

	// If the first item was removed, replicate !delsong pear-desktop cleanup.
	if idx == 0 {
		intervalDelay := time.Second
		maxRetries := 3
		pearIndex := -1
		found := false
		for range maxRetries {
			time.Sleep(intervalDelay)
			found2, videoData := helpers.FindAllVideoIDCounterparts(song.Song.VideoID)
			found = found2
			if found2 {
				pearIndex = videoData[song.Song.VideoID].Index
				break
			}
		}
		if found {
			req, _ := http.NewRequest(http.MethodDelete, "http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue/"+strconv.Itoa(pearIndex), nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusNoContent {
				log.Println("delsong API cleanup: Failed to delete song from pear-desktop, proceeding anyway...")
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
		if len(songQueue) > 0 {
			b := map[string]any{
				"videoId":        songQueue[0].Song.VideoID,
				"insertPosition": "INSERT_AFTER_CURRENT_VIDEO",
			}
			bb, _ := json.Marshal(b)
			resp, err := http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue", "application/json", bytes.NewBuffer(bb))
			if err != nil || resp.StatusCode != http.StatusNoContent {
				log.Println("delsong API cleanup: Failed to add next song in queue to pear-desktop. https://youtu.be/" + songQueue[0].Song.VideoID)
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}

	log.Printf("Control panel: removed song #%d %s - %s from queue\n", idx+1, song.Song.Title, song.Song.Artist)

	a.emit("QUEUE_INFO", QueueInfoPayload{SongQueue: append([]SongQueueItem{}, songQueue...)})
	return nil
}
