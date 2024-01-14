package main

import (
	"fmt"
	"os"
	"ta-kasir/helper"
	"ta-kasir/model"
	"ta-kasir/model/request"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)


func main() {
	e := godotenv.Load()
	if e!= nil {
		fmt.Println("Error:", e)
		return
	}
	userData := model.User{
		Iduser:   123,
		Username: "john_doe",
		Role:     1,
	}
	tokenString, err := helper.GenerateToken(userData)

	if err != nil{
		fmt.Println("Error:", err)
		return
	}

	parsedToken, err := jwt.ParseWithClaims(tokenString, &request.JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})
	
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return
	}
	
	if claims, ok := parsedToken.Claims.(*request.JwtClaim); ok && parsedToken.Valid {
		fmt.Println("IssuedAt from token claims:", claims.IssuedAt)
	} else {
		fmt.Println("Invalid token")
	}
}

	// tokene, err := helper.GenerateToken(userData)

	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// fmt.Println("JWT Token:", tokene)