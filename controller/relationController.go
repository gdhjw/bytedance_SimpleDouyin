package controller

import (
	"bytedance_SimpleDouyin/common"
	"bytedance_SimpleDouyin/middleware"
	"bytedance_SimpleDouyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ReturnFollower 关注表与粉丝表共用的用户数据模型
type ReturnFollower struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// FollowingListResponse 关注表相应结构体
type FollowingListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
}

// FollowerListResponse 粉丝表相应结构体
type FollowerListResponse struct {
	common.Response
	UserList []ReturnFollower `json:"user_list"`
}

// RelationAction 关注/取消关注操作
func RelationAction(c *gin.Context) {
	//1.取数据
	//1.1 从token中获取用户id
	strToken := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(strToken)
	hostId := tokenStruct.UserId
	//1.2 获取待关注的用户id
	getToUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	guestId := uint(getToUserId)
	//1.3 获取关注操作（关注1，取消关注2）
	getActionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	actionType := uint(getActionType)

	//2.自己关注/取消关注自己不合法
	if hostId == guestId {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 405,
			StatusMsg:  "无法关注自己",
		})
		c.Abort()
		return
	}

	//3.service层进行关注/取消关注处理
	err := service.FollowAction(hostId, guestId, actionType)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  "关注/取消关注成功！",
		})
	}
}
