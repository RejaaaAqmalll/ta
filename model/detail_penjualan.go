package model

type DetailPenjualan struct {
	IdDetailPenjualan    string  `json:"id_detail_penjualan" gorm:"primary_key;index"`
	PenjualanIdPenjualan string  `json:"penjualan_id_penjualan" gorm:"index"`
	ProdukIdProduk       int     `json:"produk_id_produk" gorm:"index"`
	JumlahProduk         int     `json:"jumlah_produk"`
	SubTotal             float64 `json:"sub_total"`
	Produk               *Produk `json:"produk" gorm:"foreignkey:id_produk;association_foreignkey:produk_id_produk;references:produk_id_produk"`
	BaseModel
}