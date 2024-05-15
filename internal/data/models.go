package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Films FilmModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Films: FilmModel{DB: db},
	}
}