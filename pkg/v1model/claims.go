package v1model

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT claims struct
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GetClaims(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}
	return claims, nil
}
