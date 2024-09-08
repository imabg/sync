package uuid

import (
	"github.com/matoous/go-nanoid/v2"
"github.com/google/uuid"
)


func GenerateUUID() string {
	return uuid.NewString()
}

// GenerateShortId take id length and return a unique id(default to 21)
func GenerateShortId(len int) string {
	if len == 0 {
		len = 8
	}
	id, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", len)
	if err != nil {
		panic(err)
	}
	return id
}