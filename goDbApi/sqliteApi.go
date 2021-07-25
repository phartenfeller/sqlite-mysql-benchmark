package goDbApi

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func initSqlite() (*sql.DB, error)  {
	dbPath := os.Getenv("SQLITE_PATH")

	log.Println("Accessing DB from: ", dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
			return nil, errors.New("Error initializing SQLite" + err.Error())

	}

	log.Println("Initilized SQLite...")

	return db, nil
}
