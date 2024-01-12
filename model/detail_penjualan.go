package model

type DetailPenjualan struct {
	IdDetailPenjualan    int     `json:"id_detail_penjualan" gorm:"primary_key;index"`
	PenjualanIdPenjualan int     `json:"penjualan_id_penjualan" gorm:"index"`
	ProdukIdProduk       int     `json:"produk_id_produk" gorm:"index"`
	JumlahProduk         int     `json:"jumlah_produk"`
	SubTotal             float64 `json:"sub_total"`
}