package model

import "time"

type Carrier struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"carrier_id"`
	Name            string    `gorm:"type:varchar(50);not null" json:"name"`
	PreviewImageURL string    `gorm:"type:varchar(500);not null" json:"preview_image_url"`
	ModelCode       string    `gorm:"type:varchar(100);not null;uniqueIndex" json:"model_code"`
	Status          string    `gorm:"type:varchar(20);not null;default:active;index" json:"status"` 
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}