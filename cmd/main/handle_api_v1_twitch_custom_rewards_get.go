package main

//lint:file-ignore ST1001 Dot imports by jet
import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nicklaw5/helix/v2"
)

type selectableReward struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Cost int    `json:"cost"`
}

func (a *App) twitchCustomRewardsGET(c echo.Context) error {
	if a.helix == nil {
		return c.NoContent(http.StatusUnauthorized)
	}
	// fetch all custom rewards
	rewards, err := a.helix.GetCustomRewards(&helix.GetCustomRewardsParams{
		BroadcasterID: a.twitchDataStruct.userID,
	})
	if err != nil {
		e := "Failed to list custom rewards from channel " + a.twitchDataStruct.login
		log.Println(e)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": e,
		})
	}
	selectableRewards := []selectableReward{}
	for _, v := range rewards.Data.ChannelCustomRewards {
		if v.IsUserInputRequired {
			selectableRewards = append(selectableRewards, selectableReward{
				ID:   v.ID,
				Name: v.Title,
				Cost: v.Cost,
			})
		}
	}

	return c.JSON(http.StatusOK, selectableRewards)
}
