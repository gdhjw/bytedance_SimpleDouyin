package service

import (
	"bytedance_SimpleDouyin/dao"
	"bytedance_SimpleDouyin/model"
	"github.com/jinzhu/gorm"
)

// CheckFavorite 检查当前用户对当前视频是否点赞
func CheckFavorite(uid uint, vid uint) bool {
	var num int64
	dao.SqlSession.Table("favorites").Where("user_id = ? and video_id = ? and state = 1", uid, vid).Count(&num)

	if num == 0 {
		return false
	}
	return true
}

// ActionFavorite 点赞
func ActionFavorite(userId uint, videoId uint, actionType uint) (err error) {
	// actionType 1点赞，2取消点赞
	if actionType == 1 {
		newFavorite := model.Favorite{
			UserId:  userId,
			VideoId: videoId,
			State:   1,
		}
		var favoriteExist = &model.Favorite{} //找不到时会返回错误
		rs := dao.SqlSession.Table("favorites").Where("user_id = ? and video_id = ?", userId, videoId).First(favoriteExist)
		if rs.Error != nil { // 不存在
			if err := dao.SqlSession.Table("favorites").Create(&newFavorite).Error; err != nil { //创建记录
				return err
			}
			dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
		} else { // 存在
			if favoriteExist.State == 0 { //state为0-video的favorite_count加1
				dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + 1"))
				dao.SqlSession.Table("favorites").Where("video_id = ?", videoId).Update("state", 1)
			}
			//state为1-video的favorite_count不变
			return nil
		}
	} else { //2-取消点赞
		var favoriteCancel = &model.Favorite{}
		favoriteActionCancel := model.Favorite{
			UserId:  userId,
			VideoId: videoId,
			State:   0, //0-未点赞
		}
		if err := dao.SqlSession.Table("favorites").Where("user_id = ? AND video_id = ?", userId, videoId).First(&favoriteCancel).Error; err != nil { //找不到这条记录，取消点赞失败，创建记录
			dao.SqlSession.Table("favorites").Create(&favoriteActionCancel)
			return err
		}
		//存在
		if favoriteCancel.State == 1 { //state为1-video的favorite_count减1
			dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - 1"))
			//更新State
			dao.SqlSession.Table("favorites").Where("video_id = ?", videoId).Update("state", 0)
		}
		//state为0-video的favorite_count不变
		return nil
	}
	return nil
}

// ListFavorite 获取点赞列表
func ListFavorite(userId uint) ([]model.Video, error) {
	//查询当前id用户的所有点赞视频
	var favoriteList []model.Favorite
	videoList := make([]model.Video, 0)
	if err := dao.SqlSession.Table("favorites").Where("user_id=? AND state=?", userId, 1).Find(&favoriteList).Error; err != nil { //找不到记录
		return videoList, nil
	}
	for _, m := range favoriteList {

		var video = model.Video{}
		if err := dao.SqlSession.Table("videos").Where("id=?", m.VideoId).Find(&video).Error; err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}
