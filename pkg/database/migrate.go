package database

import (
	"fmt"

	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/cockroachdb"

	// Needed for the source file driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
	pq "github.com/lib/pq"
)

const (
	dbErrorDuplicateDatabase = "42P04"
	migrationTable           = "migration_info"
	dialect                  = "postgres"
	defaultLimit             = -1
)

// Init the migration process
func Init(host string, port string, user string, password string, dbname string, sslmode string) (bool, error) {
	connectStringBase := fmt.Sprintf(
		"host=%s port=%s sslmode=%s",
		host,
		port,
		sslmode,
	)

	connectString := fmt.Sprintf(
		"%s user=%s password=%s",
		connectStringBase,
		user,
		password,
	)

	// Setup the database
	_, err := setupDatabase(connectString, dbname)
	if err != nil {
		return false, err
	}

	databaseConnectString := fmt.Sprintf(
		"%s dbname=%s user=%s password=%s",
		connectStringBase,
		dbname,
		user,
		password,
	)

	// Migrate the database
	_, err = migrateDatabase(databaseConnectString)
	if err != nil {
		return false, err
	}

	return true, nil
}

func setupDatabase(connectString string, databaseName string) (bool, error) {
	// Setup the database
	db, err := sql.Open(
		"postgres",
		connectString,
	)

	if err != nil {
		return false, err
	}
	defer db.Close()

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", databaseName)); err != nil {
		if e, ok := err.(*pq.Error); ok && e.Code == dbErrorDuplicateDatabase {
			fmt.Println("Using existing database")
		} else {
			return false, err
		}
	} else {
		fmt.Println("Creating new database")
	}

	return true, nil
}

func migrateDatabase(connectString string) (bool, error) {
	db, err := sql.Open("postgres", connectString)
	if err != nil {
		return false, err
	}
	defer db.Close()

	driver, err := cockroachdb.WithInstance(db, &cockroachdb.Config{})
	if err != nil {
		return false, err
	}

	instance, err := migrate.NewWithDatabaseInstance(
		"file://../../migrations/",
		"cockroachdb",
		driver,
	)
	if err != nil {
		return false, err
	}

	instance.Up()

	return true, nil
}
