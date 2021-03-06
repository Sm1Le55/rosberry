package database

import (
	"errors"
	"fmt"
	"rosberry/model"

	"github.com/lib/pq"
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
		return queryProfilesWorld(r.UserId, []int{r.UserId}, r.Limit, r.Offset)
	case 2:
		return queryProfilesCountry(r.UserId, []int{r.UserId}, r.Limit, r.Offset)
	case 3:
		return queryProfilesNearby(r.UserId, []int{r.UserId}, r.Limit, r.Offset)
	}
	return nil, errors.New("Unexpected error")
}

func queryProfilesWorld(userID int, excludeIds []int, limit int, offset int) ([]model.Profile, error) {
	fmt.Println("World list query. Exclude: ", excludeIds)
	q := `SELECT 
		userID, name, photo, birthday, lastVisit, country, profileIntrst
		FROM
			(SELECT
				p.userID,
				name,
				'data:image/png;base64,' || encode(COALESCE(photo, (select avatar from rosberry_fsm.empty limit 1)),'base64') as photo,
				birthday,
				(SELECT time FROM rosberry_fsm.authhistory a WHERE a.userID = p.userID ORDER BY time LIMIT 1) as lastVisit,
				country,
				array(
					select theme
					from rosberry_fsm.userinterest t 
					where t.userID = p.userID
				) as profileIntrst,
				array(
					select theme
					from rosberry_fsm.hideinterestssettings t 
					where t.userID = $1
				) as notShowIntrst
			FROM 
				rosberry_fsm.profile p
			) t
			WHERE 
			NOT (userID = ANY ($2)) AND 
			NOT(profileIntrst && notShowIntrst)
			ORDER BY lastVisit
			OFFSET $3
		`

	if limit == 0 {
		q += " LIMIT ALL"
	} else {
		q += fmt.Sprintf(" LIMIT %v", limit)
	}

	rows, err := db.Query(q, userID, pq.Array(excludeIds), offset)
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
			&profile.Country,
			pq.Array(&profile.Interests))
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}

		result = append(result, profile)
	}

	return result, nil
}

func queryProfilesCountry(userID int, excludeIds []int, limit int, offset int) ([]model.Profile, error) {
	fmt.Println("Country list query. Exclude: ", excludeIds)
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
		userID, name, photo, birthday, lastVisit, country, profileIntrst
		FROM
		(SELECT
			profile.userID,
			name,
			'data:image/png;base64,' || encode(COALESCE(photo, (select avatar from rosberry_fsm.empty limit 1)),'base64') as photo,
			birthday,
			auth.time as lastVisit,
			country,
			array(
				select theme
				from rosberry_fsm.userinterest t 
				where t.userID = profile.userID
			) as profileIntrst,
			 array(
				select theme
				from rosberry_fsm.hideinterestssettings t 
				where t.userID = $2
			) as notShowIntrst
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
		NOT (profile.userID = ANY ($1))) t
		WHERE
		NOT(profileIntrst && notShowIntrst)
		ORDER BY lastVisit
		OFFSET $4
	`

	if limit == 0 {
		q += " LIMIT ALL"
	} else {
		q += fmt.Sprintf(" LIMIT %v", limit)
	}

	rows, err := db.Query(q, pq.Array(excludeIds), userID, userCountry, offset)
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
			&profile.Country,
			pq.Array(&profile.Interests))
		if err != nil {
			fmt.Printf("Error row scan: %v\n", err)
			continue
		}

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

	exIds := make([]int, 0, len(result)+1)
	exIds = append(exIds, userID)
	for _, val := range result {
		exIds = append(exIds, val.UserID)
	}
	additResult, err := queryProfilesWorld(userID, exIds, limit, offset)
	if err != nil {
		return nil, err
	}
	result = append(result, additResult...)

	return result, nil
}

func queryProfilesNearby(userID int, excludeIds []int, limit int, offset int) ([]model.Profile, error) {
	fmt.Println("Nearby list query. Exclude: ", excludeIds)
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
		userID, name, photo, birthday, lastVisit, country, profileIntrst
	FROM
	(SELECT
		profile.userID,
		name,
		'data:image/png;base64,' || encode(COALESCE(photo, (select avatar from rosberry_fsm.empty limit 1)),'base64') as photo,
		birthday,
		time as lastVisit,
		country,
		array(
			select theme
			from rosberry_fsm.userinterest t 
			where t.userID = profile.userID
		) as profileIntrst,
		 array(
			select theme
			from rosberry_fsm.hideinterestssettings t 
			where t.userID = $2
		) as notShowIntrst
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
	NOT (profile.UserID = ANY ($1)) AND
	rosberry_fsm.pointDistance(coord, point($3, $4)) < $5) t
	WHERE NOT(profileIntrst && notShowIntrst)
	ORDER BY lastVisit
	OFFSET $6
	`

	if limit == 0 {
		q += " LIMIT ALL"
	} else {
		q += fmt.Sprintf(" LIMIT %v", limit)
	}

	nerabyRadiusKM := 150
	rows, err := db.Query(q, pq.Array(excludeIds), userID, userCoord.X, userCoord.Y, nerabyRadiusKM, offset)
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
		/*
		interests := getUserInterestList(profile.UserID)
		profile.Interests = interests
*/
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
	exIds := make([]int, 0, len(result)+1)
	exIds = append(exIds, userID)
	for _, val := range result {
		exIds = append(exIds, val.UserID)
	}
	additResultCountry, err := queryProfilesCountry(userID, exIds, limit, offset)
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
	for _, val := range result {
		exIds = append(exIds, val.UserID)
	}
	additResultWorld, err := queryProfilesWorld(userID, exIds, limit, offset)
	if err != nil {
		return nil, err
	}
	result = append(result, additResultWorld...)

	return result, nil
}
