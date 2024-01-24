package request

type AddProduk struct {
	NamaProduk string  `json:"nama_produk" form:"nama_produk" binding:"required"`
	Harga      float64 `json:"harga" form:"harga" binding:"required"`
	Stok       int     `json:"stok" form:"stok" binding:"required"`
}

type EditProduk struct {
	NamaProduk string  `json:"nama_produk" form:"nama_produk" binding:"required"`
	Harga      float64 `json:"harga" form:"harga" binding:"required"`
	Stok       int     `json:"stok" form:"stok" binding:"required"`
}
