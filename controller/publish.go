package controller

import (
	"bytedance_SimpleDouyin/dao"
	"bytedance_SimpleDouyin/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ReturnAuthor struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	//IsFollow      bool   `json:"is_follow"` // follow接口做好后加上
}

type ReturnVideo struct {
	VideoId       uint         `json:"video_id"`
	Author        ReturnAuthor `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount uint         `json:"favorite_count"`
	CommentCount  uint         `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

type VideoListResponse struct {
	StatusCode int32         `json:"status_code"`
	StatusMsg  string        `json:"status_msg,omitempty"`
	VideoList  []ReturnVideo `json:"video_list"`
}

// List 获取发布列表
func List(c *gin.Context) {
	//1.中间件鉴权token
	getHostId, _ := c.Get("user_id")
	var HostId uint
	if v, ok := getHostId.(uint); ok {
		HostId = v
	}
	// 获取用户id
	userId, _ := strconv.Atoi(c.Query("user_id"))
	GuestId := uint(userId)

	//根据用户id查找用户
	getUser, err := dao.GetUser(GuestId)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": 1,
			"status_msg":  "Not find this person",
		})
		c.Abort()
		return
	}

	returnAuthor := ReturnAuthor{
		AuthorId:      GuestId,
		Name:          getUser.Name,
		FollowCount:   getUser.FollowCount,
		FollowerCount: getUser.FollowerCount,
		//IsFollow:      dao.IsFollowing(HostId, GuestId), // 是否关注
	}
	var videoList []model.Video
	var num int
	num, videoList = dao.GetVideoList(GuestId)
	// 无视频
	if num == 0 {
		c.JSON(200, gin.H{
			"status_code": 1,
			"status_msg":  "null",
			"video_list":  videoList,
		})
	} else { // 返回视频列表
		var returnVideoList []ReturnVideo
		for i := 0; i < num; i++ {
			returnVideo := ReturnVideo{
				VideoId:       videoList[i].ID,
				Author:        returnAuthor,
				PlayUrl:       videoList[i].PlayUrl,
				CoverUrl:      videoList[i].CoverUrl,
				FavoriteCount: videoList[i].FavoriteCount,
				CommentCount:  videoList[i].CommentCount,
				IsFavorite:    dao.CheckFavorite(HostId, videoList[i].ID),
				Title:         videoList[i].Title,
			}
			returnVideoList = append(returnVideoList, returnVideo)
		}
		c.JSON(200, gin.H{
			"status_code": 0,
			"status_msg":  "success",
			"video_list":  returnVideoList,
		})
	}
}
