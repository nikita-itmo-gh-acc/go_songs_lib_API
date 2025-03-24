package storage

import (
	"database/sql"
	"os"
	"songsapi/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Migrator struct {
	MigrationTool *migrate.Migrate
}

func CreateMigrator(db *sql.DB) (*Migrator, error) {
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Err.Println("can't get postres driver")
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", os.Getenv("DB_NAME"), dbDriver)
	if err != nil {
		logger.Err.Println("can't create migrator instance")
		return nil, err
	}
	
	return &Migrator{MigrationTool: m}, nil
}
	
func (migrator *Migrator) MakeMigrations() error {
	if err := migrator.MigrationTool.Up(); err != nil {
		if err == migrate.ErrNoChange {
			logger.Debug.Println("no migrations to apply")
			return nil
		}

		logger.Err.Println("migration failed - ", err)
		return err
	}

	logger.Debug.Println("migrations applied!")
	return nil
}

func (migrator *Migrator) RollBack(steps int) error {
	if err := migrator.MigrationTool.Steps(-steps); err != nil {
		if err == migrate.ErrNoChange {
			logger.Debug.Println("no migrations to rollback")
			return nil
		}

		logger.Err.Println("rollback failed - ", err)
		return err
	}

	logger.Debug.Println("migrations rolled back!")
	return nil
}
