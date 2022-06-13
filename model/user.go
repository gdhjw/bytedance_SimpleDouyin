package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name          string `json:"name"`
	Password      string `json:"password"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
}

type Users struct {
	Id            uint   `json:"id" db:"id"`
	Name          string `json:"name" db:"username"`
	Password      string `json:"password" db:"password"`
	FollowCount   uint   `json:"followCount" db:"follow_count"`
	FollowerCount uint   `json:"followerCount" db:"follower_count"`
}
