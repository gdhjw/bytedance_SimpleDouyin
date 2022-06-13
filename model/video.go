package model

import (
	"github.com/jinzhu/gorm"
)

// Video 视频
type Video struct {
	gorm.Model
	AuthorId      uint   `json:"author"`         // 所属用户id
	PlayUrl       string `json:"play_url"`       // 播放地址
	CoverUrl      string `json:"cover_url"`      // 封面地址
	FavoriteCount uint   `json:"favorite_count"` // 点赞数
	CommentCount  uint   `json:"comment_count"`  // 评论数
	IsFavorite    bool   `json:"is_favorite"`    // true 已点赞  false 未点赞
	Title         string `json:"title"`          // 标题
}
