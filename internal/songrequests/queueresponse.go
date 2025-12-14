package songrequests

type QueueResponse struct {
	Items []struct {
		PlaylistPanelVideoRenderer *struct {
			VideoId         string `json:"videoId"`
			Selected        bool   `json:"selected"`
			ShortBylineText struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"shortBylineText"`
			// LongBylineText struct {
			// 	Runs []struct {
			// 		Text string `json:"text"`
			// 	} `json:"runs"`
			// } `json:"longBylineText"`
			Title struct {
				Runs []struct {
					Text string `json:"text"`
				} `json:"runs"`
			} `json:"title"`
		} `json:"playlistPanelVideoRenderer"`
	} `json:"items"`
}
