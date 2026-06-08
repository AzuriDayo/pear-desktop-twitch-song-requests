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
				// songQueueMutex is held for the entire case: we must not allow concurrent
				// queue mutations (songRequestLogic, !delsong, API delete) to interleave
				// while we coordinate pear-desktop player state with the internal queue.
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

				if len(songQueue) == 0 {
					songQueueMutex.Unlock()
					continue
				}

				// Determine if the queue head is now playing.
				isQueueHeadPlayingNow := songQueue[0].Song.VideoID == newVideoId
				if !isQueueHeadPlayingNow {
					found, allVideoIDCounterparts := helpers.FindAllVideoIDCounterparts(songQueue[0].Song.VideoID)
					if found {
						if _, ok := allVideoIDCounterparts[newVideoId]; ok {
							isQueueHeadPlayingNow = true
						}
					}
				}

				if !isQueueHeadPlayingNow {
					log.Println("Failed to play next song in queue: " + songQueue[0].Song.Title + " - " + songQueue[0].Song.Artist)
					log.Println("Make sure it plays next in pear desktop to resume song requests!")
					songQueueMutex.Unlock()
					continue
				}

				// Handle unknown song metadata: fill in live data from pear-desktop.
				if songQueue[0].Song.IsUnknown {
					songQueue[0].Song.Artist = string(v.GetStringBytes("song", "artist"))
					songQueue[0].Song.ImageUrl = string(v.GetStringBytes("song", "imageSrc"))
					songQueue[0].Song.Title = string(v.GetStringBytes("song", "title"))
					// Value-copy the head before the queue shifts so the goroutine
					// doesn't alias into the slice after songQueue[1:].
					headSnapshot := songQueue[0]
					go saveSongHistory(&headSnapshot.Song, headSnapshot.RequestedBy, headSnapshot.RequestedByUserID, headSnapshot.IsNinja)
				}

				// Shift queue and promote the next song.
				// Lock is held through HTTP calls and polling to prevent any concurrent
				// mutation from observing or modifying the queue mid-transition.
				nextInQueueIsAdded := false
				for !nextInQueueIsAdded {
					if len(songQueue) == 0 {
						break
					}
					songQueue = songQueue[1:]

					// Broadcast QUEUE_SHIFT to control-panel clients.
					// clientsMu is safe to acquire here: no code path holds clientsMu
					// and then tries to acquire songQueueMutex (verified in ws handler).
					queueInfoOnShift, _ := json.Marshal(echo.Map{
						"type": "QUEUE_SHIFT",
					})
					a.clientsMu.Lock()
					for ws := range a.clients {
						websocket.Message.Send(ws, string(queueInfoOnShift))
					}
					a.clientsMu.Unlock()

					if len(songQueue) == 0 {
						break
					}

					// Post next song to pear desktop.
					b := echo.Map{
						"videoId":        songQueue[0].Song.VideoID,
						"insertPosition": "INSERT_AFTER_CURRENT_VIDEO",
					}
					bb, _ := json.Marshal(b)
					resp, err := http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue", "application/json", bytes.NewBuffer(bb))
					if err != nil || resp.StatusCode != http.StatusNoContent {
						log.Println("Internal error when adding song to pear desktop", err)
					}
					if resp != nil {
						resp.Body.Close()
					}

					// Poll to verify the song was actually added.
					intervalDelay := time.Second
					maxRetries := 5
					for range maxRetries {
						time.Sleep(intervalDelay)
						found, _ := helpers.FindAllVideoIDCounterparts(songQueue[0].Song.VideoID)
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
							Message:       "Sorry " + songQueue[0].RequestedBy + " , I failed to add https://youtu.be/" + songQueue[0].Song.VideoID + " next and will be removed from queue.",
						})
					}
				}
				songQueueMutex.Unlock()
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
