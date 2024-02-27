package model

type Penjualan struct {
	IdPenjualan          string            `json:"id_penjualan" gorm:"primary_key;type:varchar(255)"`
	PelangganIdPelanggan string            `json:"pelanggan_id_pelanggan" gorm:"type:varchar(255)"`
	UserIduser           int               `json:"user_iduser"`
	TotalHarga           float64           `json:"total_harga"`
	User                 *User             `json:"user" gorm:"foreignkey:UserIduser;association_foreignkey:Iduser;references:Iduser"`
	Pelanggan            *Pelanggan        `json:"pelanggan" gorm:"foreignkey:PelangganIdPelanggan;association_foreignkey:IdPelanggan;references:IdPelanggan"`
	DetailPenjualan      []DetailPenjualan `json:"detail_penjualan" gorm:"foreignkey:PenjualanIdPenjualan;association_foreignkey:IdPenjualan;references:IdPenjualan"`
	BaseModel
}
