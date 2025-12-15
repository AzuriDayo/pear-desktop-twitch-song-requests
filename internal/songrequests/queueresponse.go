package songrequests

type QueueResponsePlaylistPanelVideoRenderer struct {
	VideoId         string `json:"videoId"`
	Selected        bool   `json:"selected"`
	ShortBylineText struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"shortBylineText"`
	Title struct {
		Runs []struct {
			Text string `json:"text"`
		} `json:"runs"`
	} `json:"title"`
	NavigationEndpoint struct {
		WatchEndpoint struct {
			Index int `json:"index"`
		} `json:"watchEndpoint"`
	} `json:"navigationEndpoint"`
}

type QueueResponse struct {
	Items []struct {
		PlaylistPanelVideoRenderer        *QueueResponsePlaylistPanelVideoRenderer `json:"playlistPanelVideoRenderer"`
		PlaylistPanelVideoWrapperRenderer *struct {
			PrimaryRenderer struct {
				PlaylistPanelVideoRenderer QueueResponsePlaylistPanelVideoRenderer `json:"playlistPanelVideoRenderer"`
			} `json:"primaryRenderer"`
		} `json:"playlistPanelVideoWrapperRenderer"`
	} `json:"items"`
}
