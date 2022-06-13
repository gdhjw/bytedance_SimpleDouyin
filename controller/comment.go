package controller

import (
	"bytedance_SimpleDouyin/dao"
	"bytedance_SimpleDouyin/model"
	"github.com/jinzhu/gorm"
)

type CommentResponse struct {
	StatusCode int32   `json:"status_code,omitempty"`
	StatusMsg  string  `json:"status_msg,omitempty"`
	Comment    Comment `json:"comment,omitempty"`
}
type CommentListResponse struct {
	StatusCode  int32     `json:"status_code,omitempty"`
	StatusMsg   string    `json:"status_msg,omitempty"`
	CommentList []Comment `json:"comment_list,omitempty"`
}

func CommentAction(ctx *gin.Context) {
	var request struct {
		VideoID     uint   `form:"video_id" binding:"required"`
		ActionType  uint   `form:"action_type" binding:"required,min=1,max=2"`
		CommentText string `form:"comment_text" binding:"omitempty"`
		CommentID   uint   `form:"comment_id" binding:"omitempty"`
	}

	if !BindAndValid(ctx, &request) {
		return
	}
	user := GetUserFromCtx(ctx)
	if request.ActionType == 1 {
		comment := dao.CreateComment(&models.Comment{
			UserID:  user.ID,
			VideoID: request.VideoID,
			Content: request.CommentText,
		})
		ctx.JSON(http.StatusOK, CommentResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Comment success!",
			},
			Comment: CommentModelChange(*comment),
		})
	} else {
		dao.DeleteComment(request.CommentID)
		ctx.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "Comment delete!",
		})
	}

}

func CommentList(ctx *gin.Context) {
	videoID := QueryIDAndValid(ctx, "video_id")
	if videoID == 0 {
		return
	}

	ctx.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: CommentsModelChange(dao.GetComments(videoID)),
	})
}
