package database

import (
	"rosberry/model"
	"database/sql"
	"fmt"
	"time"
	"errors"
	"crypto/rand"

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
 
	err := db.QueryRow("SELECT accessKey, accessKeyExpireDate FROM users WHERE ID = $1", userID).Scan(&validAccessKey, &accessKeyDateExpired)
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

func Registration(data model.UserRegData) error {
	if !emailValidation(data.Email) {
		return errors.New("Not valid email")
	}

	if getUserID(data.Email) != 0 {
		return errors.New("User alredy exists")
	}

	fmt.Printf("Email: %v passw: %v\n", data.Email, data.Password)
	_, err := db.Exec("INSERT INTO rosberry_fsm.users (email, password) VALUES ($1, $2)", data.Email, data.Password)
	if err != nil {
		return errors.New(fmt.Sprintf("Query error: %v\n", err))
	}
	return nil
}

func emailValidation(email string) bool {
	return true
}

//userID - search user with email in db, return user id
func getUserID(email string) int {
	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return 0
	}
	return userID
}

func Login(data model.UserLoginData) (model.UserAuthInfo, error) {
	if !emailValidation(data.Email) {
		return model.UserAuthInfo{}, errors.New("Not valid email")
	}

	userID := getUserID(data.Email)
	if userID == 0{
		return model.UserAuthInfo{}, errors.New("User not exists")
	}

	if !checkPassword(data.Email, data.Password) {
		return model.UserAuthInfo{}, errors.New("Wrong password")
	}

	//update coord
	_, err := db.Exec("INSERT AuthHistory SET latitude = $1, longitude =$2 where userID = $3", data.Latitude, data.Longitude, userID)
	if err != nil {
		return model.UserAuthInfo{}, errors.New("Auth history update error: " + err.Error())
	}

	//set access key
	accessKey := generateAccessKey()
	accessKeyExpireDate := time.Now().Add(time.Hour*24*7)
	_, err = db.Exec("UPDATE users SET accessKey = $1, accessKeyExpireDate =$2 where email = $3", accessKey, accessKeyExpireDate, data.Email)
	if err != nil {
		return model.UserAuthInfo{}, errors.New("Key issue error: " + err.Error())
	}

	return model.UserAuthInfo{userID, accessKey}, nil
}

func checkPassword(email string, pass string) bool {
	var validPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email = $1", email).Scan(&validPassword)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return false
	}
	if validPassword == "" {
		fmt.Println("Empty password")
		return false
	}
	if validPassword == pass {
		return true
	}
	return false
}

func generateAccessKey() string {
	length := 4
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		fmt.Println("Error key generate")	
	}

	fmt.Printf("Generate key: %x\n", key)
	return fmt.Sprintf("%x", key)
}


func ThemesListQuery() ([]model.Theme, error) {
	rows, err := db.Query("SELECT ID, title FROM themes")
	if err != nil {
		return nil, errors.New("Themes query error: " + err.Error())
	}
	defer rows.Close()

	result := make([]model.Theme, 0)
	for rows.Next() {
		theme := model.Theme{}
		err := rows.Scan(&theme.Id, &theme.Title)
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}
		result = append(result, theme)
	}

	return result, nil
}

