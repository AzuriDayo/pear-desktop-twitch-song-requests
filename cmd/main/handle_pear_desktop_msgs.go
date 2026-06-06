package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/labstack/echo/v4"
	"github.com/nicklaw5/helix/v2"
	"github.com/valyala/fastjson"
	"golang.org/x/net/websocket"
)

func (a *App) handlePearDesktopMsgs() {
	var p fastjson.Parser
	for {
		select {
		case <-a.ctx.Done():
			return
		case msg := <-a.pearDesktopIncomingMsgs:
			v, err := p.ParseBytes(msg)
			if err != nil {
				log.Printf("Received non-json: %s\n", msg)
				continue
			}
			msgType := string(v.GetStringBytes("type"))
			switch msgType {
			case "POSITION_CHANGED":
				songQueueMutex.Lock()
				playerInfo.Position = v.GetInt("position")
				songQueueMutex.Unlock()
			case "PLAYER_INFO":
				songQueueMutex.Lock()
				playerInfo.IsPlaying = v.GetBool("isPlaying")
				playerInfo.Position = v.GetInt("position")
				songinfo := playerSonginfo{
					ImageSrc:         string(v.GetStringBytes("song", "imageSrc")),
					Artist:           string(v.GetStringBytes("song", "artist")),
					SongDuration:     v.GetInt("song", "songDuration"),
					AlternativeTitle: string(v.GetStringBytes("song", "alternativeTitle")),
					VideoId:          string(v.GetStringBytes("song", "videoId")),
				}
				playerInfo.Song = songinfo
				songQueueMutex.Unlock()
			case "VIDEO_CHANGED":
				// Phase 1: Update playerInfo and capture queue state under a brief lock.
				songQueueMutex.Lock()
				newVideoId := string(v.GetStringBytes("song", "videoId"))
				playerInfo.Position = v.GetInt("position")
				playerInfo.Song = playerSonginfo{
					ImageSrc:         string(v.GetStringBytes("song", "imageSrc")),
					Artist:           string(v.GetStringBytes("song", "artist")),
					SongDuration:     v.GetInt("song", "songDuration"),
					AlternativeTitle: string(v.GetStringBytes("song", "alternativeTitle")),
					VideoId:          newVideoId,
				}
				var queueHeadCopy SongQueueItem
				hasQueue := len(songQueue) > 0
				if hasQueue {
					queueHeadCopy = songQueue[0] // value copy so we can release the lock
				}
				songQueueMutex.Unlock()

				if !hasQueue {
					continue
				}

				// Phase 2: Determine if queue head is now playing — no lock needed (read-only copy).
				isQueueHeadPlayingNow := queueHeadCopy.Song.VideoID == newVideoId
				if !isQueueHeadPlayingNow {
					found, allVideoIDCounterparts := helpers.FindAllVideoIDCounterparts(queueHeadCopy.Song.VideoID)
					if found {
						if _, ok := allVideoIDCounterparts[newVideoId]; ok {
							isQueueHeadPlayingNow = true
						}
					}
				}

				if !isQueueHeadPlayingNow {
					log.Println("Failed to play next song in queue: " + queueHeadCopy.Song.Title + " - " + queueHeadCopy.Song.Artist)
					log.Println("Make sure it plays next in pear desktop to resume song requests!")
					continue
				}

				// Phase 3: Handle unknown song metadata — update actual queue entry under lock.
				if queueHeadCopy.Song.IsUnknown {
					queueHeadCopy.Song.Artist = string(v.GetStringBytes("song", "artist"))
					queueHeadCopy.Song.ImageUrl = string(v.GetStringBytes("song", "imageSrc"))
					queueHeadCopy.Song.Title = string(v.GetStringBytes("song", "title"))
					songQueueMutex.Lock()
					if len(songQueue) > 0 {
						songQueue[0].Song = queueHeadCopy.Song
					}
					songQueueMutex.Unlock()
					go saveSongHistory(&queueHeadCopy.Song, queueHeadCopy.RequestedBy, queueHeadCopy.RequestedByUserID, queueHeadCopy.IsNinja)
				}

				// Phase 4: Shift queue and promote the next song.
				// HTTP calls and polling happen outside the mutex to avoid holding it for seconds.
				nextInQueueIsAdded := false
				for !nextInQueueIsAdded {
					// Shift the queue under a brief lock and capture the new head.
					songQueueMutex.Lock()
					if len(songQueue) == 0 {
						songQueueMutex.Unlock()
						break
					}
					songQueue = songQueue[1:]
					var nextHeadCopy *SongQueueItem
					if len(songQueue) > 0 {
						c := songQueue[0]
						nextHeadCopy = &c
					}
					songQueueMutex.Unlock()

					// Broadcast QUEUE_SHIFT to control-panel clients.
					queueInfoOnShift, _ := json.Marshal(echo.Map{
						"type": "QUEUE_SHIFT",
					})
					a.clientsMu.Lock()
					for ws := range a.clients {
						websocket.Message.Send(ws, string(queueInfoOnShift))
					}
					a.clientsMu.Unlock()

					if nextHeadCopy == nil {
						break
					}

					// Post next song to pear desktop (no lock held).
					b := echo.Map{
						"videoId":        nextHeadCopy.Song.VideoID,
						"insertPosition": "INSERT_AFTER_CURRENT_VIDEO",
					}
					bb, _ := json.Marshal(b)
					resp, err := http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue", "application/json", bytes.NewBuffer(bb))
					if err != nil || resp.StatusCode != http.StatusNoContent {
						emsg := "Internal error when adding song to pear desktop"
						log.Println(emsg, err)
					}
					if resp != nil {
						resp.Body.Close()
					}

					// Poll to verify the song was actually added (no lock held).
					intervalDelay := time.Second
					maxRetries := 5
					for range maxRetries {
						time.Sleep(intervalDelay)
						found, _ := helpers.FindAllVideoIDCounterparts(nextHeadCopy.Song.VideoID)
						if found {
							nextInQueueIsAdded = true
							break
						}
					}

					// Notify chat if the song failed to queue, then try the next one.
					if !nextInQueueIsAdded {
						senderID := a.twitchDataStruct.userID
						if a.twitchDataStructBot.isAuthenticated {
							senderID = a.twitchDataStructBot.userID
						}
						a.helix.SendChatMessage(&helix.SendChatMessageParams{
							BroadcasterID: a.twitchDataStruct.userID,
							SenderID:      senderID,
							Message:       "Sorry " + nextHeadCopy.RequestedBy + " , I failed to add https://youtu.be/" + nextHeadCopy.Song.VideoID + " next and will be removed from queue.",
						})
					}
				}
			case "PLAYER_STATE_CHANGED":
				songQueueMutex.Lock()
				playerInfo.Position = v.GetInt("position")
				playerInfo.IsPlaying = v.GetBool("isPlaying")
				songQueueMutex.Unlock()
			default:
				// Nothing, ignore non important
			}
		}
	}
}
