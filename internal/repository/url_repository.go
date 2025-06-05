package repository

import (
	"database/sql"
)

type URLRepository interface {
	StoreURL(short, long string) error
	GetLongURL(short string) (string, error)
}

type urlRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) StoreURL(short, long string) error {
	_, err := r.db.Exec("INSERT INTO urls (short, long) VALUES (?, ?)", short, long)
	return err
}

func (r *urlRepository) GetLongURL(short string) (string, error) {
	var longURL string
	err := r.db.QueryRow("SELECT long FROM urls WHERE short = ?", short).Scan(&longURL)
	return longURL, err
}
