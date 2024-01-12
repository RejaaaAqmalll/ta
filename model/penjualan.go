package model

type Penjualan struct {
	IdPenjualan          int     `json:"id_penjualan" gorm:"primary_key;index"`
	PelangganIdPelanggan int     `json:"pelanggan_id_pelanggan" gorm:"index"`
	TotalHarga           float64 `json:"total_harga"`
	BaseModel
}