package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"dev.hartenfeller.sqlite-mysql-benchmark/goDbApi"
	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// SendJSONError returns an Error response in JSON format
func SendJSONError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	errResponse := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(errResponse)
}

// SendJSONErrorLog returns an Error response in JSON format
func SendJSONErrorLog(w http.ResponseWriter, message string, code int, errMessage string) {
	log.Println(errMessage)
	SendJSONError(w, message, code)
}

var reqcnt int;

func blogPostHandler(w http.ResponseWriter, r *http.Request) {
	reqcnt++
	log.Println("req", reqcnt);
	params := mux.Vars(r)

	post, err := goDbApi.GetBlogpostById(params["id"])

	if err != nil {
		SendJSONErrorLog(w, "Unexpected server error", http.StatusInternalServerError, "Cannot query blogpost: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(post)
}

func randomBlogPostHandler(w http.ResponseWriter, r *http.Request) {
	reqcnt++
	log.Println("req", reqcnt);

	id := rand.Intn(goDbApi.Posts - 1) + 1
	log.Println("id", id)

	post, err := goDbApi.GetBlogpostById(fmt.Sprint(id))

	if err != nil {
		SendJSONErrorLog(w, "Unexpected server error - " + fmt.Sprint(id), http.StatusInternalServerError, "Cannot query blogpost: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(post)
}


func StandardHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	reqcnt = 0

	goDbApi.InitDb()

	r := mux.NewRouter()
	r.Use(StandardHeadersMiddleware)

	r.HandleFunc("/api/blogPost/{id}", blogPostHandler).Methods("GET")
	r.HandleFunc("/api/randomBlogPost", randomBlogPostHandler).Methods("GET")

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
