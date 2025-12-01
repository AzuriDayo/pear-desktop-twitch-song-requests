package helpers

import (
	"bufio"
	"log"
	"os"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
)

func PreflightTest() {
	log.Println("Starting Main Process")

	// Test Connection Postgres Database
	log.Println("Testing sqlite")
	// Apply sqlite migrations
	log.Println("Applying sqlite migrations")
	err := databaseconn.Migrate()
	if err != nil {
		log.Println("Failed to apply sqlite migrations")
		log.Println(err)
		log.Println("Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}

	// None for now
	log.Println("sqlite migrations successfully applied")
	log.Println("Preflight Test Completed")
}
