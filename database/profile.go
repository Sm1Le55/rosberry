package database

import (
	"errors"
	"fmt"
	"rosberry/model"
)

func ProfileQuery(userID int) (*model.Profile, error) {
	profile := model.Profile{}

	qProfile := `SELECT userID,name,photo,birthday
	FROM rosberry_fsm.Profile WHERE userID = $1`

	err := db.QueryRow(qProfile, userID).Scan(&profile.UserID, &profile.Name, &profile.Avatar, &profile.Birthday)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return nil, err
	}

	profile.Interests = interestsQuery(userID)

	return &profile, nil
}

func interestsQuery(userID int) []int {
	q := `SELECT theme
	FROM rosberry_fsm.UserInterest, rosberry_fsm.Profile
	WHERE profile.userID =userInterest.userID AND profile.userID = $1`

	rows, err := db.Query(q, userID)
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

func UpdateProfile(profile *model.Profile) error {
	err := updProfileQuery(profile)
	if err != nil {
		return err
	}

	err = updUserIntr(profile.UserID, profile.Interests)
	if err != nil {
		return err
	}

	return nil
}

func updProfileQuery(profile *model.Profile) error {
	q := `INSERT INTO rosberry_fsm.profile (userID, name, photo, birthday)
			VALUES	($1, $2, $3, $4) ON CONFLICT (userID) DO UPDATE SET (name, photo, birthday) = ($2, $3, $4);`

	_, err := db.Exec(q, profile.UserID, profile.Name, profile.Avatar, profile.Birthday) //Strings! must be ints
	if err != nil {
		return errors.New("Profile update error: " + err.Error())
	}
	return nil
}

func updUserIntr(userID int, intrs []int) error {
	qDel := "DELETE FROM UserInterest WHERE userID = $1"
	_, err := db.Exec(qDel, userID)
	if err != nil {
		return errors.New("Update show interest error (del): " + err.Error())
	}

	qIns := "INSERT INTO UserInterest (userID, theme) VALUES ($1, $2)"
	for _, theme := range intrs {
		_, err := db.Exec(qIns, userID, theme)
		if err != nil {
			return errors.New("Update show interest error (ins): " + err.Error())
		}
	}

	return nil
}

func ProfilesListQuery(r *model.ListRequest) ([]model.Profile, error) {
	return []model.Profile{}, nil
}
