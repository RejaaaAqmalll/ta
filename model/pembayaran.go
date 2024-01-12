package model

type Pembayaran struct {
	Idpembayaran int     `json:"idpembayaran" gorm:"primary_key"`
	Amount       float64 `json:"amount"`
	BiayaAdmin   float64 `json:"biaya_admin"`
	GrandTotal   float64 `json:"grandtotal"`
	BaseModel
}