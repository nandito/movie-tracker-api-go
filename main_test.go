package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

