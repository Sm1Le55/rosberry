package database

import (
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
	return nil
}

func ProfilesListQuery(r *model.ListRequest) ([]model.Profile, error) {
	return []model.Profile{}, nil
}
