package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mail.v2"
)


type OptGetEmail struct {
	NamaPenerima	string
	NamaKasir 		string
	TotalHarga 		float64
}

func SendEmail(email, imagePath string) error  {
	godotenv.Load()

	mailer := mail.NewMessage()
	mailer.SetHeader("From", os.Getenv("SENDER_NAME"))
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "      This is your transaction receipt      ")
	body := "<html><body style='text-align: center; color: white; font-size: 16px; font-weight: bold;'>Please read and check your receipt</body></html>"
	mailer.SetBody("text/html", body)
	// path get image receipt
	mailer.Attach(imagePath)

	host := os.Getenv("SMTP_HOST")
	port := 587
	authEmail := os.Getenv("AUTH_EMAIL")
	authPass := os.Getenv("AUTH_PASSWORD")

	dialer := mail.NewDialer(host, port, authEmail, authPass)

	if err := dialer.DialAndSend(mailer); err != nil {
		log.Fatal(err)
		return err
	}
	
	return nil
}
