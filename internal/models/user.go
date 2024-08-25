package models

// User schema 
type User struct {
	Base
	Email     *string            `json:"email"`
	UserId    string             `json:"userId"`
	Name      *string            `json:"name"`
	Password  string             `json:"password"`
}
