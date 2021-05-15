package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func initTables(db *sql.DB) {
	dbTables, err := ioutil.ReadFile("./db/tables.sql")

	log.Println("tables", string(dbTables))

	res, err := db.Exec(string(dbTables))
	if err != nil {
			log.Panic(err)
	}

	log.Println("res", res)
}

func insertData(db *sql.DB) {
	inserts, err := ioutil.ReadFile("./db/inserts.sql")

	if err != nil {
		log.Panic(err)
	}

	// log.Println("inserts", string(inserts))

	tx, err := db.Begin()

	if err != nil {
		log.Panic(err)
	}

	res2, err := tx.Exec(string(inserts))
  if err != nil {
      log.Panic(err)
  }

	log.Println("res2", res2)

	err = tx.Commit()

	if err != nil {
		log.Panic(err)
	}
}

func queryDataCount(db *sql.DB) {
	queryStmt := "SELECT count(*) FROM blog_posts"
	rows, err := db.Query(queryStmt)
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
			var count int64
			err = rows.Scan(&count)
			if err != nil {
					log.Panic(err)
			}
			log.Println("Rowcount:", count)
	}
}


func main() {
	dbPath := "simple.sqlite"

	os.Remove(dbPath)

	log.Println("Hello world")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
			log.Panic(err)
	}

	initTables(db)
	insertData(db)
	queryDataCount(db)
}
