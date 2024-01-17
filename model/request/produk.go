package request

type AddProduk struct {
	NamaProduk string  `json:"nama"`
	Harga      float64 `json:"harga"`
	Stok       int     `json:"stok"`
}