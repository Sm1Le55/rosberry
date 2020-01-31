package database

import "rosberry/model"

func ProfileQuery(userID int) (*model.Profile, error) {
	return nil, nil
}

func UpdateProfile(profile *model.Profile) error {
	return nil
}

func ProfilesListQuery(r *model.ListRequest) ([]model.Profile, error) {
	return nil, []model.Profile{}
}
