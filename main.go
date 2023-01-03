package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Open a connection to the database
	db, err := sql.Open("sqlite3", "./movies.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the "movies" table
	query := `
	CREATE TABLE IF NOT EXISTS movies (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		year INTEGER NOT NULL,
		watched INTEGER NOT NULL
	);`
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
}

