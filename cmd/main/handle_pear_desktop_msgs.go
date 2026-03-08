package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/labstack/echo/v4"
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
					if queueHead.Song.VideoID == newVideoId {
						songQueue = songQueue[1:]
						queueInfoOnShift, _ := json.Marshal(echo.Map{
							"type": "QUEUE_SHIFT",
						})
						a.clientsMu.Lock()
						for ws := range a.clients {
							websocket.Message.Send(ws, string(queueInfoOnShift))
						}
						a.clientsMu.Unlock()

						if len(songQueue) > 0 {
							b := echo.Map{
								"videoId":        songQueue[0].Song.VideoID,
								"insertPosition": "INSERT_AFTER_CURRENT_VIDEO",
							}
							bb, _ := json.Marshal(b)
							resp, err := http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/queue", "application/json", bytes.NewBuffer(bb))
							if err != nil || resp.StatusCode != http.StatusNoContent {
								emsg := "Internal error when adding song to pear desktop"
								log.Println(emsg, err)
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
