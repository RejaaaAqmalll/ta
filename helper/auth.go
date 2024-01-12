package helper

import (
	"errors"
	"fmt"
	"os"
	"ta-kasir/model"
	"ta-kasir/model/request"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func GenerateToken(data model.User) (string, error) {
	godotenv.Load()

	secret := []byte(os.Getenv("SECRET"))
	fmt.Println(string(secret))
	if len(secret) == 0 {
		return "", errors.New("JWT_SECRET is empty")
	}
	issuedAt := time.Now().Unix()

	claims := request.JwtClaim{
		UserId: data.Iduser,
		Nama: data.Username,
		Role: data.Role,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: issuedAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString(secret)
	fmt.Println(string(secret))

	if err != nil {
		return "", err	
	}

	return tokenString, nil
}