package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
	"github.com/joeyak/go-twitch-eventsub/v3"
	"github.com/nicklaw5/helix/v2"
)

var checkMainChannelUserStatusMutex = sync.RWMutex{}
var checkMainChannelUserStatus = map[string]struct {
	isSub       bool
	isModerator bool
	timeExpiry  time.Time
}{}

func (a *App) SetSubscriptionHandlersBot() {
	a.twitchWSBotService.Client().OnEventChannelChatMessage(func(event twitch.EventChannelChatMessage) {
		isSub := false
		isBroadcaster := false
		isModerator := false
		isVip := false
		isSelf := false
		useProperHelix := a.helixBot
		properUserID := a.twitchDataStructBot.userID
		realBroadcasterID := a.twitchDataStruct.userID
		trimmedText := strings.TrimSpace(event.Message.Text)
		trimmedText = strings.Trim(trimmedText, " ͏") // idk why twitch adds this character

		if strings.EqualFold(event.ChatterUserLogin, a.twitchDataStructBot.login) {
			isSelf = true
		} else if strings.EqualFold(event.ChatterUserLogin, a.twitchDataStruct.login) {
			isSub = true
			isBroadcaster = true
			isModerator = true
		} else {
			checkMainChannelUserStatusMutex.RLock()
			if v, ok := checkMainChannelUserStatus[event.ChatterUserLogin]; ok && !time.Now().After(v.timeExpiry) {
				isSub = v.isSub
				isModerator = v.isModerator
				checkMainChannelUserStatusMutex.RUnlock()
			} else {
				checkMainChannelUserStatusMutex.RUnlock()

				subsResponse, err := a.helix.GetSubscriptions(&helix.SubscriptionsParams{
					UserID:        []string{event.ChatterUserId},
					BroadcasterID: realBroadcasterID,
				})
				if err != nil {
					emsg := "Internal error when checking if " + event.ChatterUserLogin + " is a sub"
					log.Println(emsg, err)
					a.helixBot.SendChatMessage(&helix.SendChatMessageParams{
						BroadcasterID: event.BroadcasterUserId,
						SenderID:      properUserID,
						Message:       emsg,
					})
					return
				}
				if len(subsResponse.Data.Subscriptions) > 0 {
					isSub = true
				}

				vipsResponse, err := a.helix.GetChannelVips(&helix.GetChannelVipsParams{
					UserID:        event.ChatterUserId,
					BroadcasterID: realBroadcasterID,
				})
				if err != nil {
					emsg := "Internal error when checking if " + event.ChatterUserLogin + " is a vip"
					log.Println(emsg, err)
					a.helixBot.SendChatMessage(&helix.SendChatMessageParams{
						BroadcasterID: event.BroadcasterUserId,
						SenderID:      properUserID,
						Message:       emsg,
					})
					return
				}
				if len(vipsResponse.Data.ChannelsVips) > 0 {
					isVip = true
				}

				modsResponse, err := a.helix.GetModerators(&helix.GetModeratorsParams{
					UserIDs:       []string{event.ChatterUserId},
					BroadcasterID: realBroadcasterID,
				})
				if err != nil {
					emsg := "Internal error when checking if " + event.ChatterUserLogin + " is a moderator"
					log.Println(emsg, err)
					a.helixBot.SendChatMessage(&helix.SendChatMessageParams{
						BroadcasterID: event.BroadcasterUserId,
						SenderID:      properUserID,
						Message:       emsg,
					})
					return
				}
				if len(modsResponse.Data.Moderators) > 0 {
					isSub = true
					isModerator = true
				}

				checkMainChannelUserStatusMutex.Lock()
				checkMainChannelUserStatus[event.ChatterUserLogin] = struct {
					isSub       bool
					isModerator bool
					timeExpiry  time.Time
				}{
					isSub:       isSub,
					isModerator: isModerator,
					timeExpiry:  time.Now().Add(time.Hour * 2),
				}
				checkMainChannelUserStatusMutex.Unlock()
			}
		}
		if isSelf {
			// do not process bot self msgs.
			return
		}

		if (isSub || isVip) && len(trimmedText) > 4 && strings.EqualFold(trimmedText[:4], "!sr ") {
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
			if time.Now().After(lastSkipped.Add(time.Second * -10)) {
				hasSkipped = true
				songQueueMutex.Lock()
				http.Post("http://"+songrequests.GetPearDesktopHost()+"/api/v1/next", "application/json", nil)
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
				a.helixBot.SendChatMessage(&helix.SendChatMessageParams{
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
			currentSongMutexBot.Lock()
			if !time.Now().After(lastUsedCurrentSongBot.Add(time.Second * 10)) {
				currentSongMutexBot.Unlock()
				return
			}
			lastUsedCurrentSongBot = time.Now()
			currentSongMutexBot.Unlock()

			resp, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/song")
			if err == nil {
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
				a.helixBot.SendChatMessage(&helix.SendChatMessageParams{
					BroadcasterID:        event.BroadcasterUserId,
					SenderID:             properUserID,
					Message:              "Internal failure to get song details!",
					ReplyParentMessageID: event.MessageId,
				})
			} else {
				a.helixBot.SendChatMessage(&helix.SendChatMessageParams{
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
			queueCmdMutexBot.Lock()
			if !time.Now().After(lastUsedQueueCmdBot.Add(time.Second * 10)) {
				queueCmdMutexBot.Unlock()
				return
			}
			lastUsedQueueCmdBot = time.Now()
			queueCmdMutexBot.Unlock()

			songQueueMutex.RLock()
			defer songQueueMutex.RUnlock()

			song := songrequests.SongResult{}

			resp, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/song")
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

		if (isModerator || isBroadcaster) && strings.EqualFold(trimmedText, "!version") {
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
	})
}

var currentSongMutexBot = sync.Mutex{}
var lastUsedCurrentSongBot = time.Now().Add(time.Second * -10)

var queueCmdMutexBot = sync.Mutex{}
var lastUsedQueueCmdBot = time.Now().Add(time.Second * -10)
