package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"log"
	"net/http"
	"strconv"

	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
	"github.com/go-jet/jet/v2/sqlite"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/labstack/echo/v4"
)

type PaginationCommon struct {
	PerPage string `query:"perPage"`
	Page    string `query:"page"`
}

func (a *App) handleRequestersHistory(c echo.Context) error {
	var err error
	p := PaginationCommon{}

	if err = c.Bind(&p); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	page := 0
	if page, err = strconv.Atoi(p.Page); err != nil || page < 1 {
		return c.NoContent(http.StatusBadRequest)
	}
	perPage := 0
	if perPage, err = strconv.Atoi(p.PerPage); err != nil || perPage < 10 || perPage > 100 {
		return c.NoContent(http.StatusBadRequest)
	}

	db, err := databaseconn.NewDBConnection()
	if err != nil {
		log.Println("handleRequestersHistory: failed to open database connection")
		return c.NoContent(http.StatusInternalServerError)
	}
	defer db.Close()

	maxResults := struct {
		MaxResults int64
	}{}
	queryStmt := SELECT(MAX(sqlite.IntegerColumn("rowid")).AS("max_results")).FROM(SongRequestRequesters)
	err = queryStmt.QueryContext(c.Request().Context(), db, &maxResults)
	if err != nil {
		log.Println("handleRequestersHistory: failed to query max page")
		return c.NoContent(http.StatusInternalServerError)
	}
	maxPage := maxResults.MaxResults / int64(perPage)
	if maxResults.MaxResults%int64(perPage) > 0 {
		maxPage++
	}

	if int64(page) > maxPage {
		log.Println("handleRequestersHistory: page > max_page")
		return c.NoContent(http.StatusBadRequest)
	}

	results := []struct {
		VideoID        string `json:"video_id"`
		TwitchUsername string `json:"twitch_username"`
		RequestedAt    string `json:"requested_at"`
		IsNinja        bool   `json:"is_ninja"`
		SongTitle      string `json:"song_title"`
		ArtistName     string `json:"artist_name"`
		ImageURL       string `json:"image_url"`
	}{}
	queryStmt = SELECT(SongRequestRequesters.VideoID.AS("video_id"), SongRequestRequesters.TwitchUsername.AS("twitch_username"), SongRequestRequesters.RequestedAt.AS("requested_at"), SongRequestRequesters.IsNinja.AS("is_ninja"), SongRequests.SongTitle.AS("song_title"), SongRequests.ArtistName.AS("artist_name"), SongRequests.ImageURL.AS("image_url")).FROM(SongRequestRequesters.LEFT_JOIN(SongRequests, SongRequests.VideoID.EQ(SongRequestRequesters.VideoID))).ORDER_BY(sqlite.IntegerColumn("rowid").DESC()).LIMIT(int64(perPage)).OFFSET(int64((page - 1) * perPage))
	err = queryStmt.QueryContext(c.Request().Context(), db, &results)
	if err != nil {
		log.Println("handleRequestersHistory: failed to query data", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"max_pages": maxPage,
		"items":     results,
	})

}
