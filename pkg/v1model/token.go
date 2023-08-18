package v1model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phluxx/gogert/internal/service/config"
)

func CreateToken(username string, cfg config.JwtConfig) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		// TODO: Pull other claims from LDAP
		"smell": "bad",
		"nbf":   time.Now().Add(time.Minute * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(cfg.Secret))
}
