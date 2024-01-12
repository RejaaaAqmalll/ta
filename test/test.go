package main

import (
	"fmt"
	"ta-kasir/helper"
	"ta-kasir/model"

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
		// ... tambahkan field lain sesuai struktur model Anda
	}

	tokene, err := helper.GenerateToken(userData)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("JWT Token:", tokene)

}