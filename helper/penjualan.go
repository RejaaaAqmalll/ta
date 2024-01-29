package helper

import (
	"math/rand"
	"time"
)

func GenerateRandomNumber(length int) string {
	const charset = "0123456789"

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	return string(b)
}