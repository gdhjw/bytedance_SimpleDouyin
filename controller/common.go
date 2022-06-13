package controller

import (
	"net/http"
	"strconv"

	"bytedance_SimpleDouyin/dao"
	"bytedance_SimpleDouyin/model"
	"github.com/gin-gonic/gin"
)

func BindAndValid(ctx *gin.Context, target interface{}) bool {
	if err := ctx.ShouldBindQuery(target); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, Response{
			StatusCode: 1,
			StatusMsg:  "Wrong match!",
		})
		return false
	}
	return true
}

func QueryIDAndValid(ctx *gin.Context, queryName string) uint {
	id, err := strconv.ParseUint(ctx.Query(queryName), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  queryName + "is not number",
		})
		return 0
	}
	return uint(id)
}

func GetUserFromCtx(ctx *gin.Context) *models.User {
	user_, _ := ctx.Get("User")
	user, _ := user_.(*models.User)
	return user
}

func CommentsModelChange(commentList []models.Comment) []Comment {
	var comments []Comment
	for _, comment := range commentList {
		comments = append(comments, CommentModelChange(comment))
	}
	return comments
}

func CommentModelChange(comment models.Comment) Comment {
	user, _ := crud.GetUserByID(comment.UserID)
	return Comment{
		ID:         int64(comment.ID),
		Content:    comment.Content,
		CreateDate: comment.CreatedAt.Format("01-02"),
		User:       UserModelChange(*user),
	}
}


