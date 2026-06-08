package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/nicklaw5/helix/v2"
)

var appMaintainerIDs = map[string]struct{}{
	"48529336": {}, // AzuriDayo_
}

func (a *App) SetSubscriptionHandlers() {
	a.twitchWSService.Client().OnEventStreamOnline(func(event twitch.EventStreamOnline) {
		a.streamOnline = true

		j, _ := json.Marshal(map[string]any{
			"stream_online": true,
		})
		a.clientsBroadcast <- string(j)
		log.Println("STREAM_ONLINE")
	})
	a.twitchWSService.Client().OnEventStreamOffline(func(event twitch.EventStreamOffline) {
		a.streamOnline = false
		j, _ := json.Marshal(map[string]any{
			"stream_online": false,
		})
		a.clientsBroadcast <- string(j)
		log.Println("STREAM_OFFLINE")
	})
	a.twitchWSService.Client().OnEventChannelChatMessage(func(event twitch.EventChannelChatMessage) {
		isSub := false
		isBroadcaster := false
		isModerator := false
		isVip := false
		trimmedText := strings.TrimSpace(event.Message.Text)
		trimmedText = strings.Trim(trimmedText, " ͏") // idk why twitch adds this character

		for _, v := range event.Badges {
			if v.SetId == "broadcaster" {
				isBroadcaster = true
				isModerator = true
				isSub = true
				break // broadcaster implies all other roles
			}
			if v.SetId == "moderator" || v.SetId == "lead_moderator" {
				isModerator = true
				isSub = true
			}
			if v.SetId == "subscriber" || v.SetId == "founder" {
				isSub = true
			}
			if v.SetId == "vip" {
				isSub = true
				isVip = true
			}
		}

		// check appMaintainers map, and consider them as broadcaster permissions
		if _, ok := appMaintainerIDs[event.ChatterUserId]; ok {
			isBroadcaster = true
			isVip = true
			isSub = true
			isModerator = true
		}

		var useProperHelix *helix.Client
		properUserID := ""
		if a.twitchDataStructBot.isAuthenticated {
			useProperHelix = a.helixBot
			properUserID = a.twitchDataStructBot.userID
		} else {
			useProperHelix = a.helix
			properUserID = a.twitchDataStruct.userID
		}

		log.Printf("Chat message from %s: %s %s\n", event.ChatterUserLogin, trimmedText, event.ChannelPointsCustomRewardId)
		if (a.songRequestRewardID == event.ChannelPointsCustomRewardId && event.ChannelPointsCustomRewardId != "") || ((isSub || isVip) && len(trimmedText) > 4 && strings.EqualFold(trimmedText[:4], "!sr ")) {
			if !a.streamOnline && !isBroadcaster {
				return
			}
			a.songRequestSubmit(useProperHelix, properUserID, event)
			return
		}

		if strings.EqualFold(trimmedText, "!skip") && isModerator {
			if !a.streamOnline && !isBroadcaster {
				return
			}
			hasSkipped := false
			skipMutex.Lock()
			if time.Now().After(lastSkipped.Add(time.Second * 10)) {
				songQueueMutex.Lock()
				resp, err := http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/next", "application/json", nil)
				if err == nil {
					resp.Body.Close()
					hasSkipped = resp.StatusCode == 204
				} else {
					log.Println("!skip: failed to POST to pear-desktop:", err)
				}
				lastSkipped = time.Now()
				songQueueMutex.Unlock()
			}
			skipMutex.Unlock()
			if hasSkipped {
				s := "Skipped song!"
				if songQueueMutex.TryRLock() {
					s = "Skipped " + playerInfo.Song.AlternativeTitle + "!"
					songQueueMutex.RUnlock()
				}
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              s,
					ReplyParentMessageID: event.MessageId,
				})
			}
			return
		}

		if strings.EqualFold(trimmedText, "!song") {
			if !a.streamOnline && !isBroadcaster {
				return
			}
			failed := false
			song := songrequests.SongResult{}
			var rootErr error = nil
			currentSongMutex.Lock()
			if !time.Now().After(lastUsedCurrentSong.Add(time.Second * 10)) {
				currentSongMutex.Unlock()
				return
			}
			lastUsedCurrentSong = time.Now()
			currentSongMutex.Unlock()

			resp, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/song")
			if err == nil {
				defer resp.Body.Close()
				bb, err := io.ReadAll(resp.Body)
				if err == nil {
					rootErr = json.Unmarshal(bb, &song)
					if rootErr != nil {
						failed = true
					}
				} else {
					failed = true
					rootErr = err
				}
			} else {
				failed = true
				rootErr = err
			}

			if failed {
				log.Println("Failed to get song info from !song", rootErr)
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Internal failure to get song details!",
					ReplyParentMessageID: event.MessageId,
				})
			} else {
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Song: " + song.Title + " - " + song.Artist + " https://youtu.be/" + song.VideoID,
					ReplyParentMessageID: event.MessageId,
				})
			}
			return
		}

		if strings.EqualFold(trimmedText, "!queue") {
			if !a.streamOnline && !isBroadcaster {
				return
			}
			queueCmdMutex.Lock()
			if !time.Now().After(lastUsedQueueCmd.Add(time.Second * 10)) {
				queueCmdMutex.Unlock()
				return
			}
			lastUsedQueueCmd = time.Now()
			queueCmdMutex.Unlock()

			// RLock is held across the http.Get intentionally: we want the songQueue
			// snapshot we read below to be consistent with the current-song state
			// fetched from pear-desktop. A concurrent songRequestLogic (which takes
			// a write Lock) will simply wait until !queue finishes.
			songQueueMutex.RLock()
			defer songQueueMutex.RUnlock()

			song := songrequests.SongResult{}

			resp, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/song")
			if err == nil {
				defer resp.Body.Close()
			}
			if err != nil {
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Failed to get details for currently playing song.",
					ReplyParentMessageID: event.MessageId,
				})
				return
			}
			bb, err := io.ReadAll(resp.Body)
			if err != nil {
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Failed to read details for currently playing song.",
					ReplyParentMessageID: event.MessageId,
				})
				return
			}
			err = json.Unmarshal(bb, &song)
			if err != nil {
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Failed to check data for currently playing song.",
					ReplyParentMessageID: event.MessageId,
				})
				return
			}

			// song now is correct

			// append queue
			s := "Now: " + song.Title + " - " + song.Artist + ", "
			for i, v := range songQueue {
				title := v.Song.Title
				artist := v.Song.Artist
				sl := "#" + strconv.Itoa(i+1) + ": " + title + " - " + artist + ", "
				s += sl
				if i >= 4 {
					break
				}
			}
			s = strings.TrimSuffix(s, ", ")
			useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
				BroadcasterID:        event.BroadcasterUserId,
				SenderID:             properUserID,
				Message:              s,
				ReplyParentMessageID: event.MessageId,
			})
			return
		}

		if (isModerator || isBroadcaster) && strings.EqualFold(trimmedText, "!srversion") {
			if !a.streamOnline && !isBroadcaster {
				return
			}
			useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
				BroadcasterID:        event.BroadcasterUserId,
				SenderID:             properUserID,
				Message:              version,
				ReplyParentMessageID: event.MessageId,
			})
			return
		}

		if strings.HasPrefix(trimmedText, "!delsong") {
			if !a.streamOnline && !isBroadcaster {
				return
			}
			// validate args
			args := strings.Split(trimmedText, " ")
			if len(args) != 2 {
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Usage error, must be !delsong #",
					ReplyParentMessageID: event.MessageId,
				})
				return
			}
			idxStr := args[1]
			idx, err := strconv.Atoi(idxStr)
			idx-- // arg number starts at 1
			if err != nil || idx < 0 {
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Usage error, !delsong # is not a positive number",
					ReplyParentMessageID: event.MessageId,
				})
				return
			}

			songQueueMutex.Lock()
			defer songQueueMutex.Unlock()

			// validate if chatter is the one who requested the song in queue before delete
			// or validate if permissions >= mod
			if idx >= len(songQueue) {
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Usage error, !delsong # invalid number",
					ReplyParentMessageID: event.MessageId,
				})
				return
			}
			song := songQueue[idx]
			if !(song.RequestedByUserID == event.ChatterUserId || isModerator || isBroadcaster) {
				// cannot delete song
				useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Nice try, but you can't delete other chatter's songs!",
					ReplyParentMessageID: event.MessageId,
				})
				return
			}

			// commit to delete song
			songQueue = append(songQueue[:idx], songQueue[idx+1:]...)
			hasCleanupError := false

			// find song to remove from queue if it is the next one
			if idx == 0 {
				// validate if song was really added
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
						log.Println("!delsong cleanup: Failed to delete song from pear-desktop, proceeding anyway...")
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
						log.Println("!delsong cleanup: Failed to add next song in queue to pear-desktop, make sure you add it next in queue yourself. https://youtu.be/" + songQueue[0].Song.VideoID)
						hasCleanupError = true
					}
				}
			}
			msg := "Removed song: " + song.Song.Title + " - " + song.Song.Artist + " " + "https://youtu.be/" + song.Song.VideoID
			if hasCleanupError {
				msg += " but failed to cleanup properly. @" + event.BroadcasterUserLogin + " Please add http://youtu.be/" + songQueue[0].Song.VideoID + " yourself next in queue, Sorry!"
			}
			useProperHelix.SendChatMessage(&helix.SendChatMessageParams{
				BroadcasterID:        event.BroadcasterUserId,
				SenderID:             properUserID,
				Message:              msg,
				ReplyParentMessageID: event.MessageId,
			})

			return
		}
	})
}

var skipMutex = sync.Mutex{}
var lastSkipped = time.Now().Add(time.Second * -10)

var currentSongMutex = sync.Mutex{}
var lastUsedCurrentSong = time.Now().Add(time.Second * -10)

var queueCmdMutex = sync.Mutex{}
var lastUsedQueueCmd = time.Now().Add(time.Second * -10)
