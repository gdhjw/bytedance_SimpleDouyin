package router

import (
	"bytedance_SimpleDouyin/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 未加入jwt中间件!!!

	// 主路由
	douyin := r.Group("/douyin")

	publish := douyin.Group("/publish")
	{
		// 发布列表路由
		publish.GET("/list/", controller.List)
		publish.Post("/action/", controller.Publish)
	}
	comment := douyin.Group("/comment")
	{
		publish.Post("/action/", controller.CommentAction)
		publish.get("/list/", controller.commentList)
	}
	favorite := douyin.Group("/favorite")
	{
		// 点赞路由
		favorite.POST("/action", controller.ActionFavorite)
		// 获取点赞列表
		favorite.POST("/list", controller.ListFavorite)
	}

	return r

}
