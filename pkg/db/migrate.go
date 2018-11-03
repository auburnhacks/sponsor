package db

import (
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const migrationSourceFiles = "file://migrations"

func MigrateDB(quit chan os.Signal) error {
	if Conn == nil {
		return errors.New("sql.DB instance not initialized to db.Conn cannot migrate database")
	}
	driver, err := postgres.WithInstance(Conn, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "could not create driver with instance")
	}
	m, err := migrate.NewWithDatabaseInstance(migrationSourceFiles, "postgres", driver)
	if err != nil {
		return errors.Wrap(err, "could not migrate")
	}
	m.Steps(2)
	go func() {
		<-quit
		if err := m.Drop(); err != nil {
			log.Fatalf("error dropping migrations: %v", err)
		}
	}()
	return nil
}
