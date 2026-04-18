package helpers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
)

type FoundVideoData struct {
	Index int
}

func FindAllVideoIDCounterparts(videoID string) (bool, map[string]FoundVideoData) {
	// get queue
	queue := songrequests.QueueResponse{}
	resp, err := http.Get("http://" + songrequests.GetPearDesktopHost() + "/api/v1/queue")
	if err == nil && resp.StatusCode == http.StatusOK {
		var qb []byte
		qb, err = io.ReadAll(resp.Body)
		if err == nil {
			err = json.Unmarshal(qb, &queue)
			if err == nil {
				resp.Body.Close()
				// compare all videoid
				// get all video ids and counterparts
				for i := len(queue.Items) - 1; i >= 0; i-- {
					v := queue.Items[i]
					compareVideoIDs := map[string]FoundVideoData{}
					// all counterparts
					if v.PlaylistPanelVideoWrapperRenderer != nil {
						compareVideoIDs[v.PlaylistPanelVideoWrapperRenderer.PrimaryRenderer.PlaylistPanelVideoRenderer.VideoId] = FoundVideoData{}
						for _, v2 := range v.PlaylistPanelVideoWrapperRenderer.Counterpart {
							compareVideoIDs[v2.CounterpartRenderer.PlaylistPanelVideoRenderer.VideoId] = FoundVideoData{}
						}
					}
					// native videoid
					if v.PlaylistPanelVideoRenderer != nil {
						compareVideoIDs[v.PlaylistPanelVideoRenderer.VideoId] = FoundVideoData{ Index: i }
					}

					// compare newVideoId and compareVideoIds and queueHead set in isQueueHeadPlayingNow
					okVideoId := false
					_, okVideoId = compareVideoIDs[videoID]
					if okVideoId {
						return true, compareVideoIDs
					}
				}
			}
		}
	}
	return false, nil
}
