package handler

import (
	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService UserService
}

func NewUserHandler(userService UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 注册
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body model.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response
// @Router /users/register [post]
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

// Login 用户登录
// @Summary 登录
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "登录信息"
// @Success 200 {object} response.Response
// @Router /users/login [post]
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

// GetProfile 获取个人资料
// @Summary 个人资料
// @Tags 用户
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /users/profile [get]
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

	response.Success(c, user)
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

// ForgotPassword 忘记密码
// @Summary 获取重置令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body model.ForgotPasswordRequest true "邮箱"
// @Success 200 {object} response.Response
// @Router /users/forgot-password [post]
func (h *UserHandler) ForgotPassword(c *gin.Context) {
	var req model.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "请输入邮箱地址")
		return
	}
	if err := h.userService.ForgotPassword(req.Email); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body model.ResetPasswordRequest true "令牌+新密码"
// @Success 200 {object} response.Response
// @Router /users/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}
	if err := h.userService.ResetPassword(req.Token, req.NewPassword); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}
