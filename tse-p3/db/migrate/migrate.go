package migrate

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func Run(connstr string) error {
	var (
		drvr	source.Driver
		mgrt	*migrate.Migrate
		err		error
	)
	drvr, err = iofs.New(migrationFiles, "migrations")
	if err != nil {
		return err
	}

	mgrt, err = migrate.NewWithSourceInstance("iofs", drvr, connstr)
	if err != nil {
		return err
	}
	defer mgrt.Close()

	if err = mgrt.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
			return nil
		}
		if err == migrate.ErrLocked {
			log.Println("Database is locked by another migration process")
			return err
		}
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}