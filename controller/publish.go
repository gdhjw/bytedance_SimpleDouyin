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
type PublishActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}
type VideoListResponse struct {
	PublishActionResponse
	VideoList []ReturnVideo `json:"video_list"`
}

func Publish(c *gin.Context) {

	//校验token
	// token := ctx.PostForm("token")
	// userID, err := middlewares.Parse(ctx, token)
	// if err != nil {
	// 	return
	// }
	// 判断 token 中的 user 是否存在
	username, _ := c.Get("username")
	var user db.User
	var video db.Video

	result := db.Mysql.Where(" name = ?", username).Find(&user)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, PublishActionResponse{
			StatusCode: 1,
			StatusMsg:  "User does not exist!",
		})
		return
	}

	// 读取视频数据 data
	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusBadRequest, PublishActionResponse{
			StatusCode: 2,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 加工保存视频数据
	filename := filepath.Base(data.Filename)

	finalName := fmt.Sprintf("%d_%s", user.ID, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusBadRequest, PublishActionResponse{
			StatusCode: 3,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 读取视频标题title
	title := c.PostForm("title")
	if title == "" {
		ctx.JSON(http.StatusBadRequest, PublishActionResponse{
			StatusCode: 4,
			StatusMsg:  "Empty title!",
		})
	}

	video.UserID = user.ID
	// video.PlayUrl = "localhost:8080/static/" + finalName
	// video.CoverUrl = "localhost:8080/static/20190330000110_360x480_55.jpg"

	// fmt.Print("title = ", video.Title)
	db.Mysql.Save(&video)

	c.JSON(http.StatusOK, PublishActionResponse{
		StatusCode: 0,
		StatusMsg:  finalName + "is uploaded successfully",
	})
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
