package model

import "time"

type Rafiki struct {
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(150);"  json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(150);"  json:"updated_by"`
}
