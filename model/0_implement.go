package model

import "time"

type BaseModel struct {
	CreatedAt *time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP;INDEX"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	Hapus int `json:"hapus" gorm:"default:0;type:int(3)"`
}