package models

type Cred struct {
	UserId string `json:"userId"`
	Salt int `json:"-"`
	Password string `json:"password"`
}