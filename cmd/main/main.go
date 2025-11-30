package main

import (
	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/helpers"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	helpers.PreflightTest()
}
