package model

type Pembayaran struct {
	Idpembayaran         string     `json:"idpembayaran" gorm:"primary_key;type:varchar(255)"`
	PenjualanIdPenjualan string     `json:"penjualan_id_penjualan" gorm:"type:varchar(255)"`
	Amount               float64    `json:"amount"`
	BiayaAdmin           float64    `json:"biaya_admin"`
	GrandTotal           float64    `json:"grandtotal"`
	Penjualan            *Penjualan `json:"penjualan" gorm:"foreignkey:id_penjualan;association_foreignkey:penjualan_id_penjualan;references:id_penjualan"`
	BaseModel
}