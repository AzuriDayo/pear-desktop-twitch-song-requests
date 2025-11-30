package db

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"github.com/azuridayo/pear-desktop-twitch-song-requests/internal/data"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var mu = &sync.Mutex{}

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

func NewDBConnection() (*sql.DB, func(), error) {
	mu.Lock()
	db, err := sql.Open("sqlite3", "file:"+dbName)
	if err != nil {
		return nil, nil, err
	}

	return db, func() { db.Close(); mu.Unlock() }, nil
}

func getMigrator() (*migrate.Migrate, error) {
	var m *migrate.Migrate
	var err error
	// use migrate with sqlite3
	d, err := iofs.New(data.FS, "iofs/migrations")
	if err != nil {
		return nil, err
	}
	m, err = migrate.NewWithSourceInstance(
		"iofs",
		d,
		"sqlite3://"+dbName)
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
	uperr := m.Up()
	if uperr != nil && uperr != migrate.ErrNoChange {
		// attempt recovery force down
		// log.Println("Migration failed, recovering...")
		// stmt := SELECT(SchemaMigrations.Version, SchemaMigrations.Dirty).FROM(SchemaMigrations).LIMIT(1)
		// schema := model.SchemaMigrations{}
		// err := stmt.Query(db, &schema)
		// if err != nil {
		// 	log.Println("Migration recovery failed to get database migration version")
		// 	return err
		// }
		// if !schema.Dirty {
		// 	log.Println("Migration recovery not necessary, schema is not dirty")
		// 	return uperr
		// }
		// forceVersion := schema.Version - 1
		// log.Println("Migration recovery will force schema back 1 version, to version", forceVersion)
		// err = m.Force(int(forceVersion))
		// if err != nil {
		// 	log.Println("Migration recovery failed to force version to", forceVersion)
		// 	return err
		// }
		// log.Println("Migration recovery forced version to", forceVersion)
		return uperr
	}
	return nil
}
