package main

import (
	"context"
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func initTables(db *sql.DB) {
	dbTables, err := ioutil.ReadFile("../db/tables.sql")

	if err != nil {
		log.Panic(err)
}

	log.Println("tables", string(dbTables))

	res, err := db.Exec(string(dbTables))
	if err != nil {
			log.Panic(err)
	}

	log.Println("res", res)
}

func insertData(db *sql.DB) {
	inserts, err := ioutil.ReadFile("../db/inserts.sql")

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

	log.Println("inserted all data...")
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


func blogPostHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Write([]byte(params["id"]))
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

	r := mux.NewRouter()

	r.HandleFunc("/api/blogPost/{id}", blogPostHandler).Methods("GET")

	address := "localhost:8098"

	srv := &http.Server{
		Addr: address,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	log.Println("Server started @", address, "...")

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	wait := time.Second*15

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
