package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService *service.CommentService
}

func NewCommentHandler(commentService *service.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	commentID, err := h.commentService.CreateComment(userID.(uint), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, gin.H{"comment_id": commentID})
}

func (h *CommentHandler) GetCommentList(c *gin.Context) {
	var req model.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.commentService.GetCommentList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的评论ID")
		return
	}

	var req model.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	err = h.commentService.UpdateComment(userID.(uint), uint(commentID), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的评论ID")
		return
	}

	err = h.commentService.DeleteComment(userID.(uint), uint(commentID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *CommentHandler) LikeComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的评论ID")
		return
	}

	err = h.commentService.LikeComment(userID.(uint), uint(commentID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *CommentHandler) UnlikeComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}

	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseUint(commentIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的评论ID")
		return
	}

	err = h.commentService.UnlikeComment(userID.(uint), uint(commentID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}
