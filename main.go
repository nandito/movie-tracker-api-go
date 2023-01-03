package main

import (
	"database/sql"
    "time"
	"log"
    "net/http"

	_ "github.com/mattn/go-sqlite3"
    "github.com/gorilla/mux"
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


    router := mux.NewRouter()
    
    // Register the HTTP handlers
    router.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPost:
            AddMovieHandler(db, w, r)
        case http.MethodGet:
            ListMoviesHandler(db, w, r)
        default:
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    })
    router.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodPut:
            UpdateMovieHandler(db, w, r)
        case http.MethodDelete:
            DeleteMovieHandler(db, w, r)
        default:
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    })

    srv := &http.Server{
        Handler:      router,
        Addr:         "127.0.0.1:8080",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}

