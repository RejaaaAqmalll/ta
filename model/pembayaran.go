package model

type Pembayaran struct {
	Idpembayaran         string     `json:"idpembayaran" gorm:"primary_key;type:varchar(255)"`
	PenjualanIdPenjualan string     `json:"penjualan_id_penjualan" gorm:"type:varchar(255)"`
	Amount               float64    `json:"amount"`
	BiayaAdmin           float64    `json:"biaya_admin"`
	Grandtotal           float64    `json:"grandtotal"`
	Penjualan            *Penjualan `json:"penjualan" gorm:"foreignkey:PenjualanIdPenjualan;association_foreignkey:IdPenjualan;references:IdPenjualan"`
	BaseModel
}
