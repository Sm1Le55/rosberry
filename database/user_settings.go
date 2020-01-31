package database

import (
	"errors"
	"fmt"
	"rosberry/model"
)

func DisplaySettingsQuery(userID int) (*model.DisplaySettings, error) {
	showRange, hideRange := ageSettingsQuery(userID)
	showIntr := showIntrQuery(userID)
	hideIntr := hideIntrQuery(userID)
	location := locationQuery(userID)
	result := model.DisplaySettings{
		UserID:         userID,
		ShowMeAges:     showRange,
		HideMeFromAges: hideRange,
		ShowThemesID:   showIntr,
		HideThemesID:   hideIntr,
		Location:       location,
	}

	return &result, nil
}

func ageSettingsQuery(userID int) (model.AgeRange, model.AgeRange) {
	qAge := `SELECT showRangeForMe, hideMeByRange
				FROM rosberry_fsm.users, rosberry_fsm.AgeSettings
				WHERE AgeSettings.userID = users.ID AND users.ID = $1`

	var showRange, hideRange model.AgeRange
	err := db.QueryRow(qAge, userID).Scan(&showRange, &hideRange)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return 0, 0
	}
	return showRange, hideRange
}

func showIntrQuery(userID int) []int {
	qShowIntr := `SELECT theme
					FROM rosberry_fsm.users, rosberry_fsm.ShowInterestsSettings
					WHERE ShowInterestsSettings.userID = users.ID AND users.ID = $1`

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
	qHideIntr := `SELECT theme
					FROM rosberry_fsm.users, rosberry_fsm.HideInterestsSettings
					WHERE HideInterestsSettings.userID = users.ID AND users.ID = $1`

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

func locationQuery(userID int) model.LocationType {
	q := `SELECT location
			FROM rosberry_fsm.users, rosberry_fsm.locationSettings
			WHERE locationSettings.userID = users.ID and users.ID = $1`

	var location model.LocationType
	err := db.QueryRow(q, userID).Scan(&location)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return 0
	}
	return location
}

func SaveDisplaySettings(settings *model.DisplaySettings) error {
	err := updAgeSettings(settings)
	if err != nil {
		return err
	}

	err = updShowIntr(settings.UserID, settings.ShowThemesID)
	if err != nil {
		return err
	}

	err = updHideIntr(settings.UserID, settings.HideThemesID)
	if err != nil {
		return err
	}

	err = updLocationSettings(settings.UserID, settings.Location)
	if err != nil {
		return err
	}

	return nil
}

func updAgeSettings(sett *model.DisplaySettings) error {
	q := `INSERT INTO rosberry_fsm.AgeSettings (userID, showRangeForMe, hideMeByRange)
			VALUES	($1, $2, $3) ON CONFLICT (userID) DO UPDATE SET (showRangeForMe, hideMeByRange) = ($2, $3);`

	_, err := db.Exec(q, sett.UserID, sett.ShowMeAges, sett.HideMeFromAges) //Strings! must be ints
	if err != nil {
		return errors.New("Age range settings update error: " + err.Error())
	}
	return nil
}

func updShowIntr(userID int, intrs []int) error {
	qDel := "DELETE FROM ShowInterestsSettings WHERE userID = $1"
	_, err := db.Exec(qDel, userID)
	if err != nil {
		return errors.New("Update show interest error (del): " + err.Error())
	}

	qIns := "INSERT INTO ShowInterestsSettings (userID, theme) VALUES ($1, $2)"
	for _, theme := range intrs {
		_, err := db.Exec(qIns, userID, theme)
		if err != nil {
			return errors.New("Update show interest error (ins): " + err.Error())
		}
	}

	return nil
}

func updHideIntr(userID int, intrs []int) error {
	qDel := "DELETE FROM HideInterestsSettings WHERE userID = $1"
	_, err := db.Exec(qDel, userID)
	if err != nil {
		return errors.New("Update show interest error (del): " + err.Error())
	}

	qIns := "INSERT INTO HideInterestsSettings (userID, theme) VALUES ($1, $2)"
	for _, theme := range intrs {
		_, err := db.Exec(qIns, userID, theme)
		if err != nil {
			return errors.New("Update show interest error (ins): " + err.Error())
		}
	}

	return nil
}

func updLocationSettings(userID int, location model.LocationType) error {
	q := `INSERT INTO LocationSettings (userID, location)
			VALUES	($1, $2) ON CONFLICT (userID) DO UPDATE SET location = $2;`

	_, err := db.Exec(q, userID, location) //Strings! must be ints
	if err != nil {
		return errors.New("Location settings update error: " + err.Error())
	}
	return nil
}
