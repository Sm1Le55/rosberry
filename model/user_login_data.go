package model

type UserLoginData struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
