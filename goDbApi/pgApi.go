package goDbApi

import (
	"database/sql"
	"errors"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func initPg() (*sql.DB, error)  {
	db, err := sql.Open("pgx", os.Getenv("PG_DATABASE_URL"))
	if err != nil {
			return nil, errors.New("Error initializing Postgres" + err.Error())

	}

	db.SetMaxOpenConns(100)

	err = db.Ping()
	if err != nil {
		return nil, errors.New("Error initializing Postgres" + err.Error())
	}

	log.Println("Initilized Postgres...")

	return db, nil
}
