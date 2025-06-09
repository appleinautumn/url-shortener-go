package storage

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"short" TEXT UNIQUE,
		"long" TEXT
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CloseDB closes the database connection.
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
