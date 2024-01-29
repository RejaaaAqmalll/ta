package model

type Penjualan struct {
	IdPenjualan          string            `json:"id_penjualan" gorm:"primary_key;type:varchar(255)"`
	PelangganIdPelanggan string            `json:"pelanggan_id_pelanggan" gorm:"type:varchar(255)"`
	TotalHarga           float64           `json:"total_harga"`
	Pelanggan            *Pelanggan        `json:"pelanggan" gorm:"foreignkey:id_pelanggan;association_foreignkey:pelanggan_id_pelanggan;references:id_pelanggan"`
	DetailPenjualan      []DetailPenjualan `json:"detail_penjualan" gorm:"foreignkey:penjualan_id_penjualan;association_foreignkey:id_penjualan;references:id_penjualan"`
	BaseModel
}