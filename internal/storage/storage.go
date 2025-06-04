package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	db, err = sql.Open("sqlite", dataSourceName)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"short" TEXT UNIQUE,
		"long" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	return err
}

// CloseDB closes the database connection.
func CloseDB() {
	if db != nil {
		db.Close()
	}
}

// StoreURL saves a short URL and its corresponding long URL to the database.
func StoreURL(short, long string) error {
	_, err := db.Exec("INSERT INTO urls (short, long) VALUES (?, ?)", short, long)
	return err
}

// GetLongURL retrieves the long URL for a given short ID.
// It returns sql.ErrNoRows if the short ID is not found.
func GetLongURL(shortID string) (string, error) {
	var longURL string
	err := db.QueryRow("SELECT long FROM urls WHERE short = ?", shortID).Scan(&longURL)
	return longURL, err
}
