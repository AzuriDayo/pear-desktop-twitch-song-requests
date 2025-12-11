package helpers

import (
	"bufio"
	"log"
	"os"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/databaseconn"
)

func PreflightTest() {
	log.Println("Starting preflight tests")

	// Test Connection Postgres Database
	log.Println("Testing database")
	// Apply database migrations
	log.Println("Applying database upgrades")
	err := databaseconn.Migrate()
	if err != nil {
		log.Println("Failed to apply database upgrades")
		log.Println(err)
		log.Println("Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}

	// None for now
	log.Println("Database upgrades successfully applied")
	log.Println("Preflight Test Completed")
}
