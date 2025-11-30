package helpers

import (
	"log"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/db"
)

func PreflightTest() {
	log.Println("Starting Main Process")

	// Test Connection Postgres Database
	log.Println("Testing sqlite")
	dB, closeAndUnlock, err := db.NewDBConnection()
	if err != nil {
		log.Println("Failed to connect to sqlite db")
		log.Fatal(err)
	}
	defer closeAndUnlock()
	if err := dB.Ping(); err != nil {
		log.Println("Failed to Ping sqlite")
		log.Fatal(err)
	}
	log.Println("SQLite successfully initialized")

	// Apply sqlite migrations
	log.Println("Applying sqlite migrations")
	err = db.Migrate()
	if err != nil {
		log.Println("Failed to apply sqlite migrations")
		log.Fatal(err)
	}

	// None for now
	log.Println("sqlite migrations successfully applied")
	log.Println("Preflight Test Completed")
}
