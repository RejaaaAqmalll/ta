package request

type Transaksi struct {
	Email       string    `json:"email" form:"email" binding:"email"`
	Nama        string    `json:"nama" form:"nama" binding:"required"`
	NoTelp      string    `json:"no_telp" form:"no_telp" binding:"required"`
	Alamat      string    `json:"alamat" form:"alamat" binding:"required"`
	DataPesanan []Pesanan `json:"data_pesanan" form:"data_pesanan" binding:"required"`
	Pembayaran  Bayar     `json:"pembayaran" form:"pembayaran" binding:"required"`
}

type Pesanan struct {
	IdProduk     int     `json:"id_produk" form:"id_produk" binding:"required"`
	NamaProduk   string  `json:"nama_produk" form:"nama_produk" binding:"required"`
	JumlahProduk int     `json:"jumlah_produk" form:"jumlah_produk" binding:"required"`
	SubTotal     float64 `json:"sub_total" form:"sub_total" binding:"required"`
}

type Bayar struct {
	Amount     float64 `json:"amount" form:"amount" binding:"required"`
	BiayaAdmin float64 `json:"biaya_admin" form:"biaya_admin" binding:"required"`
	Grandtotal float64 `json:"grandtotal" form:"grandtotal" binding:"required"`
}
