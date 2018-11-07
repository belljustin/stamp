package db

import (
	"database/sql"
	"fmt"

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
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}
