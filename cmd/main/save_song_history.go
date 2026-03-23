package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"log"
	"time"

	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/gen/model"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/songrequests"
)

func saveSongHistory(song *songrequests.SongResult, chatterUserLogin, chatterUserId string, isNinja bool) {
	db, err := databaseconn.NewDBConnection()
	if err != nil {
		log.Println("Somehow failed to add !sr history to database")
		return
	}
	defer db.Close()
	srData := model.SongRequests{
		VideoID:    song.VideoID,
		SongTitle:  song.Title,
		ArtistName: song.Artist,
		ImageURL:   song.ImageUrl,
	}
	stmt := SongRequests.INSERT(SongRequests.AllColumns).MODEL(srData).ON_CONFLICT(SongRequests.VideoID).DO_NOTHING()
	_, err = stmt.Exec(db)
	if err != nil {
		log.Println("Somehow failed to save !sr song to database")
		return
	}
	srrData := model.SongRequestRequesters{
		VideoID:        song.VideoID,
		TwitchUsername: chatterUserLogin,
		RequestedAt:    time.Now().Local().Format(data.TWITCH_SERVER_DATE_LAYOUT),
		UserID:         chatterUserId,
		IsNinja:        isNinja,
	}
	stmt = SongRequestRequesters.INSERT(SongRequestRequesters.AllColumns).MODEL(srrData)
	_, err = stmt.Exec(db)
	if err != nil {
		log.Println("Somehow failed to save !sr requester name to database")
		return
	}
}
