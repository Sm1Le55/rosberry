package database

import (
	"database/sql"
	"fmt"
)

type AuthResult int

const (
	AuthResultSuccess AuthResult = iota + 1
	AuthResultUserNotFound
	AuthResultAccessKeyExpired
	AuthResultBadAccessKey
)

var db *sql.DB

func NewDB() *sql.DB {
	connStr := "user=postgres password=123 dbname=postgres search_path=rosberry_fsm sslmode=disable"
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Open db connect error: %v\n", err)
		return nil
	}
	db = DB
	return db
}
