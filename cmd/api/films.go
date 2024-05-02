package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cliftondavies/filmapi/internal/data"
	"github.com/cliftondavies/filmapi/internal/validator"
)

func (app *application) createFilmHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
		Year int32 `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	film := &data.Film{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	}

	v := validator.New()

	if data.ValidateFilm(v, film); !v.Valid() {
		app.failedValidationResponse(w, r , v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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