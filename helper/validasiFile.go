package helper

import "strings"

func IsSupportedImageFormat(mime string) bool {
	return strings.HasPrefix(mime, "image/jpeg") || strings.HasPrefix(mime, "image/png") || strings.HasPrefix(mime, "image/jpg")
}