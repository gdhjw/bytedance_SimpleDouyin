package controller

import (
	"bytedance_SimpleDouyin/common"
	"bytedance_SimpleDouyin/model"
	"bytedance_SimpleDouyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"path/filepath"
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
type PublishActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type VideoListResponse struct {
	PublishActionResponse
	VideoList []ReturnVideo `json:"video_list"`
}

func Publish(c *gin.Context) {

	//1.中间件验证token后，获取userId
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	//2.接收请求参数信息
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	//3.返回至前端页面的展示信息
	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, filename)
	saveFile := filepath.Join("../resources/videos/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		StatusCode: 0,
		StatusMsg:  finalName + "--uploaded successfully",
	})
	//尝试用远程服务器但未部署成功，此时用本地静态资源服务器
	var playUrl string
	playUrl = "http://172.22.108.166:8000/" + "videos/" + finalName
	//封面url已写死
	var coverUrl string
	coverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"

	//4.保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
	video := model.Video{
		Model:         gorm.Model{},
		AuthorId:      userId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	service.CreateVideo(&video)
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
	getUser, err := service.GetUser(GuestId)
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
		//IsFollow:      service.IsFollowing(HostId, GuestId), // 是否关注
	}
	var videoList []model.Video
	var num int
	num, videoList = service.GetVideoList(GuestId)
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
				IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
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
