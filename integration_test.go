package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
    "github.com/gorilla/mux"
)

func TestAddMovieHandler(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expect a single INSERT statement to be executed
	mock.ExpectPrepare("INSERT INTO movies").ExpectExec().WithArgs("Movie Title", 2020, 1).WillReturnResult(sqlmock.NewResult(1, 1))

    // Create a request body
	body := []byte(`{"title":"Movie Title","year":2020,"watched":1}`)

	// Create a request to pass to the handler
	req, err := http.NewRequest("POST", "/movies", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := AddMovieHandler

	// Call the handler with the request, the ResponseRecorder, and the mock database
	defer func() {
		if r := recover(); r != nil {
            debug.PrintStack()
			t.Errorf("handler panicked: %v", r)
		}
	}()
	handler(db, rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Ensure that all expected database queries were executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestListMoviesHandler(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expect a single SELECT statement to be executed
	rows := sqlmock.NewRows([]string{"id", "title", "year", "watched"}).
		AddRow(1, "Movie 1", 2020, 120).
		AddRow(2, "Movie 2", 2021, 90)
	mock.ExpectPrepare("SELECT id, title, year, watched FROM movies").ExpectQuery().WillReturnRows(rows)

	// Create a request to pass to the handler
	req, err := http.NewRequest("GET", "/movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := ListMoviesHandler

	// Call the handler with the request, the ResponseRecorder, and the mock database
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			t.Errorf("handler panicked: %v", r)
		}
	}()
	handler(db, rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := `[{"id":1,"title":"Movie 1","year":2020,"watched":120},{"id":2,"title":"Movie 2","year":2021,"watched":90}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	// Ensure that all expected database queries were executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdateMovieHandler(t *testing.T) {
	// Create a mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Expect a single SELECT statement to be executed
	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)
	mock.ExpectQuery("SELECT id FROM movies WHERE id=?").WithArgs(1).WillReturnRows(rows)

	// Expect a single UPDATE statement to be executed
	mock.ExpectPrepare("UPDATE movies SET").ExpectExec().WithArgs("Updated Movie", 2021, 0, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Create a request body
	body := []byte(`{"title":"Updated Movie","year":2021,"watched":0}`)

	// Create a request to pass to the handler
	req, _ := http.NewRequest("PUT", "/movies/1", bytes.NewReader(body))

    //Hack to try to fake gorilla/mux vars
    vars := map[string]string{
        "id": "1",
    }
    req = mux.SetURLVars(req, vars)

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := UpdateMovieHandler

	// Call the handler with the request, the ResponseRecorder, and the mock database
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			t.Errorf("handler panicked: %v", r)
		}
	}()
	handler(db, rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	if rr.Body.String() != "" {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), "")
	}

	// Ensure that all expected database queries were executed
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
