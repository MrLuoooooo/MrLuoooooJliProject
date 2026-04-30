package handler

import (
	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	userID, err := h.userService.Register(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, err.Error())
		return
	}

	response.Success(c, gin.H{"user_id": userID})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	resp, err := h.userService.Login(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeUnauthorized, err.Error())
		return
	}

	response.Success(c, resp)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}

	user, err := h.userService.GetUserByID(userID.(uint))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeNotFound, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"avatar":   user.Avatar,
		"bio":      user.Bio,
		"email":    user.Email,
		"status":   user.Status,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}

	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	err := h.userService.UpdateProfile(userID.(uint), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}
