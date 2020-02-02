package database

import (
	"errors"
	"fmt"
	"rosberry/model"
)

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
