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

func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection.
func CloseDB() {
	if db != nil {
		db.Close()
	}
}
