package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService CommentService
}

func NewCommentHandler(commentService CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

// CreateComment 发表评论
// @Summary 发表评论
// @Tags 评论
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.CreateCommentRequest true "评论内容"
// @Success 200 {object} response.Response
// @Router /comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
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

// GetCommentList 评论列表
// @Summary 评论列表
// @Tags 评论
// @Produce json
// @Param post_id query uint true "帖子ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /comments [get]
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

// UpdateComment 编辑评论
// @Summary 编辑评论
// @Tags 评论
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "评论ID"
// @Param body body model.UpdateCommentRequest true "更新内容"
// @Success 200 {object} response.Response
// @Router /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
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

// DeleteComment 删除评论
// @Summary 删除评论
// @Tags 评论
// @Security BearerAuth
// @Produce json
// @Param id path uint true "评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
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

// LikeComment 点赞评论
// @Summary 点赞评论
// @Tags 评论
// @Security BearerAuth
// @Produce json
// @Param id path uint true "评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id}/like [post]
func (h *CommentHandler) LikeComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
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

// UnlikeComment 取消点赞评论
// @Summary 取消点赞评论
// @Tags 评论
// @Security BearerAuth
// @Produce json
// @Param id path uint true "评论ID"
// @Success 200 {object} response.Response
// @Router /comments/{id}/like [delete]
func (h *CommentHandler) UnlikeComment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
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
