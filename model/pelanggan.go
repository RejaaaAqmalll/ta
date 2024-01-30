package model

type Pelanggan struct {
	IdPelanggan string `json:"id_pelanggan" gorm:"primary_key;type:varchar(255)"`
	Email       string `json:"email" gorm:"type:varchar(255)"`
	Nama        string `json:"nama"`
	NoTelp      string `json:"no_telp"`
	Alamat      string `json:"alamat" gorm:"type:text"`
	BaseModel
}