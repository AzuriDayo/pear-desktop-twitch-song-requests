//go:build !windows

package main

import (
	"log"
	"time"
)

// notifyTokenExpiresSoon logs a warning if the given Twitch token is
// authenticated and expires within 7 days.
func notifyTokenExpiresSoon(account string, expiresDate time.Time, isAuthenticated bool) {
	if !isAuthenticated {
		return
	}
	if !time.Now().Add(7 * 24 * time.Hour).After(expiresDate) {
		return
	}
	days := int(time.Until(expiresDate).Hours() / 24)
	log.Printf("ALERT! %s Twitch token expires in %d day(s). Consider refreshing the token at http://localhost:3999/", account, days)
}

// notifyIfTokenExpired logs a warning if a token was stored previously but is
// no longer valid (expired since last use).
func notifyIfTokenExpired(account string, accessToken string, isAuthenticated bool) {
	if accessToken == "" || isAuthenticated {
		return
	}
	log.Printf("ALERT! %s Twitch token has expired. Re-authenticate at http://localhost:3999/", account)
}
