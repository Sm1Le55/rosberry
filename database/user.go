package database

import (
	"crypto/rand"
	"errors"
	"fmt"
	"rosberry/model"
	"time"
)

type AgeRange string

const (
	AgeRangeAll     AgeRange = "all"
	AgeRangeBt18n24 AgeRange = "18-24"
	AgeRangeBt25n40 AgeRange = "25-40"
	AgeRangeMore40  AgeRange = "40+"
)

func AuthQuery(userID int, accessKey string) AuthResult {
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
	if validAccessKey != accessKey {
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

func Login(data model.UserLoginData) (*model.UserAuthInfo, error) {
	if !emailValidation(data.Email) {
		return nil, errors.New("Not valid email")
	}

	userID := getUserID(data.Email)
	if userID == 0 {
		return nil, errors.New("User not exists")
	}

	if !checkPassword(data.Email, data.Password) {
		return nil, errors.New("Wrong password")
	}

	//update coord
	_, err := db.Exec("INSERT AuthHistory SET latitude = $1, longitude =$2 where userID = $3", data.Latitude, data.Longitude, userID)
	if err != nil {
		return nil, errors.New("Auth history update error: " + err.Error())
	}

	//set access key
	accessKey := generateAccessKey()
	accessKeyExpireDate := time.Now().Add(time.Hour * 24 * 7)
	_, err = db.Exec("UPDATE users SET accessKey = $1, accessKeyExpireDate =$2 where userID = $3", accessKey, accessKeyExpireDate, userID)
	if err != nil {
		return nil, errors.New("Key issue error: " + err.Error())
	}

	return &model.UserAuthInfo{userID, accessKey}, nil
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

func Logout(email string) error {
	//!!TODO: logout only current user
	if !emailValidation(email) {
		return errors.New("Not valid email")
	}

	userID := getUserID(email)
	if userID == 0 {
		return errors.New("User not exists")
	}

	_, err := db.Exec("UPDATE users SET accessKeyExpireDate = $1 where id = $2", time.Now(), userID)
	if err != nil {
		return errors.New("Key issue error: " + err.Error())
	}

	return nil
}

func DisplaySettingsQuery(userID int) (*model.DisplaySettings, error) {
	showRange, hideRange := ageSettingsQuery(userID)
	showIntr := showIntrQuery(userID)
	hideIntr := hideIntrQuery(userID)
	location := locationQuery(userID)
	result := model.DisplaySettings{
		UserId:         userID,
		ShowMeAges:     showRange,
		HideMeFromAges: hideRange,
		ShowThemesID:   showIntr,
		HideThemesID:   hideIntr,
		Location:       location,
	}

	return &result, nil
}

func ageSettingsQuery(userID int) (string, string) {
	qAge := "SELECT userID, showRangeForMe, hideMeByRange" +
		"FROM rosberry_fsm.profile, rosberry_fsm.AgeSettings" +
		"WHERE AgeSettings.profile = profile.ID AND profile.userID = $1"

	var showRange, hideRange string
	err := db.QueryRow(qAge, userID).Scan(&showRange, &hideRange)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return "", ""
	}
	return showRange, hideRange
}

func showIntrQuery(userID int) []int {
	qShowIntr := "SELECT theme" +
		"FROM rosberry_fsm.profile, rosberry_fsm.ShowInterestsSettings" +
		"WHERE ShowInterestsSettings.profile = profile.ID AND profile.userID = $1"

	rows, err := db.Query(qShowIntr, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]int, 0)
	for rows.Next() {
		var intrID int
		err := rows.Scan(&intrID)
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}
		result = append(result, intrID)
	}
	return result
}

func hideIntrQuery(userID int) []int {
	qHideIntr := "SELECT theme" +
		"FROM rosberry_fsm.profile, rosberry_fsm.HideInterestsSettings" +
		"WHERE HideInterestsSettings.profile = profile.ID AND profile.userID = $1"

	rows, err := db.Query(qHideIntr, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]int, 0)
	for rows.Next() {
		var intrID int
		err := rows.Scan(&intrID)
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}
		result = append(result, intrID)
	}
	return result
}

func locationQuery(userID int) string {
	q := "SELECT title" +
		"FROM rosberry_fsm.profile, rosberry_fsm.locationSettings, rosberry_fsm.locations" +
		"WHERE locationSettings.profile = profile.ID and" +
		"locationSettings.location = locations.ID and profile.userID = $1"

	var location string
	err := db.QueryRow(q, userID).Scan(&location)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return ""
	}
	return location
}
