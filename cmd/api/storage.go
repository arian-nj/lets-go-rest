package main

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/arian-nj/site/back/internal/data"
	"github.com/lib/pq"
)

type storage interface {
	Init() error
	InsertMovie(*data.Movie) error
	GetMovieById(int64) (*data.Movie, error)
	Update(movie *data.Movie) error
	DeleteMovie(int64) error
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return &PostgresStorage{}, fmt.Errorf("unable to connect to database: %v", err)
	}
	// defer conn.Close(context.Background())

	err = conn.Ping()
	if err != nil {
		return &PostgresStorage{}, err
	}
	return &PostgresStorage{
		db: conn,
	}, err
}

func (s *PostgresStorage) Init() error {
	return s.CreateMovieTable()
}

func (s *PostgresStorage) CreateMovieTable() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	query := `CREATE TABLE IF NOT EXISTS movies (
		id bigserial PRIMARY KEY,
		created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
		title text NOT NULL,
		year integer NOT NULL,
		runtime integer NOT NULL,
		genres text[] NOT NULL,
		version integer NOT NULL DEFAULT 1
		);`

	_, err = s.db.Exec(query)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

func (s *PostgresStorage) InsertMovie(movie *data.Movie) error {
	fmt.Println("making ", movie.Title, " in db")
	statment := `INSERT INTO movies 
	(title,year,runtime,genres)
	VALUES 
	($1,$2,$3,$4) 
	RETURNING id, created_at, version
	`
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
	}

	return s.db.QueryRow(statment, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (s *PostgresStorage) GetMovieById(id int64) (*data.Movie, error) {
	if id < 1 {
		return nil, data.ErrRecordNotFound
	}
	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`
	var movie data.Movie

	err := s.db.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, data.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Otherwise, return a pointer to the Movie struct.
	return &movie, nil
}

func (s *PostgresStorage) Update(movie *data.Movie) error {
	query := `
	UPDATE movies
	SET title = $1, year = $2, runtime = $3, genres = $4, version =
	version + 1
	WHERE id = $5
	RETURNING version`
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return s.db.QueryRow(query, args...).Scan(&movie.Version)
}

func (s *PostgresStorage) DeleteMovie(id int64) error {
	if id < 1 {
		return data.ErrRecordNotFound
	}
	query := `
	DELETE FROM movies
	WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsEffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsEffected == 0 {
		return data.ErrRecordNotFound
	}

	return nil
}
