package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService FollowService
}

func NewFollowHandler(followService FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

// FollowUser 关注用户
// @Summary 关注用户
// @Tags 关注
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.FollowRequest true "关注目标"
// @Success 200 {object} response.Response
// @Router /follows [post]
func (h *FollowHandler) FollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}
	var req model.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.followService.FollowUser(userID.(uint), &req); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnfollowUser 取消关注
// @Summary 取消关注
// @Tags 关注
// @Security BearerAuth
// @Produce json
// @Param id path uint true "被关注用户ID"
// @Success 200 {object} response.Response
// @Router /follows/{id} [delete]
func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}
	followIDStr := c.Param("id")

	followID, err := strconv.ParseUint(followIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	if err := h.followService.UnfollowUser(userID.(uint), uint(followID)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetFollowers 粉丝列表
// @Summary 粉丝列表
// @Tags 关注
// @Produce json
// @Param user_id query uint true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /follows/followers [get]
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	var req model.FollowListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	result, err := h.followService.GetFollowers(req.UserID, req.Page, req.PageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

// GetFollowing 关注列表
// @Summary 关注列表
// @Tags 关注
// @Produce json
// @Param user_id query uint true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /follows/following [get]
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	var req model.FollowListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	result, err := h.followService.GetFollowing(req.UserID, req.Page, req.PageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

// IsFollowing 是否已关注
// @Summary 关注状态
// @Tags 关注
// @Security BearerAuth
// @Produce json
// @Param id path uint true "目标用户ID"
// @Success 200 {object} response.Response
// @Router /follows/{id}/status [get]
func (h *FollowHandler) IsFollowing(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}
	followIDStr := c.Param("id")

	followID, err := strconv.ParseUint(followIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	isFollowing, err := h.followService.IsFollowing(userID.(uint), uint(followID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, "查询失败")
		return
	}
	response.Success(c, gin.H{"is_following": isFollowing})
}

// GetFollowCounts 关注数
// @Summary 关注统计
// @Tags 关注
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} response.Response
// @Router /follows/{id}/counts [get]
func (h *FollowHandler) GetFollowCounts(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	followers, following := h.followService.GetFollowCounts(uint(userID))
	response.Success(c, gin.H{
		"followers": followers,
		"following": following,
	})
}
