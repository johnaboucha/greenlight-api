package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.johnboucha.com/internal/data"
)

// handler for "POST /v1/movies"
func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// anonymous struct to hold HTTP request body
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	// decode JSON from body with cmd/api/helpers.go -> readJSON method
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

// handler for "GET /v1/movies/:id"
func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r) // helper function that gets ID
	if err != nil {
		app.notFoundResponse(w, r) // error response from cmd/api/errors.go
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		// something went wrong, throw error
		app.serverErrorResponse(w, r, err)
	}
}
