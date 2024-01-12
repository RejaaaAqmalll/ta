package model

type Penjualan struct {
	IdPenjualan          int               `json:"id_penjualan" gorm:"primary_key;index"`
	PelangganIdPelanggan int               `json:"pelanggan_id_pelanggan" gorm:"index"`
	TotalHarga           float64           `json:"total_harga"`
	DetailPenjualan      []DetailPenjualan `json:"detail_penjualan" gorm:"foreignkey:penjualan_id_penjualan;association_foreignkey:id_penjualan;references:id_penjualan"`
	BaseModel
}