//go:build windows

package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/toast.v1"
)

func pushTokenToast(title, message string) {
	go func() {
		notification := toast.Notification{
			AppID:   "Pear Desktop Twitch Song Requests",
			Title:   title,
			Message: message,
		}
		if err := notification.Push(); err != nil {
			log.Printf("Failed to show token notification: %v", err)
		}
	}()
}

// notifyTokenExpiresSoon fires a native Windows toast notification if the
// given Twitch token is authenticated and expires within 7 days.
func notifyTokenExpiresSoon(account string, expiresDate time.Time, isAuthenticated bool) {
	if !isAuthenticated {
		return
	}
	if !time.Now().Add(7 * 24 * time.Hour).After(expiresDate) {
		return
	}
	days := int(time.Until(expiresDate).Hours() / 24)
	pushTokenToast(
		"Twitch Token Expiring Soon",
		fmt.Sprintf("Your %s token expires in %d day(s). Re-authenticate in the app.", account, days),
	)
}

// notifyIfTokenExpired fires a native Windows toast notification if a token
// was stored previously but is no longer valid (expired since last use).
func notifyIfTokenExpired(account string, accessToken string, isAuthenticated bool) {
	if accessToken == "" || isAuthenticated {
		return
	}
	pushTokenToast(
		"Twitch Token Expired",
		fmt.Sprintf("Your %s token has expired. Re-authenticate in the app.", account),
	)
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
	pushTokenToast(
		"Twitch Re-auth Required Soon",
		fmt.Sprintf("Your %s account must be re-authenticated within %d day(s) (30-day inactivity limit).", account, days),
	)
}
