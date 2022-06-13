package controller

import (
	"bytedance_SimpleDouyin/common"
	"bytedance_SimpleDouyin/middleware"
	"bytedance_SimpleDouyin/model"
	"bytedance_SimpleDouyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserIdTokenResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserRegisterResponse struct {
	common.Response
	UserIdTokenResponse
}

type UserLoginResponse struct {
	common.Response
	UserIdTokenResponse
}

//用户登录

func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	userLoginResponse, err := UserLoginService(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		Response:            common.Response{StatusCode: 0},
		UserIdTokenResponse: userLoginResponse,
	})

}
func UserLoginService(username string, password string) (UserIdTokenResponse, error) {

	//数据准备
	var userResponse = UserIdTokenResponse{}

	var login model.Users
	tempPassword := service.FindPasswordByName(username)
	if !service.ComparePasswords(tempPassword, password) {
		return userResponse, common.ErrorPasswordFalse
	}

	//颁发token
	token, err := middleware.CreateToken(login.Id, login.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = UserIdTokenResponse{
		UserId: login.Id,
		Token:  token,
	}
	return userResponse, nil

}

//用户注册

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	RegisterResponse, err := UserRegisterService(username, password)
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return

	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response:            common.Response{StatusCode: 0},
		UserIdTokenResponse: RegisterResponse,
	})
	return

}

func UserRegisterService(UserName string, Password string) (UserIdTokenResponse, error) {
	//数据准备
	var response = UserIdTokenResponse{}
	//检验合法性
	err := service.ISUserLegal(UserName, Password)
	if err != nil {
		return response, err
	}
	//新建用户
	service.CreateRegisterUser2(UserName, Password)

	id := service.FindIdByName(UserName)

	//颁发token
	token, err := middleware.CreateToken(id, UserName)
	if err != nil {
		return response, err
	}
	response = UserIdTokenResponse{
		UserId: id,
		Token:  token,
	}
	return response, nil
}

type UserInfoResponse struct {
	common.Response
	UserList UserInfoQueryResponse `json:"user"`
}
type UserInfoQueryResponse struct {
	UserId        uint   `json:"user_id"`
	UserName      string `json:"name"`
	FollowCount   uint   `json:"follow_count"`
	FollowerCount uint   `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// UserInfo 用户信息主函数
func UserInfo(c *gin.Context) {
	//根据user_id查询
	rawId := c.Query("user_id")
	userInfoResponse, err := UserInfoService(rawId)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: common.Response{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		UserList: userInfoResponse,
	})

}

// UserInfoService 用户信息处理函数
func UserInfoService(rawId string) (UserInfoQueryResponse, error) {
	//0.数据准备
	var userInfoQueryResponse = UserInfoQueryResponse{}
	userId, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		return userInfoQueryResponse, err
	}

	//1.获取用户信息
	var user model.Users
	user, err = service.GetUserById2(uint(userId))
	if err != nil {
		return userInfoQueryResponse, err
	}

	userInfoQueryResponse = UserInfoQueryResponse{
		UserId:        user.Id,
		UserName:      user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      false,
	}
	return userInfoQueryResponse, nil
}
