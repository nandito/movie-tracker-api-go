package main

import (
	"database/sql"
    "fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Movie struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Year   int    `json:"year"`
    Watched int    `json:"watched"`
}

func addMovie(db *sql.DB, movie Movie) (int, error) {
    // Prepare the INSERT statement
    stmt, err := db.Prepare("INSERT INTO movies (title, year, watched) VALUES (?, ?, ?)")
    if err != nil {
        return 0, err
    }
    defer stmt.Close()

    // Execute the statement
    result, err := stmt.Exec(movie.Title, movie.Year, movie.Watched)
    if err != nil {
        return 0, err
    }

    // Get the ID of the newly-inserted row
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}

func getMovies(db *sql.DB) ([]Movie, error) {
    // Prepare the SELECT statement
    stmt, err := db.Prepare("SELECT id, title, year, watched FROM movies")
    if err != nil {
        return nil, err
    }
    defer stmt.Close()

    // Execute the statement
    rows, err := stmt.Query()
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    // Iterate over the rows
    var movies []Movie
    for rows.Next() {
        var movie Movie
        if err := rows.Scan(&movie.ID, &movie.Title, &movie.Year, &movie.Watched); err != nil {
            return nil, err
        }
        movies = append(movies, movie)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return movies, nil
}

func updateMovie(db *sql.DB, id int, movie Movie) error {
    err := db.QueryRow("SELECT id FROM movies WHERE id=?", id).Scan(&id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("Movie not found")
	}
	if err != nil {
		return err
	}
    // Prepare the UPDATE statement
    stmt, err := db.Prepare("UPDATE movies SET title = ?, year = ?, watched = ? WHERE id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the statement
    _, err = stmt.Exec(movie.Title, movie.Year, movie.Watched, id)
    if err != nil {
        return err
    }

    return nil
}

func deleteMovie(db *sql.DB, id int) error {
    err := db.QueryRow("SELECT id FROM movies WHERE id=?", id).Scan(&id)
	if err == sql.ErrNoRows {
		return fmt.Errorf("Movie not found")
	}
	if err != nil {
		return err
	}
    // Prepare the DELETE statement
    stmt, err := db.Prepare("DELETE FROM movies WHERE id = ?")
    if err != nil {
        return err
    }
    defer stmt.Close()

    // Execute the statement
    _, err = stmt.Exec(id)
    if err != nil {
        return err
    }

    return nil
}

