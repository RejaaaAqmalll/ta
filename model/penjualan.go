package model

type Penjualan struct {
	IdPenjualan          string            `json:"id_penjualan" gorm:"primary_key;type:varchar(255)"`
	PelangganIdPelanggan string            `json:"pelanggan_id_pelanggan" gorm:"type:varchar(255)"`
	TotalHarga           float64           `json:"total_harga"`
	Pelanggan            *Pelanggan        `json:"pelanggan" gorm:"foreignkey:PelangganIdPelanggan;association_foreignkey:IdPelanggan;references:IdPelanggan"`
	DetailPenjualan      []DetailPenjualan `json:"detail_penjualan" gorm:"foreignkey:PenjualanIdPenjualan;association_foreignkey:IdPenjualan;references:IdPenjualan"`
	BaseModel
}