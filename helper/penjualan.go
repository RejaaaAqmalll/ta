package helper

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
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

func GenerateImage(width, height int, idPenjualan string) (string, error) {

	dc := gg.NewContext(width, height)

	// set background color
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetColor(color.Black)
	
	// HEADER 1
	err := dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 18)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("Simple Cash", 100, 30,)
	
	// HEADER 2
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 16)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("Jl Tanimbar No.22 Malang", 50, 55)

	// GARIS PEMISAH 1
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 10)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("-----------------------------------------------", 20, 75)
	
	// TANGGAL DAN NAMA KASIR
	timeString := time.Now().Format("02.01.2006")
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 11)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString(timeString + "/", 20, 95)
	dc.DrawString("", 84, 95)
	dc.DrawString("CASHIER 01", 220, 95)

	// GARIS PEMISAH 2
	err = dc.LoadFontFace("./storage/fonts/Poppins-Regular.ttf", 10)
	if err != nil {
		log.Fatal(err)
	}
	dc.DrawString("-----------------------------------------------", 20, 115)

	receiptPath := fmt.Sprintf("./storage/receipt/Receipt Simple Cash <%s>.png", idPenjualan) 
	err = dc.SavePNG(receiptPath)
	if err != nil {
		return "", err
	}
	return receiptPath, nil
}