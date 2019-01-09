package db

import (
	"os"

	//"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"

	// this is a required import by the migration package. oh well!
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var (
	migrator *migrate.Migrate
)

const migrationSourceFiles = "file://migrations"

// Migrate is a function that always runs the first time when the application
// starts to sync all database fields and run migrations
func Migrate(quit chan os.Signal) error {
	var err error
	log.Debug("starting to perform migrations")
	if Conn == nil {
		return errors.New("sql.DB instance not initialized to db.Conn cannot migrate database")
	}
	log.Debug("creating postgres driver for migrations")
	driver, err := postgres.WithInstance(Conn.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "could not create driver with instance")
	}
	log.Debugf("migrating database to with driver: %v", driver)
	migrator, err = migrate.NewWithDatabaseInstance(migrationSourceFiles, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "could not migrate")
	}
	log.Debug("finished performing migrations")
	migrator.Steps(2)
	return nil
}

// Drop is a function that drop all the tables and their data
func Drop() error {
	log.Debug("dropping tables...")
	if err := migrator.Drop(); err != nil {
		return err
	}
	return nil
}
