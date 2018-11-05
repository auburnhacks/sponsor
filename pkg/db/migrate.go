package db

import (
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file" // this is a required import by the migration package. oh well!
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const migrationSourceFiles = "file://migrations"

// MigrateDB is a function that always runs the first time when the application
// starts to sync all database fields and run migrations
func MigrateDB(quit chan os.Signal) error {
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
	m, err := migrate.NewWithDatabaseInstance(migrationSourceFiles, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "could not migrate")
	}
	log.Debug("finished performing migrations")
	m.Steps(2)
	go func() {
		log.Debug("watching for quit message to drop db tables")
		signal := <-quit
		log.Debugf("signal %v received, dropping all tables", signal)
		if err := m.Drop(); err != nil {
			log.Fatalf("error dropping migrations: %v", err)
		}
	}()
	return nil
}
