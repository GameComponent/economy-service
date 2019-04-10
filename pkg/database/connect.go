package database

import (
	"database/sql"
	"fmt"
)

// Connect to a database
func Connect(
	host string,
	port string,
	user string,
	password string,
	dbname string,
	sslmode string,
) (*sql.DB, error) {

	connectStringBase := fmt.Sprintf(
		"host=%s port=%s sslmode=%s",
		host,
		port,
		sslmode,
	)

	fmt.Println(connectStringBase)

	connectString := fmt.Sprintf(
		"%s dbname=%s user=%s password=%s",
		connectStringBase,
		dbname,
		user,
		password,
	)

	// Setup the database
	db, err := sql.Open(
		"postgres",
		connectString,
	)

	if err != nil {
		return nil, err
	}

	return db, nil
}
