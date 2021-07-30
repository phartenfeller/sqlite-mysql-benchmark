package goDbApi

import (
	"errors"
	"math/rand"
	"strconv"
)

type LapTime struct {
	RaceID int `json:"raceId"`
	DriverID int `json:"driverId"`
	Lap int `json:"lap"`
	Position int `json:"position"`
	Time string `json:"time"`
	Milliseconds int `json:"milliseconds"`
}

func InsertLapTime() (l LapTime, err error) {
	l.RaceID = GetRandomRaceId()
	l.DriverID = GetRandomDriver()
	l.Lap = rand.Intn(90) + 100
	l.Position = rand.Intn(22)
	l.Time = "11:" + strconv.Itoa(rand.Intn(59)) + "." + strconv.Itoa(rand.Intn(999))
	l.Milliseconds = rand.Intn(999999)

	tx, err := DB.Begin()

	if err != nil {
		return l, errors.New("Cannot begin insert transaction: " + err.Error())
	}

	stmt := "INSERT INTO lapTimes (raceId, driverId, lap, position, time, milliseconds) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err = tx.Exec(stmt, l.RaceID, l.DriverID, l.Lap, l.Position, l.Time, l.Milliseconds)

	if err != nil {
		tx.Rollback()
		return l, errors.New("Cannot exec insert: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return l, errors.New("Cannot commit insert transaction: " + err.Error())
	}

	return l, err
}

func UpdateLaptime() (err error) {
	raceID := GetRandomRaceId()
	driverID := GetRandomDriver()
	lap := rand.Intn(20)

	tx, err := DB.Begin()

	if err != nil {
		return errors.New("Cannot begin update transaction: " + err.Error())
	}

	stmt := "UPDATE lapTimes SET milliseconds = milliseconds + 1, lap = lap * lap, position = position * 2 WHERE driverId = $1 and raceId = $2 and lap = $3"

	_, err = tx.Exec(stmt, driverID, raceID, lap)

	if err != nil {
		tx.Rollback()
		return errors.New("Cannot exec update: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("Cannot commit update transaction: " + err.Error())
	}

	return err
}

func DeleteLaptime() (err error) {
	raceID := GetRandomRaceId()
	driverID := GetRandomDriver()
	lap := rand.Intn(20)

	tx, err := DB.Begin()

	if err != nil {
		return errors.New("Cannot begin delete transaction: " + err.Error())
	}

	stmt := "DELETE from lapTimes WHERE driverId = $1 and raceId = $2 and lap = $3"

	_, err = tx.Exec(stmt, raceID, driverID, lap)

	if err != nil {
		tx.Rollback()
		return errors.New("Cannot exec update: " + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("Cannot commit delete transaction: " + err.Error())
	}

	return err
}
