package data

import (
	"time"

	"github.com/cliftondavies/filmapi/internal/validator"
)

type Film struct {
	ID int64 `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title string `json:"title"`
	Year int32 `json:"year,omitempty"`
	Runtime Runtime `json:"runtime,omitempty"`
	Genres []string `json:"genres,omitempty"`
	Version int32 `json:"version"`
}

func ValidateFilm(v *validator.Validator, film *Film) {
	v.Check(film.Title != "", "title", "must be provided")
	v.Check(len(film.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(film.Year != 0, "year", "must be provided")
	v.Check(film.Year >= 1888, "year", "must be greater than 1888")
	v.Check(film.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(film.Runtime != 0, "runtime", "must be provided")
	v.Check(film.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(film.Genres != nil, "genres", "must be provided")
	v.Check(len(film.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(film.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(film.Genres), "genres", "must not contain duplicate values")
}
	