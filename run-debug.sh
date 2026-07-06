#!/bin/bash
# Set TWITCH_CLIENT_ID from the Twitch Developer Console (no client secret required).
export TWITCH_CLIENT_ID="${TWITCH_CLIENT_ID:?set TWITCH_CLIENT_ID}"
go run -ldflags "-X github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data.twitchClientID=${TWITCH_CLIENT_ID}" \
    ./cmd/main/
