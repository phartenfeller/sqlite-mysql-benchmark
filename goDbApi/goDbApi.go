package goDbApi

import (
	"database/sql"
	"errors"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDb initializes the db
func InitDb() {
	var db *sql.DB
	var err error

	if (os.Getenv("DB_DRIVER") == "PG") {
		db, err = initPg()
	} else {
		db, err = initSqlite()
	}

	if (err != nil) {
		log.Panic(err)
	}

	DB = db

	rand.Seed(time.Now().Unix())

	queryRaceIDs()
	querySeasons()
	queryDriverIDs()
}

// Post struct
type Post struct {
	PostID int `json:"postID"`
	Content string `json:"content"`
	Title string `json:"title"`
	Slug string `json:"slug"`
	CreatedAt string `json:"createAt"`
}

type DriverStandings struct {
	Forename string `json:"forename"`
	Surname string `json:"surname"`
	Points float32 `json:"points"`
	Position int `json:"position"`
	Wins int `json:"wins"`
}

func GetLastDriverStandingsByYear (year int) (dStandings []DriverStandings, err error) {
	queryStmt := `select d.forename, d.surname, points, position, wins
									from driverStandings ds
									join drivers d
										on ds.driverId = d.driverId
								where raceId = (
								select max(raceId)
										from (
										select distinct ra.raceId, ra.year
											from races ra
											join results re
												on ra.raceId = re.raceId
											where ra.year = $1
										) season_races
									group by year
						    )`

	rows, err := DB.Query(queryStmt, year)

	
	if err != nil {
		return nil, errors.New("Cannot query getLastDriverStandingsByYear: " + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
			var dStand DriverStandings
			err := rows.Scan(&dStand.Forename, &dStand.Surname, &dStand.Points, &dStand.Position, &dStand.Wins)
			if err != nil {
				return nil, errors.New("Cannot scan getLastDriverStandingsByYear: " + err.Error())
			}
			dStandings = append(dStandings, dStand)
	}

	return dStandings, nil
}


type RaceDriverAvgPitstop struct {
	Forename string `json:"forename"`
	Surname string `json:"surname"`
	Constructor string `json:"constructor"`
	Year int `json:"year"`
	Race string `json:"race"`
	AvgPitstopSeconds float32 `json:"avgPitsotopSeconds"`
}

func GetRaceDriverAvgPitstops (raceID int) (pstops []RaceDriverAvgPitstop, err error) {
	queryStmt := `select d.forename, d.surname, c.name as constructor, ra.year, ra.name as race, round(avg(p.milliseconds)) / 1000 as avg_pitstop_s
								from pitStops p
								join drivers d
									on p.driverId = d.driverId
								join results r
									on p.raceId = r.raceId
								and p.driverid = r.driverid
								join constructors c
									on r.constructorId = c.constructorId
								join races ra
									on ra.raceId = p.raceId
							where p.milliseconds < 180000 -- 3 minutes to filter out red flas
							  and p.raceId = $1
							group by p.raceId, d.forename, d.surname, c.name, ra.year, ra.name`

	rows, err := DB.Query(queryStmt, raceID)

	if err != nil {
		return nil, errors.New("Cannot query getRaceDriverAvgPitstops: " + err.Error())
	}
		defer rows.Close()
		for rows.Next() {
			var pstop RaceDriverAvgPitstop
			err := rows.Scan(&pstop.Forename, &pstop.Surname, &pstop.Constructor, &pstop.Year, &pstop.Race, &pstop.AvgPitstopSeconds)
			if err != nil {
				return nil, errors.New("Cannot scan getRaceDriverAvgPitstops: " + err.Error())
			}
			pstops = append(pstops, pstop)
	}
		return pstops, nil		
}

type AvgLapTime struct {
	Race string `json:"race"`
	Forename string `json:"forename"`
	Surname string `json:"surname"`
	AvgLapTime float32 `json:"avgLapTime"`
	RelevantLapCount int `json:"relevantLapCount"`
}

func GetAvgBestLapTimes (raceID int) (lapTimes []AvgLapTime, err error) {
	queryStmt := `
			select race
				, forename
				, surname
				, round(avg(milliseconds)) / 1000 as avg_lapTime_s
				, round(lap_count * 0.7) as relevant_lap_count
			from (
			select l.raceId
					, l.milliseconds
					, d.forename
					, d.surname
					, r2.name as race
					, row_number() over (partition by l.raceId, l.driverId order by l.milliseconds) as lap_rank
					, lp.lap_count
				from lapTimes l
				join drivers d
					on l.driverId = d.driverId
				join results r
					on l.raceId = r.raceId
				and l.driverId = r.driverId
				join races r2 on r.raceid = r2.raceid
				cross join (select max(lap) as lap_count from lapTimes il where il.raceId = $1) as lp
			where l.milliseconds < (select min(il.milliseconds) * 1.5 from lapTimes il where il.raceId = l.raceId)
				and l.raceid = $1
		) lapTimes
		where lap_rank <= round(lap_count * 0.7) --70% best times
		group by race, forename, surname, lap_count
		having count(*) >= round(lap_count * 0.7) --filter out drivers that have less rounds than 70%
	`

	rows, err := DB.Query(queryStmt, raceID)

	if err != nil {
		return nil, errors.New("Cannot query GetAvgBestLapTimes: " + err.Error())
	}
		defer rows.Close()
		for rows.Next() {
			var lapTime AvgLapTime
			err := rows.Scan(&lapTime.Race, &lapTime.Forename, &lapTime.Surname, &lapTime.AvgLapTime, &lapTime.RelevantLapCount)
			if err != nil {
				return nil, errors.New("Cannot scan GetAvgBestLapTimes: " + err.Error())
			}
			lapTimes = append(lapTimes, lapTime)
	}
		return lapTimes, nil		
}

type RaceDetails struct {
	Round int `json:"round"`
	Race string `json:"race"`
	Year string `json:"year"`
	Circuit string `json:"circuit"`
	Country string `json:"country"`
	Location string `json:"location"`
}

func GetRaceDetails (raceID int) (RaceDetails, error) {
	var details RaceDetails

	queryStmt := `
	 select r.round, r.name as race, r.year, c.name as circuit, c.country, c.location
   from races r
   join circuits c
     on r.circuitId = c.circuitId
	where r.raceId = $1
	`
	
	row := DB.QueryRow(queryStmt, raceID)

	err := row.Scan(&details.Round, &details.Race, &details.Year, &details.Circuit, &details.Country, &details.Location)
	if err != nil {
		return details, errors.New("Cannot scan GetRaceDetails: " + err.Error())
	}
	
	return details, nil
}
