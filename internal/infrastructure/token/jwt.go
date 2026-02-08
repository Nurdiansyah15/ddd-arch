package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	Secret string
	TTL    time.Duration
}

func (j *JWTMaker) Generate(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.TTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}
