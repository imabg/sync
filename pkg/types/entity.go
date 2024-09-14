package types

import "time"

type EntityResp struct {
	UserId string `json:"userId"`
	Email string `json:"email"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LoginResp struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt time.Time `json:"expireAt"`
}

type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}