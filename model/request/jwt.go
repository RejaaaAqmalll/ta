package request

import "github.com/dgrijalva/jwt-go"

type JwtClaim struct {
	UserId int `json:"userid"`
	Nama   string `json:"nama"`
	Role   int `json:"role"`
	jwt.StandardClaims
}