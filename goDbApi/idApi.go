package goDbApi

import (
	"log"
	"math/rand"
)

var raceIDs []int
var seasons []int
var driverIDs []int

func queryRaceIDs() {
	queryStmt := `select raceId from races`

	rows, err := DB.Query(queryStmt)

	if err != nil {
		log.Panicln("Cannot query raceIds", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				log.Panicln("Cannot scan raceId", err.Error())
			}
			raceIDs = append(raceIDs, id)
	}
}

func querySeasons() {
	queryStmt := `select year from seasons`

	rows, err := DB.Query(queryStmt)

	if err != nil {
		log.Panicln("Cannot query seasons", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
			var season int
			err := rows.Scan(&season)
			if err != nil {
				log.Panicln("Cannot scan season", err.Error())
			}
			seasons = append(seasons, season)
	}
}

func queryDriverIDs() {
	queryStmt := `select driverId from drivers`

	rows, err := DB.Query(queryStmt)

	if err != nil {
		log.Panicln("Cannot query drivers", err.Error())
	}

	defer rows.Close()

	for rows.Next() {
			var driver int
			err := rows.Scan(&driver)
			if err != nil {
				log.Panicln("Cannot scan driver", err.Error())
			}
			driverIDs = append(driverIDs, driver)
	}
}

func GetRandomRaceId () int {
	id := raceIDs[rand.Intn(len(raceIDs))]
	return id
}

func GetRandomSeason () int {
	year := seasons[rand.Intn(len(seasons))]
	return year
}

func GetRandomDriver () int {
	driverID := driverIDs[rand.Intn(len(driverIDs))]
	return driverID
}
