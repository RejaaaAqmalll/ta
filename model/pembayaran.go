package model

type Pembayaran struct {
	Idpembayaran         int        `json:"idpembayaran" gorm:"primary_key"`
	PenjualanIdPenjualan int        `json:"penjualan_id_penjualan" gorm:"index"`
	Amount               float64    `json:"amount"`
	BiayaAdmin           float64    `json:"biaya_admin"`
	GrandTotal           float64    `json:"grandtotal"`
	Penjualan            *Penjualan `json:"penjualan" gorm:"foreignkey:id_penjualan;association_foreignkey:penjualan_id_penjualan;references:id_penjualan"`
	BaseModel
}