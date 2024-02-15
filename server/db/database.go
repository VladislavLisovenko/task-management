package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/stdlib" // database doesn't work without it
)

var database *sql.DB

func ExecuteQuery(command string) (*sql.Rows, error) {
	return database.Query(command)
}

func init() {
	db, err := sql.Open("pgx", "postgres://postgres:sec@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("Failed to load driver: %v", err)
	}
	database = db
}
