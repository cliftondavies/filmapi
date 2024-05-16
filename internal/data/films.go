package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cliftondavies/filmapi/internal/validator"

	"github.com/lib/pq"
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

type FilmModel struct {
	DB *sql.DB
}

func (f FilmModel) Insert(film *Film) error {
	query := `
		INSERT INTO films (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	args := []any{film.Title, film.Year, film.Runtime, pq.Array(film.Genres)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return f.DB.QueryRowContext(ctx, query, args...).Scan(&film.ID, &film.CreatedAt, &film.Version)
}

func (f FilmModel) Get(id int64) (*Film, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, title, year, runtime, genres, version
		FROM films
		WHERE id = $1`

	var film Film

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := f.DB.QueryRowContext(ctx, query, id).Scan(
		&film.ID,
		&film.CreatedAt,
		&film.Title,
		&film.Year,
		&film.Runtime,
		pq.Array(&film.Genres),
		&film.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &film, nil
}

func (f FilmModel) Update(film *Film) error {
	query := `
		UPDATE films
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []any{
		film.Title,
		film.Year,
		film.Runtime,
		pq.Array(film.Genres),
		film.ID,
		film.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := f.DB.QueryRowContext(ctx, query, args...).Scan(&film.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (f FilmModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM films
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result, err := f.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	
	return nil
}

func (f FilmModel) GetAll(title string, genres []string, filters Filters) ([]*Film, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM films
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{title, pq.Array(genres), filters.limit(), filters.offset()}

	rows, err := f.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	films := []*Film{}

	for rows.Next() {
		var film Film

		err := rows.Scan(
			&totalRecords,
			&film.ID,
			&film.CreatedAt,
			&film.Title,
			&film.Year,
			&film.Runtime,
			pq.Array(&film.Genres),
			&film.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		films = append(films, &film)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return films, metadata, nil
}