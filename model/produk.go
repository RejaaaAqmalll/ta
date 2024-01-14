package model

type Produk struct {
	IdProduk   int     `json:"id_produk" gorm:"primary_key;index"`
	NamaProduk string  `json:"nama_produk"`
	Harga      float64 `json:"harga"`
	Stok       int     `json:"stok"`
	Gambar     string  `json:"gambar" gorm:"type:text"`
	BaseModel
}