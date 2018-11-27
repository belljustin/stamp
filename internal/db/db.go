package db

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/belljustin/stamp/internal/configs"
)

type Handle interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func InitDB(config *configs.DatabaseConfig) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)
	db, err := exponentialBackoff("postgres", connStr, 3)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	return db
}

func exponentialBackoff(driverName, dataSourceName string, maxAttempts uint) (*sql.DB, error) {
	db, err := connect(driverName, dataSourceName)
	if err == nil {
		return db, nil
	}
	for i := uint(1); i < maxAttempts; i++ {
		d := math.Exp2(float64(i))
		log.Printf("Could not connect to database. Will retry in %f seconds", d)
		time.Sleep(time.Duration(d) * time.Second)
		db, err = connect(driverName, dataSourceName)
		if err == nil {
			return db, nil
		}
	}
	return nil, err
}

func connect(driverName, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil { // Configuration error return right away
		return nil, err
	}
	err = db.Ping()
	if err == nil {
		return db, nil
	}
	return nil, err
}
