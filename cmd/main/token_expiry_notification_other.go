//go:build !windows

package main

import (
	"log"
	"time"
)

func notifyTokenExpiresSoon(account string, expiresDate time.Time, isAuthenticated bool) {
	if !isAuthenticated {
		return
	}
	if !time.Now().Add(7 * 24 * time.Hour).After(expiresDate) {
		return
	}
	days := int(time.Until(expiresDate).Hours() / 24)
	log.Printf("ALERT! %s Twitch token expires in %d day(s). Re-authenticate in the app.", account, days)
}

func notifyIfTokenExpired(account string, accessToken string, isAuthenticated bool) {
	if accessToken == "" || isAuthenticated {
		return
	}
	log.Printf("ALERT! %s Twitch token has expired. Re-authenticate in the app.", account)
}

func notifyRefreshTokenExpiresSoon(account string, lastUsed time.Time, refreshToken string) {
	if refreshToken == "" || lastUsed.IsZero() {
		return
	}
	expiresAt := refreshTokenExpiresAt(lastUsed)
	if !time.Now().Add(twitchRefreshTokenWarnBefore).After(expiresAt) {
		return
	}
	days := int(time.Until(expiresAt).Hours() / 24)
	if days < 0 {
		days = 0
	}
	log.Printf("ALERT! %s Twitch refresh token inactive limit reached in %d day(s). Re-authenticate in the app.", account, days)
}
