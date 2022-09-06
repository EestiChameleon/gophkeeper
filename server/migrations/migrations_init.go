package migration

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	m *migrate.Migrate
)

// InitMigration creates needed tables and functions in the database.
func InitMigration() error {
	// connect
	if err := migrateInitConnect(); err != nil {
		return err
	}
	// create
	if err := upGophKeeperStorage(); err != nil {
		return err
	}
	// close
	srcErr, dbErr := m.Close()
	if srcErr != nil {
		return srcErr
	}
	if dbErr != nil {
		return dbErr
	}
	return nil
}

// upGophKeeperStorage migrates all the way up to the final DB level, found in .sql file description.
func upGophKeeperStorage() error {
	err := m.Up()
	if !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// downGophKeeperStorage migrates all the way up to the final DB level, found in .sql file description.
func downGophKeeperStorage() error {
	return m.Down()
}

// migrateInitConnect establishes connection to the DB for migrate operations.
func migrateInitConnect() error {
	conn, err := sql.Open("postgres", "postgres://localhost:5432/yandex_practicum_db?sslmode=disable")
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		return err
	}

	db, err := migrate.NewWithDatabaseInstance(
		"file://server/migrations/sqlscripts",
		"postgres", driver)
	if err != nil {
		return err
	}

	m = db
	return nil
}
