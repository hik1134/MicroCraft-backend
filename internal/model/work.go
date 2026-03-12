package model

import "time"

type Work struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"work_id"`
	UserID       uint      `gorm:"not null;index" json:"user_id"`
	Title        string    `gorm:"type:varchar(100)" json:"title"`
	Type         string    `gorm:"type:varchar(20);not null" json:"type"`   
	Source       string    `gorm:"type:varchar(30);not null" json:"source"` 

	FileKey      string    `gorm:"type:varchar(255);not null" json:"-"`
	FileURL      string    `gorm:"type:varchar(500);not null" json:"file_url"`
	ThumbnailKey string    `gorm:"type:varchar(255)" json:"-"`
	ThumbnailURL string    `gorm:"type:varchar(500)" json:"thumbnail_url"`

	Status    string    `gorm:"type:varchar(20);not null;default:active" json:"status"` 
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}