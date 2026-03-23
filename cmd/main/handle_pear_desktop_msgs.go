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
				songQueueMutex.Lock()
				newVideoId := string(v.GetStringBytes("song", "videoId"))
				playerInfo.Position = v.GetInt("position")

				songinfo := playerSonginfo{
					ImageSrc:         string(v.GetStringBytes("song", "imageSrc")),
					Artist:           string(v.GetStringBytes("song", "artist")),
					SongDuration:     v.GetInt("song", "songDuration"),
					AlternativeTitle: string(v.GetStringBytes("song", "alternativeTitle")),
					VideoId:          newVideoId,
				}
				playerInfo.Song = songinfo

				if len(songQueue) > 0 {
					queueHead := songQueue[0]
					isQueueHeadPlayingNow := queueHead.Song.VideoID == newVideoId
					if !isQueueHeadPlayingNow {
						found, allVideoIDCounterparts := helpers.FindAllVideoIDCounterparts(queueHead.Song.VideoID)
						if found {
							if _, ok := allVideoIDCounterparts[newVideoId]; ok {
								isQueueHeadPlayingNow = true
							}
						}
					}

					if isQueueHeadPlayingNow {
						// handle unknown with this new song detail
						if queueHead.Song.IsUnknown {
							// re-fill song details
							queueHead.Song.Artist = string(v.GetStringBytes("song", "artist"))
							queueHead.Song.ImageUrl = string(v.GetStringBytes("song", "imageSrc"))
							queueHead.Song.Title = string(v.GetStringBytes("song", "title"))
							// save history
							go saveSongHistory(&queueHead.Song, queueHead.RequestedBy, queueHead.RequestedByUserID, queueHead.IsNinja)
							// Send frontend for unknown vid update
							// No info to send because it is already shifted from queue
						}

						nextInQueueIsAdded := false
						for !nextInQueueIsAdded && len(songQueue) > 0 {
							songQueue = songQueue[1:]
							queueInfoOnShift, _ := json.Marshal(echo.Map{
								"type": "QUEUE_SHIFT",
							})
							a.clientsMu.Lock()
							for ws := range a.clients {
								websocket.Message.Send(ws, string(queueInfoOnShift))
							}
							a.clientsMu.Unlock()

							// TODO: find out why it wont add song to queue mid-queue for "problematic song"
							if len(songQueue) > 0 {
								newQueueHead := songQueue[0]
								b := echo.Map{
									"videoId":        newQueueHead.Song.VideoID,
									"insertPosition": "INSERT_AFTER_CURRENT_VIDEO",
								}
								bb, _ := json.Marshal(b)
								resp, err := http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue", "application/json", bytes.NewBuffer(bb))
								if err != nil || resp.StatusCode != http.StatusNoContent {
									emsg := "Internal error when adding song to pear desktop"
									log.Println(emsg, err)
								}

								// validate if song was really added
								intervalDelay := time.Second

								maxRetries := 5
								for range maxRetries {
									time.Sleep(intervalDelay)
									found, _ := helpers.FindAllVideoIDCounterparts(newQueueHead.Song.VideoID)
									if found {
										nextInQueueIsAdded = true
										break
									}
								}

								// Notify failed
								if !nextInQueueIsAdded {
									senderID := a.twitchDataStruct.userID
									if a.twitchDataStructBot.isAuthenticated {
										senderID = a.twitchDataStructBot.userID
									}
									a.helix.SendChatMessage(&helix.SendChatMessageParams{
										BroadcasterID: a.twitchDataStruct.userID,
										SenderID:      senderID,
										Message:       "Sorry " + newQueueHead.RequestedBy + " , I failed to add https://youtu.be/" + newQueueHead.Song.VideoID + " next and will be removed from queue.",
									})
								}

								// Failed to add this song, we go for the one after this one
							}
						}
					} else {
						log.Println("Failed to play next song in queue: " + queueHead.Song.Title + " - " + queueHead.Song.Artist)
						log.Println("Make sure it plays next in pear desktop to resume song requests!")
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
