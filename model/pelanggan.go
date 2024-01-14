package model

type Pelanggan struct {
	IdPelanggan string `json:"id_pelanggan" gorm:"primary_key;index"`
	Nama        string `json:"nama"`
	NoTelp      string `json:"no_telp"`
	Alamat      string `json:"alamat" gorm:"type:text"`
	BaseModel
}