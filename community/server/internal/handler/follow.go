package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService *service.FollowService
}

func NewFollowHandler(followService *service.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

func (h *FollowHandler) FollowUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

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

func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	userID, _ := c.Get("user_id")
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

func (h *FollowHandler) GetFollowers(c *gin.Context) {
	var req model.FollowListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	result, err := h.followService.GetFollowers(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *FollowHandler) GetFollowing(c *gin.Context) {
	var req model.FollowListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	result, err := h.followService.GetFollowing(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *FollowHandler) IsFollowing(c *gin.Context) {
	userID, _ := c.Get("user_id")
	followIDStr := c.Param("id")

	followID, err := strconv.ParseUint(followIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	isFollowing, _ := h.followService.IsFollowing(userID.(uint), uint(followID))
	response.Success(c, gin.H{"is_following": isFollowing})
}

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
