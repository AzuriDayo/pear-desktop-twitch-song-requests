package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// handleApiV1QueueDeleteDELETE removes a song from the internal song queue by
// 1-based index (matching !delsong # semantics). It replicates the full
// !delsong cleanup logic: if the first item is removed it also removes that
// song from the pear-desktop player queue and inserts the new first item.
// After mutating the queue it broadcasts QUEUE_INFO to all WS clients.
//
// DELETE /api/v1/queue/:idx   (idx is 1-based)
func (a *App) handleApiV1QueueDeleteDELETE(c echo.Context) error {
	idxStr := c.Param("idx")
	idx, err := strconv.Atoi(idxStr)
	if err != nil || idx < 1 {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "idx must be a positive integer"})
	}
	idx-- // convert to 0-based

	songQueueMutex.Lock()
	defer songQueueMutex.Unlock()

	if idx >= len(songQueue) {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "idx out of range"})
	}

	song := songQueue[idx]

	// Remove from internal queue
	songQueue = append(songQueue[:idx], songQueue[idx+1:]...)

	// If the first item was removed, replicate !delsong pear-desktop cleanup
	if idx == 0 {
		// Try to find and remove the song from pear-desktop player queue
		intervalDelay := time.Second
		maxRetries := 3
		pearIndex := -1
		found := false
		for range maxRetries {
			time.Sleep(intervalDelay)
			found2, videoData := helpers.FindAllVideoIDCounterparts(song.Song.VideoID)
			found = found2
			pearIndex = videoData[song.Song.VideoID].Index
			break
		}
		if found {
			req, _ := http.NewRequest(http.MethodDelete, "http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue/"+strconv.Itoa(pearIndex), nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusNoContent {
				log.Println("delsong API cleanup: Failed to delete song from pear-desktop, proceeding anyway...")
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
		}
	}

	log.Printf("Control panel: removed song #%d %s - %s from queue\n", idx+1, song.Song.Title, song.Song.Artist)

	// Broadcast updated queue to all WS clients
	queueInfo, _ := json.Marshal(echo.Map{
		"type":       "QUEUE_INFO",
		"song_queue": songQueue,
	})
	a.clientsMu.RLock()
	for ws := range a.clients {
		websocket.Message.Send(ws, string(queueInfo))
	}
	a.clientsMu.RUnlock()

	return c.NoContent(http.StatusNoContent)
}
