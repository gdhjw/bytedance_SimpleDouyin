package model

import "github.com/jinzhu/gorm"

// Favorite 点赞
type Favorite struct {
	gorm.Model
	UserId  uint `json:"user_id"`
	VideoId uint `json:"video_id"`
	State   uint
}
