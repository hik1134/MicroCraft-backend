package model

import "time"

type Post struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"post_id"`
	WorkID    uint      `gorm:"not null;index" json:"work_id"`
	Section   string    `gorm:"type:varchar(30);not null;index" json:"section"` 
	Status    string    `gorm:"type:varchar(20);not null;index" json:"status"`  
	LikeCount int64     `gorm:"not null;default:0" json:"like_count"`
	ViewCount int64     `gorm:"not null;default:0" json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}