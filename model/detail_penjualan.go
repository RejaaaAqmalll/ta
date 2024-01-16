package model

type DetailPenjualan struct {
	IdDetailPenjualan    string  `json:"id_detail_penjualan" gorm:"primary_key;type:varchar(255)"`
	PenjualanIdPenjualan string  `json:"penjualan_id_penjualan" gorm:"type:varchar(255)"`
	ProdukIdProduk       int     `json:"produk_id_produk" gorm:"type:int(5)"`
	JumlahProduk         int     `json:"jumlah_produk"`
	SubTotal             float64 `json:"sub_total"`
	Produk               *Produk `json:"produk" gorm:"foreignkey:id_produk;association_foreignkey:produk_id_produk;references:produk_id_produk"`
	BaseModel
}