package database

import (
	"errors"
	"fmt"
	"rosberry/model"
)

func ProfilesListQuery(r *model.ListRequest) ([]model.Profile, error) {
	var userLocationSettings model.LocationType
	qUserLocationSettings := `
	SELECT 
		location
	FROM rosberry_fsm.locationsettings
	WHERE userID = $1
	LIMIT 1
	`
	err := db.QueryRow(qUserLocationSettings, r.UserId).Scan(&userLocationSettings)
	if err != nil {
		fmt.Printf("Error database query: %\n", err)
		return nil, err
	}

	fmt.Println("userLocationSettings: ", userLocationSettings)
	switch userLocationSettings {
	case 1:
		return queryProfilesWorld(r.UserId, r.Limit, r.Offset)
	case 2:
		return queryProfilesCountry(r.UserId, r.Limit, r.Offset)
	case 3:
		return queryProfilesNearby(r.UserId, r.Limit, r.Offset)
	}
	return nil, errors.New("Unexpected error")
}

func queryProfilesWorld(userID, limit, offset int) ([]model.Profile, error) {
	fmt.Println("World list query")
	q := `SELECT
			profile.userID,
			name,
			'data:image/png;base64,' || encode(COALESCE(photo, (select avatar from rosberry_fsm.empty limit 1)),'base64') as photo,
			birthday,
			(SELECT time FROM rosberry_fsm.authhistory a WHERE a.userID = profile.userID ORDER BY time LIMIT 1) as lastVisit,
			country
		FROM 
			rosberry_fsm.profile
		WHERE userID != $1
		ORDER BY lastVisit
		OFFSET $2
		`

	if limit == 0 {
		q += " LIMIT ALL"
	} else {
		q += fmt.Sprintf(" LIMIT %v", limit)
	}

	rows, err := db.Query(q, userID, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Profile, 0)
	for rows.Next() {
		var profile model.Profile
		err := rows.Scan(&profile.UserID,
			&profile.Name,
			&profile.Avatar,
			&profile.Birthday,
			&profile.LastVisit,
			&profile.Country)
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}

		interests := getUserInterestList(profile.UserID)
		profile.Interests = interests

		result = append(result, profile)
	}

	return result, nil
}

func getUserInterestList(userID int) []int {
	q := `SELECT theme
		FROM rosberry_fsm.userinterest
		WHERE userinterest.userID = $1
	`

	rows, err := db.Query(q, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make([]int, 0)
	for rows.Next() {
		var interestID int
		err := rows.Scan(&interestID)
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}

		result = append(result, interestID)
	}

	return result
}

func queryProfilesCountry(userID, limit, offset int) ([]model.Profile, error) {
	fmt.Println("Country list query")
	var userCountry string
	qUserCountry := `
		SELECT country
		FROM rosberry_fsm.profile
		WHERE profile.userID = $1
		LIMIT 1
	`
	err := db.QueryRow(qUserCountry, userID).Scan(&userCountry)
	if err != nil {
		fmt.Printf("Error database query: %v\n", err)
		return nil, err
	}

	q := `
		SELECT
			profile.userID,
			name,
			'data:image/png;base64,' || encode(COALESCE(photo, (select avatar from rosberry_fsm.empty limit 1)),'base64') as photo,
			birthday,
			auth.time as lastVisit,
			coord
		FROM 
			rosberry_fsm.profile,
			(
				SELECT
					ID,
					USERID,
					time,
					coord
				FROM
				rosberry_fsm.authhistory
				WHERE time = 
				(SELECT a.t FROM (SELECT userID as u, MAX(TIME) as t FROM rosberry_fsm.authhistory GROUP BY userID) a where a.u = userID)
			) auth
		WHERE 
		auth.UserID = profile.UserID AND
		country = $3 AND
		profile.userID != $1
		ORDER BY lastVisit
		OFFSET $2
	`

	if limit == 0 {
		q += " LIMIT ALL"
	} else {
		q += fmt.Sprintf(" LIMIT %v", limit)
	}

	rows, err := db.Query(q, userID, offset, userCountry)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Profile, 0)
	for rows.Next() {
		var profile model.Profile
		err := rows.Scan(&profile.UserID,
			&profile.Name,
			&profile.Avatar,
			&profile.Birthday,
			&profile.LastVisit,
			&profile.Country)
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}

		interests := getUserInterestList(profile.UserID)
		profile.Interests = interests

		result = append(result, profile)
	}

	if len(result) >= limit && limit != 0 {
		//Return, if count profiles ok
		fmt.Printf("Limit ok: %v (%v)\n", len(result), limit)
		return result, nil
	}

	if limit != 0 {
		limit = limit - len(result)
	}

	//We get entries from the general list
	if len(result) != 0 {
		offset = 0
	}

	additResult, err := queryProfilesWorld(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	result = append(result, additResult...)

	return result, nil
}

func queryProfilesNearby(userID, limit, offset int) ([]model.Profile, error) {
	fmt.Println("Nearby list query")
	type UserCoord struct {
		X float32
		Y float32
	}
	var userCoord UserCoord
	qUserCoord := `
		SELECT coord[0], coord[1]
		FROM rosberry_fsm.authhistory
		WHERE authhistory.userID = $1
		ORDER BY time desc
		LIMIT 1
	`
	err := db.QueryRow(qUserCoord, userID).Scan(&userCoord.X, &userCoord.Y)
	if err != nil {
		fmt.Printf("Error database query: %v\n", err)
		return nil, err
	}

	q := `
	SELECT
		profile.userID,
		name,
		'data:image/png;base64,' || encode(COALESCE(photo, (select avatar from rosberry_fsm.empty limit 1)),'base64') as photo,
		birthday,
		time as lastVisit,
		country
	FROM 
	rosberry_fsm.profile,
	(
		SELECT
			ID,
			USERID,
			TIME,
			coord
		FROM
		rosberry_fsm.authhistory
		WHERE TIME = 
		(SELECT a.t FROM (SELECT userID as u, MAX(TIME) as t FROM rosberry_fsm.authhistory GROUP BY userID) a where a.u = userID)
	) auth
	WHERE 
	auth.UserID = profile.UserID AND
	profile.userID != $1 AND
	rosberry_fsm.pointDistance(coord, point($3, $4)) < $5
	ORDER BY lastVisit
	OFFSET $2
	`

	if limit == 0 {
		q += " LIMIT ALL"
	} else {
		q += fmt.Sprintf(" LIMIT %v", limit)
	}

	nerabyRadiusKM := 150
	rows, err := db.Query(q, userID, offset, userCoord.X, userCoord.Y, nerabyRadiusKM)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]model.Profile, 0)
	for rows.Next() {
		var profile model.Profile
		err := rows.Scan(&profile.UserID,
			&profile.Name,
			&profile.Avatar,
			&profile.Birthday,
			&profile.LastVisit,
			&profile.Country)
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}

		interests := getUserInterestList(profile.UserID)
		profile.Interests = interests

		result = append(result, profile)
	}

	//Add country
	if len(result) >= limit && limit != 0 {
		//Return, if count profiles ok
		fmt.Printf("Limit ok: %v (%v)\n", len(result), limit)
		return result, nil
	}

	if limit != 0 {
		limit = limit - len(result)
	}

	if len(result) != 0 {
		offset = 0
	}
	//We get entries from the general list
	additResultCountry, err := queryProfilesCountry(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	result = append(result, additResultCountry...)

	//Add world
	if len(result) >= limit && limit != 0 {
		//Return, if count profiles ok
		fmt.Printf("Limit ok: %v (%v)\n", len(result), limit)
		return result, nil
	}

	if limit != 0 {
		limit = limit - len(result)
	}

	if len(result) != 0 {
		offset = 0
	}
	//We get entries from the general list
	additResultWorld, err := queryProfilesWorld(userID, limit, offset)
	if err != nil {
		return nil, err
	}
	result = append(result, additResultWorld...)

	return result, nil
}
