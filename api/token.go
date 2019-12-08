package api

import (
	"bloggist/static"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte(static.PrivateKey)

func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(3 * time.Hour)

	claims := static.Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "bloggist",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	return token, err
}

func ParseToken(token string) (*static.Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &static.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*static.Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
