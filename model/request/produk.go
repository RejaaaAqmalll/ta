package request

type AddProduk struct {
	NamaProduk string  `json:"nama_produk" form:"nama_produk"`
	Harga      float64 `json:"harga" form:"harga"`
	Stok       int     `json:"stok" form:"stok"`
}

type EditProduk struct {
	NamaProduk string  `json:"nama_produk" form:"nama_produk"`
	Harga      float64 `json:"harga" form:"harga"`
	Stok       int     `json:"stok" form:"stok"`
}
