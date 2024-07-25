package main

import (
	"context"
	"fmt"
	"time"

	"github.com/arian-nj/site/back/internal/data"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type storage interface {
	Init() error
	InsertMovie(*data.Movie) error
	GetMovieById(int64) (*data.Movie, error)
	DeleteMovie(int) error
}

type PostgresStorage struct {
	db *pgx.Conn
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connStr := "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return &PostgresStorage{}, fmt.Errorf("unable to connect to database: %v", err)
	}
	// defer conn.Close(context.Background())

	err = conn.Ping(context.Background())
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
	tx, err := s.db.Begin(context.Background())
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

	_, err = s.db.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	return err
}

func (s *PostgresStorage) InsertMovie(movie *data.Movie) error {
	fmt.Println("making ", movie.Title, " in db")
	statment := `INSERT INTO movies 
	(created_at,title,year,runtime,genres,version)
	VALUES 
	($1,$2,$3,$4,$5,$6)
	`

	rows, err := s.db.Exec(context.Background(), statment,
		time.Now().UTC(),
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.Version)
	fmt.Printf("%+v", rows)
	return err
}

func (s *PostgresStorage) GetMovieById(id int64) (*data.Movie, error) {
	var movie data.Movie
	err := s.db.QueryRow(context.Background(), "SELECT * FROM movies WHERE id = $1", id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		&movie.Genres,
		&movie.Version,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &movie, nil
}

func (s *PostgresStorage) Update(movie *data.Movie) error {
	// query := `UPDATE movies
	// SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
	// WHERE id = $5
	// RETURNING version`
	// err := s.db.QueryRow(context.Background(),
	// query, movie.Title, movie.Year, movie.Runtime, movie.Genres, movie.Version,)
	return nil
}

func (s *PostgresStorage) DeleteMovie(id int) error {
	// _, err := s.db.Query("delete FROM movie where id = $1", id)
	err := fmt.Errorf("s")
	return err
}
