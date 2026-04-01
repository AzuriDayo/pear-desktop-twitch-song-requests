package main

import (
	"log"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
)

func main() {
	songs := []string{
		"!sr yena nemonemo",
		"!sr yena good morning",
		"!sr yena being a good girl hurts",
		"!sr yena smartphone",
		"!sr yena good night",              // MusicShelfRenderer instead of MusicCardShelfRenderer, search results as yena - good morning
		"!sr l-txZMzBmbA",                  // This is a podcast episode
		"!sr https://youtu.be/38xYeot-ciM", // This will never work because it is not searchable. Always results in wrong video, which is expected. The actual implementation is different than this blast searches script
		"!sr https://music.youtube.com/watch?v=b4iVv91Z6lY&si=",
		"!sr https://www.youtube.com/watch?v=9dxiTuPy9dQ",
	}
	for _, v := range songs {
		v, _ := songrequests.ParseSearchQuery(v)
		song, err := songrequests.SearchSong(v, 30, 600)
		if err != nil {
			switch err {
			case songrequests.ErrNoResults:
				fallthrough
			case songrequests.ErrSongLength:
				log.Println(v, err)
			default:
				panic(err)
			}
			continue
		}
		log.Println(song.Title, song.VideoID, song.SearchOrigin, song.ImageUrl, song.RawTimeData)
	}
}
