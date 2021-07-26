package main

import (
	"context"
	"encoding/json"
	"log"
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

func logError(err error) {
	log.Println("[error] " + err.Error())
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

func driverResultsYearHandler(w http.ResponseWriter, r *http.Request) {
	year := goDbApi.GetRandomSeason()
	standings, err := goDbApi.GetLastDriverStandingsByYear(year)

	if err != nil {
		logError(err)
		SendJSONErrorLog(w, "Unexpected server error - ", http.StatusInternalServerError, "Cannot query driver standings: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(standings)
}

func avgPitstopsHandler(w http.ResponseWriter, r *http.Request) {
	raceID := goDbApi.GetRandomRaceId()
	pstops, err := goDbApi.GetRaceDriverAvgPitstops(raceID)

	if err != nil {
		logError(err)
		SendJSONErrorLog(w, "Unexpected server error - ", http.StatusInternalServerError, "Cannot query pitstops: "+err.Error())
		return
	}

	if (len(pstops) == 0) {
		SendJSONErrorLog(w, "No pitstops found", http.StatusOK, "")
		return
	}

	json.NewEncoder(w).Encode(pstops)
}


func avgLapTimesHandler(w http.ResponseWriter, r *http.Request) {
	raceID := goDbApi.GetRandomRaceId()
	lapTimes, err := goDbApi.GetAvgBestLapTimes(raceID)

	if err != nil {
		logError(err)
		SendJSONErrorLog(w, "Unexpected server error - ", http.StatusInternalServerError, "Cannot query avg laptimes: "+err.Error())
		return
	}

	if (len(lapTimes) == 0) {
		SendJSONErrorLog(w, "No times found", http.StatusOK, "")
		return
	}

	json.NewEncoder(w).Encode(lapTimes)
}

func raceDetailsHandler(w http.ResponseWriter, r *http.Request) {
	raceID := goDbApi.GetRandomRaceId()
	details, err := goDbApi.GetRaceDetails(raceID)

	if err != nil {
		logError(err)
		SendJSONErrorLog(w, "Unexpected server error - ", http.StatusInternalServerError, "Cannot query race details: "+err.Error())
		return
	}

	json.NewEncoder(w).Encode(details)
}


func randomReadHandler(w http.ResponseWriter, r *http.Request) {
	reqcnt++
	log.Println("req", reqcnt);
	
	mod := reqcnt % 6

	switch mod {
		case 0:
			raceDetailsHandler(w, r)
		case 1:
			avgLapTimesHandler(w, r)
		case 2: 
			avgPitstopsHandler(w, r)
		case 3:
			driverResultsYearHandler(w, r)
		case 4:
			raceDetailsHandler(w, r)
		case 5:
			driverResultsYearHandler(w, r)
		default:
			log.Panicln("unhandled mod => ", mod)
	}
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

	r.HandleFunc("/api/driverResultsYear", driverResultsYearHandler).Methods("GET")
	r.HandleFunc("/api/avgPitsotps", avgPitstopsHandler).Methods("GET")
	r.HandleFunc("/api/avgLaptimes", avgLapTimesHandler).Methods("GET")
	r.HandleFunc("/api/raceDetails", raceDetailsHandler).Methods("GET")

	r.HandleFunc("/api/randomRead", randomReadHandler).Methods("GET")

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
