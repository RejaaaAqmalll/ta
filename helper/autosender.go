package helper

import (
	"fmt"
	"time"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "Admin Aplikasi Simple Cash <asimplecash@gmail.com>"
const CONFIG_AUTH_EMAIL = "asimplecash@gmail.com"
const CONFIG_AUTH_PASSWORD = "adminsimplecash2024"


// contoh template WA/Email
type optGetTextWa struct {
	NamaPenerima      string
	NamaOutlet        string
	NoInvoice         string
	Nominal           string
	TanggalJatuhTempo string
	Link              string
	TanggalPembayaran string
	DibayarMelalui    string
}

func getTextWa(opt optGetTextWa) string {	
	return fmt.Sprintf(
		"Invoice , Halo %s\n"+
			"Kami dari Laundry %s\n\n"+
			"berikut ini ada tagihan invoice yang perlu anda lunasi \n\n"+
			"nomer invoice : %s \n"+
			"tagihan : %s \n"+
			"tanggal jatuh tempo : %s\n"+
			"file invoice : %s \n\n"+
			"Terima kasih",
		opt.NamaPenerima,
		opt.NamaOutlet,
		opt.NoInvoice,
		opt.Nominal,
		opt.TanggalJatuhTempo,
		opt.Link,
	)
}

// contoh looping apabila pesan ada array nya 
type array struct{
	mulai string
	selesai string
}

var jadwalArray []array
func coba() string  {
	// var str string
	var str string
	var str2 []string
	for _, each := range jadwalArray {
		layout := "2006-01-02 15:04:05"
		jadwal := each.mulai
		jadwalParse, _ := time.Parse(layout, jadwal)
		jadwalFormat := jadwalParse.Format("02 - 01 - 2006, 15:04")
		str = ""
		str = "tanggal : " + str + jadwalFormat + "\n"
		str2 = append(str2, str)
	}
		return str
}