package model

import "github.com/jinzhu/gorm"

// Comment 评论
type Comment struct {
	ID         int64  `json:"comment_id,omitempty"`
	User       User   `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,"`
}
