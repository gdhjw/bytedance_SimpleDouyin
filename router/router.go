package router

import (
	"bytedance_SimpleDouyin/controller"
	"bytedance_SimpleDouyin/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 主路由
	douyinGroup := r.Group("/douyin")

	userGroup := douyinGroup.Group("/user")
	{

		userGroup.GET("/", controller.UserInfo)
		userGroup.POST("/login/", controller.UserLogin)
		userGroup.POST("/register/", controller.UserRegister)
	}

	publish := userGroup.Group("/publish")
	{
		// 发布列表路由
		publish.GET("/list/", controller.List)
		publish.POST("/action/", controller.Publish)
	}
	comment := userGroup.Group("/comment")
	{
		comment.POST("/action/", controller.CommentAction)
		comment.GET("/list/", controller.CommentList)
	}

	favorite := r.Group("/favorite")
	{
		// 点赞路由
		favorite.POST("/action", controller.ActionFavorite)
		// 获取点赞列表
		favorite.POST("/list", controller.ListFavorite)
	}

	relationGroup := douyinGroup.Group("relation")
	{
		relationGroup.POST("/action/", middleware.JwtMiddleware(), controller.RelationAction)

	}

	return r

}
