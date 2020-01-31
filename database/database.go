package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
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

func AuthQuery(userID int,accessKey string) AuthResult {
	var validAccessKey string
	var accessKeyDateExpired time.Time
 
	err := db.QueryRow("SELECT accessKey, accessKeyExpireDate FROM profile WHERE ID = $1", userID).Scan(&validAccessKey, &accessKeyDateExpired)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return AuthResultUserNotFound
	}

	fmt.Printf("Query return: %v %v\n", validAccessKey, accessKeyDateExpired)

	if validAccessKey == "" || accessKeyDateExpired.Before(time.Now()) {
		return AuthResultAccessKeyExpired
	}
	if validAccessKey != accessKey{
		return AuthResultBadAccessKey
	}
	return AuthResultSuccess
}
