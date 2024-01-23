package helper

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

func IsSupportedImageFormat(mime string) bool {
	return strings.HasPrefix(mime, "image/jpeg") || strings.HasPrefix(mime, "image/png") || strings.HasPrefix(mime, "image/jpg")
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GenerateFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	filename := fmt.Sprintf("%s%s", GenerateRandomString(10), ext)
	return filename
}


func GetImageSavePath(filename string) string  {
	saveDir := "./storage/foto"

	err := os.MkdirAll(saveDir, 0755)
    if err != nil {
		log.Fatal(err)
    }

	savePath := filepath.Join(saveDir, filename)

	return savePath
}