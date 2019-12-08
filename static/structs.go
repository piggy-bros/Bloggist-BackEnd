package static

import (
	"github.com/dgrijalva/jwt-go"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type BlogInfo struct {
	BlogID      int
	Author      string
	LikedNum    int
	BlogTitle   string
	BlogContent string
}
type PublishInfo struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var PrivateKey string = "private key"
