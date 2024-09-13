package types

import "time"

type EntityResp struct {
	UserId string `json:"userId"`
	Email string `json:"email"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type LoginDTO struct {
	Email string `json:"email"`
	Password string `json:"password"`
}