package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func AddMovieHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Parse the request body and bind it to a Movie struct
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Add the movie to the database
	id, err := addMovie(db, movie)
	if err != nil {
		http.Error(w, "Failed to add movie", http.StatusInternalServerError)
		return
	}

	// Set the ID of the movie and write the response
	movie.ID = id
	response, err := json.Marshal(movie)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func ListMoviesHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Get a list of movies from the database
	movies, err := getMovies(db)
	if err != nil {
		http.Error(w, "Error getting movies", http.StatusInternalServerError)
		return
	}

	// Write the response
	response, err := json.Marshal(movies)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func UpdateMovieHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Get the movie ID from the URL parameters
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Parse the request body and bind it to a Movie struct
	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the movie in the database
	if err := updateMovie(db, id, movie); err != nil {
        if err.Error() == "Movie not found" {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error updating movie", http.StatusInternalServerError)
		return
	}
}

func DeleteMovieHandler(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// Get the movie ID from the URL parameters
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Delete the movie from the database
	if err := deleteMovie(db, id); err != nil {
        if err.Error() == "Movie not found" {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error deleting movie", http.StatusInternalServerError)
		return
	}
}

