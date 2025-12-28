package databaseconn

import (
	"database/sql"
	"log"
	"os"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"
)

const dbName = "pear-desktop-twitch-song-requests.db"

func init() {
	// create db if missing
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// file does not exist
		file, err := os.Create(dbName)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
	}
}

func NewDBConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "file:"+dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func getMigrator() (*migrate.Migrate, error) {
	var m *migrate.Migrate
	var err error
	d, err := iofs.New(data.GetMigrationFS(), "iofs/migrations")
	if err != nil {
		return nil, err
	}
	m, err = migrate.NewWithSourceInstance("iofs", d, "sqlite://"+dbName)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func Migrate() error {
	// create tables if not exist
	var m *migrate.Migrate
	var err error
	m, err = getMigrator()
	if err != nil {
		return err
	}
	defer m.Close()

	log.Println("Please wait while the data structure updates, this can take a few minutes...")
	uperr := m.Up()
	if uperr != nil && uperr != migrate.ErrNoChange {
		log.Println("Failed to upgrade data structure")
		version, dirty, err := m.Version()
		if err != nil {
			log.Println("Failed to get data structure version", err)
			// nothing to do, cannot get version
			return uperr
		}
		if !dirty {
			log.Println("No data corrupted, simply re-launch app")
			// nothing to do, cannot revert non-dirty db
			return uperr
		}
		if version > 0 {
			log.Println("0/2 Attempting recovery...")
			// apply recovery
			err := m.Force(int(version))
			log.Println("1/2 Force data version", version, "successful:", err == nil)
			err = m.Steps(-1)
			log.Println("2/2 Downgrade data version successful:", err == nil)
			// db ready for migration re-attempt
		}
		log.Println("Sorry, something is terribly wrong with the app.")
		log.Println("Take a screenshot of this window and please report this bug to Azuri.")
		return uperr
	}
	log.Println("Data structure updates complete.")
	return nil
}
