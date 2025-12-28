package helpers

import (
	"bufio"
	"log"
	"os"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
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
	log.Println("Database upgrades successfully applied")

	log.Println("Applying Database Transforms...")
	db, err := databaseconn.NewDBConnection()
	if err != nil {
		log.Println("Failed to start database for database transforms")
		log.Println(err)
		log.Println("Press Enter to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		os.Exit(1)
	}
	defer db.Close()
	for _, v := range data.GetDataTransformTypes() {
		isTransformed, err := data.IsDataTransformed(db, v)
		if err != nil {
			if v.IsNecessary() {
				log.Println("Failed to startup database transforms", v.GetKey())
				log.Println(err)
				log.Println("Press Enter to exit...")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				os.Exit(1)
			}
			log.Println("Failed to startup database transforms", v.GetKey())
			log.Println(err)
			log.Println("It is not yet necessary, will re-attempt on next app launch")
		}
		if !isTransformed {
			err := v.Transform(db)
			if err != nil {
				if v.IsNecessary() {
					log.Println("Failed to apply database transforms", v.GetKey())
					log.Println(err)
					log.Println("Press Enter to exit...")
					bufio.NewReader(os.Stdin).ReadBytes('\n')
					os.Exit(1)
				}
				log.Println("Failed to apply database transforms", v.GetKey())
				log.Println(err)
				log.Println("It is not yet necessary, will re-attempt on next app launch")
			}
			err = data.SetDataTransformed(db, v)
			if err != nil {
				if v.IsNecessary() {
					log.Println("Failed to save database transform results", v.GetKey())
					log.Println(err)
					log.Println("Press Enter to exit...")
					bufio.NewReader(os.Stdin).ReadBytes('\n')
					os.Exit(1)
				}
				log.Println("Failed to save database transform results", v.GetKey())
				log.Println(err)
				log.Println("It is not yet necessary, will re-attempt on next app launch")
			}
		}
	}
	log.Println("Done Database Transforms")
	log.Println("Preflight Test Completed")
}
