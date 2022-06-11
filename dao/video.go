package dao

import "bytedance_SimpleDouyin/model"

// GetVideoList 获取当前用户的发布视频列表
func GetVideoList(userId uint) (int, []model.Video) {
	var videoList []model.Video
	rs := db.Table("videos").Where("author_id=?", userId).Find(&videoList)
	num := rs.RowsAffected
	return int(num), videoList
}
