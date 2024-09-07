package models

type Cred struct {
	Base
	UserId string `json:"userId"`
	Salt int `json:"-"`
	Password string `json:"password"`
}