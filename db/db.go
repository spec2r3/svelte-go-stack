package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error // Use a new local variable for error

	// Assign the database connection to the package-level DB variable
	DB, err = sql.Open("sqlite3", "api.db")

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	DB.SetMaxIdleConns(5)
	DB.SetMaxOpenConns(10)

	createTables()
}

func createTables() {

	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		api_key TEXT NOT NULL UNIQUE,
		alias TEXT NOT NULL,
		admin INTEGER NOT NULL DEFAULT 0 CHECK (admin IN (0, 1))
	)
	`
	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Could not create table: %v", err)
	}

	createCarTables := `
	CREATE TABLE IF NOT EXISTS cars(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		brand TEXT NOT NULL,
		model TEXT NOT NULL,
		engine TEXT NOT NULL,
		gearbox TEXT NOT NULL,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)
	`

	_, err = DB.Exec(createCarTables)
	if err != nil {
		log.Fatalf("Could not create table: %v", err)
	}
}
