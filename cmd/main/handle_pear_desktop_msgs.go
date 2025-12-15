package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/valyala/fastjson"
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
				if playerInfo.Song.VideoId != newVideoId {
					songinfo := playerSonginfo{
						ImageSrc:         string(v.GetStringBytes("song", "imageSrc")),
						Artist:           string(v.GetStringBytes("song", "artist")),
						SongDuration:     v.GetInt("song", "songDuration"),
						AlternativeTitle: string(v.GetStringBytes("song", "alternativeTitle")),
						VideoId:          newVideoId,
					}
					playerInfo.Song = songinfo
					if len(songQueue) > 1 && songQueue[0].song.VideoID != newVideoId {
						// queue invalid now, wiping queue
						log.Println("App internal queue order and ytm queue order mismatched, attempting to recover queue...")
						recoverVideoId := songQueue[len(songQueue)-1].song.VideoID
						queue := songrequests.QueueResponse{
							Items: []struct {
								PlaylistPanelVideoRenderer        *songrequests.QueueResponsePlaylistPanelVideoRenderer `json:"playlistPanelVideoRenderer"`
								PlaylistPanelVideoWrapperRenderer *struct {
									PrimaryRenderer struct {
										PlaylistPanelVideoRenderer songrequests.QueueResponsePlaylistPanelVideoRenderer `json:"playlistPanelVideoRenderer"`
									} `json:"primaryRenderer"`
								} `json:"playlistPanelVideoWrapperRenderer"`
							}{},
						}
						failed := false
						resp, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/queue")
						if err != nil || resp.StatusCode != http.StatusOK {
							// failed recovery
							failed = true
						}
						if !failed {
							qb, err := io.ReadAll(resp.Body)
							if err != nil {
								failed = true
							}
							if !failed {
								resp.Body.Close()
								err = json.Unmarshal(qb, &queue)
								if err != nil {
									failed = true
								}
							}
						}

						// TODO: need to change this loop to start from the back, since songs might have been re-requested
						fromId := -1
						toId := -1
						for i := len(queue.Items) - 1; i >= 0; i-- {
							if queue.Items[i].PlaylistPanelVideoWrapperRenderer != nil {
								queue.Items[i].PlaylistPanelVideoRenderer = &queue.Items[i].PlaylistPanelVideoWrapperRenderer.PrimaryRenderer.PlaylistPanelVideoRenderer
							}
							if queue.Items[i].PlaylistPanelVideoRenderer.Selected {
								fromId = i
							}
							if queue.Items[i].PlaylistPanelVideoRenderer.VideoId == recoverVideoId {
								toId = i
							}
							if fromId != -1 && toId != -1 {
								break
							}
						}
						if fromId == -1 || toId == -1 {
							failed = true
						}
						if !failed {
							songQueue = []struct {
								requestedBy string
								song        songrequests.SongResult
							}{}
							for i := fromId; i <= toId; i++ {
								if queue.Items[i].PlaylistPanelVideoRenderer.VideoId != newVideoId {
									songQueue = append(songQueue, struct {
										requestedBy string
										song        songrequests.SongResult
									}{
										requestedBy: "recovered",
										song: songrequests.SongResult{
											Title:   queue.Items[i].PlaylistPanelVideoRenderer.Title.Runs[0].Text,
											Artist:  queue.Items[i].PlaylistPanelVideoRenderer.ShortBylineText.Runs[0].Text,
											VideoID: queue.Items[i].PlaylistPanelVideoRenderer.VideoId,
										},
									})
								}
							}
							// recover success
							log.Println("Recovery successful, internal queue maintained")
						} else {
							log.Println("Recovery unsuccessful, internal queue is wiped, but your queue in ytm is intact")
						}
					} else {
						if len(songQueue) > 0 {
							songQueue = songQueue[1:]
						}
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
