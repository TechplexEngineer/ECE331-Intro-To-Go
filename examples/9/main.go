package main

import (
	"database/sql"
	"fmt"
	"github.com/d2r2/go-i2c"
	logger "github.com/d2r2/go-logger"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

var log = logger.NewPackageLogger("main",
	logger.DebugLevel,
)

func main() {

	dataSource := "./data.db"
	db, err := StartDb(dataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = logger.ChangePackageLogLevel("i2c", logger.FatalLevel)
	if err != nil {
		log.Fatal(err)
	}

	// Create new connection to I2C bus
	dataBus, err := i2c.NewI2C(0x4d, 1)
	if err != nil {
		log.Fatal(err)
	}
	// Free I2C connection on exit
	defer dataBus.Close()

	// Here goes code specific for sending and reading data
	// to and from device connected via I2C bus, like:
	buf1 := make([]byte, 1)
	_, err = dataBus.ReadBytes(buf1)
	if err != nil {
		log.Fatal(err)
	}
	temp := float64(buf1[0])
	if temp > 127 {
		temp = temp - 256
	}
	temp = temp*9./5. + 32

	dt := time.Now()

	d := dt.Format("2006-01-02")
	t := dt.Format("15:04:05")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("insert into temp(date,time,temp) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(d, t, temp)
	if err != nil {
		log.Fatal(err)
	}
}

// StartDb attempts to open a sqlite database specified by dataSource
//
// Callers should call db.Close() to free resources when all database operations are complete.
func StartDb(dataSource string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, fmt.Errorf("unable to open database %s - %w", dataSource, err)
	}

	// make sure the database table exists
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction to create table - %w", err)
	}
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS "temp" (
		date TEXT
		time TEXT
		temp REAL
	)`)
	if err != nil {
		return nil, fmt.Errorf("error executing create table - %w", err)
	}
	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error comitting create table tranaction - %w", err)
	}
	return db, err
}
