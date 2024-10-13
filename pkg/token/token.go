package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/imabg/sync/internal/models"
	"github.com/imabg/sync/pkg/uuid"
)

type CustomClaim struct {
	Claims CustomClaimData
	jwt.RegisteredClaims
}

type CustomClaimData struct {
	UserId string `json:"userId"`
	Email string `json:"email"`
	Setting models.Setting `json:"setting"`
}

type NewToken struct {
	Claim CustomClaim
	SignedKey string
}

type TokenRes struct {
	Token string 
	ExpireAt time.Time
}

func New(key string) *NewToken {
	tCtx :=  &NewToken{
		Claim: CustomClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer: "Sync-platform",
				ID: uuid.GenerateUUID(),
			},
		},
		SignedKey: key,
	}
	return tCtx
}

func(t *NewToken) Generate(claimData CustomClaimData, expireIn time.Duration) (TokenRes, error) {
	expireTime := time.Now().Add(expireIn)
	t.Claim = CustomClaim{
		Claims: claimData,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t.Claim)
	ss, err := token.SignedString([]byte(t.SignedKey))
	if err != nil {
		return TokenRes{}, err
	}
	return TokenRes{Token: ss, ExpireAt: expireTime}, nil
}


func(t *NewToken) Validate(tokenString string) (*CustomClaim, error) {
    token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(t.SignedKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*CustomClaim); ok && token.Valid {
        return claims, nil
    } else {
        return nil, jwt.ErrInvalidKey
    }
}
