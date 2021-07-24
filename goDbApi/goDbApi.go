package goDbApi

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB
var Posts int

func initTables() {
	dbTables, err := ioutil.ReadFile("./db/tables.sql")

	if err != nil {
		log.Panic(err)
}

	log.Println("tables", string(dbTables))

	res, err := DB.Exec(string(dbTables))
	if err != nil {
			log.Panic(err)
	}

	log.Println("res", res)
}

func insertData() {
	inserts, err := ioutil.ReadFile("./db/inserts.sql")

	if err != nil {
		log.Panic(err)
	}

	// log.Println("inserts", string(inserts))

	tx, err := DB.Begin()

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

func queryDataCount() {
	queryStmt := "SELECT count(*) FROM blog_posts"
	rows, err := DB.Query(queryStmt)
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
			Posts = int(count)
	}
}

// InitDb initializes the db
func InitDb() {
	dbPath := "simple.sqlite"

	os.Remove(dbPath)

	log.Println("Hello world")

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
			log.Panic(err)
	}

	DB = db

	initTables()
	insertData()
	queryDataCount()
}

// Post struct
type Post struct {
	PostID int `json:"postID"`
	Content string `json:"content"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	CreatedAt string `json:"createAt"`
}

func GetBlogpostById (id string) (Post, error) {
	var post Post;

	queryStmt := `SELECT post_id
	                   , content
	                   , title
	                   , slug
	                   , created_dt 
								  FROM blog_posts where post_id = ?`
	row := DB.QueryRow(queryStmt, id)

	err := row.Scan(&post.PostID, &post.Content, &post.Title, &post.Slug, &post.CreatedAt)
	if err != nil {
		return post, err
	}
	
	return post, nil
}
