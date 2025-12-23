package data

import (
	"database/sql"
	"errors"
	"log"

	. "github.com/azuridayo/pear-desktop-twitch-song-requests/gen/table"
	"github.com/go-jet/jet/v2/sqlite"
	. "github.com/go-jet/jet/v2/sqlite"
	"github.com/nicklaw5/helix/v2"
)

type dataTransformTypeBackportRequestersUserID struct{}

var dataTransformBackportRequestersUserID = dataTransformTypeBackportRequestersUserID{}

func (d dataTransformTypeBackportRequestersUserID) GetKey() string {
	return "BACKPORT_REQUESTERS_USER_ID"
}

func (d dataTransformTypeBackportRequestersUserID) Transform(db *sql.DB) error {
	userAccessToken := struct {
		Value string
	}{}
	stmt := SELECT(Settings.Value).FROM(Settings).WHERE(Settings.Key.EQ(String(DB_KEY_TWITCH_ACCESS_TOKEN))).LIMIT(1)
	err := stmt.Query(db, &userAccessToken)
	if err != nil {
		return err
	}

	helixClient, err := helix.NewClient(&helix.Options{
		ClientID: GetTwitchClientID(),
	})
	if err != nil {
		return err
	}

	valid, _, err := helixClient.ValidateToken(userAccessToken.Value)
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("BACKPORT_REQUESTERS_USER_ID: invalid user access token")
	}
	helixClient.SetUserAccessToken(userAccessToken.Value)

	// fetch all srr with empty user_id
	rowsNeedFix := []struct {
		RowID          int64
		TwitchUsername string
	}{}
	stmt = SELECT(sqlite.RawString("rowid").AS("row_id"), SongRequestRequesters.TwitchUsername.AS("twitch_username")).FROM(SongRequestRequesters).WHERE(SongRequestRequesters.UserID.EQ(String("")))
	err = stmt.Query(db, &rowsNeedFix)
	if err != nil {
		return err
	}

	// paginate get users from chatter login
	logins := []string{}
	mData := map[string]string{}
	for _, v := range rowsNeedFix {
		if _, ok := mData[v.TwitchUsername]; !ok {
			mData[v.TwitchUsername] = v.TwitchUsername
			logins = append(logins, v.TwitchUsername)
		}
	}

	batches := len(logins) / 100
	if len(logins)%100 > 0 {
		batches++
	}

	if batches == 0 {
		return nil
	}
	for i := range batches {
		batchStart := i * 100
		batchEnd := min(batchStart+100, len(logins))

		resp, err := helixClient.GetUsers(&helix.UsersParams{
			Logins: logins[batchStart:batchEnd],
		})
		if err != nil {
			return err
		}
		for _, v := range resp.Data.Users {
			if v2, ok := mData[v.Login]; ok && v2 == "" {
				v2 = v.ID
			}
		}
	}

	caseExpr := CASE(SongRequestRequesters.TwitchUsername)
	for k, v := range mData {
		caseExpr = caseExpr.WHEN(SongRequestRequesters.TwitchUsername.EQ(String(k))).THEN(String(v))
	}
	caseExpr = caseExpr.ELSE(String(""))

	// save to db
	updateStmt := SongRequestRequesters.UPDATE(SongRequestRequesters.UserID).SET(caseExpr).WHERE(SongRequestRequesters.UserID.EQ(String("")))
	log.Println(updateStmt.DebugSql())
	// _, err = updateStmt.Exec(db)
	err = errors.New("Not yet implemented")
	if err != nil {
		return err
	}

	return nil
}

func (d dataTransformTypeBackportRequestersUserID) IsNecessary() bool {
	return false
}
