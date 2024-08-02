package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movie MovieModel
	Users UserModel
	Token TokenModel
}

func NewModels() (*Models, error) {
	connStr := "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)
	conn.SetConnMaxIdleTime(time.Minute * 15)

	if err != nil {
		return &Models{}, fmt.Errorf("unable to connect to database: %v", err)
	}
	// defer conn.Close(context.Background())

	err = conn.Ping()
	if err != nil {
		return &Models{}, err
	}
	return &Models{
		Movie: MovieModel{
			DB: conn,
		},
		Users: UserModel{
			DB: conn,
		},
		Token: TokenModel{
			DB: conn,
		},
	}, err
}
