package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cliftondavies/filmapi/internal/data"
)

func (app *application) createFilmHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

func (app *application) showFilmHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	
	film := envelope{
		"film": data.Film{
			ID: id,
			CreatedAt: time.Now(),
			Title: "Casablanca",
			Runtime: 102,
			Genres: []string{"drama", "romance", "war"},
			Version: 1,
		},
	}

	err = app.writeJson(w, http.StatusOK, film, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}