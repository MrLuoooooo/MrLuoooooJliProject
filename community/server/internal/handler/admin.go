package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
	}
}

func (h *AdminHandler) GetUserList(c *gin.Context) {
	var req model.AdminUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	result, err := h.adminService.GetUserList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	if err := h.adminService.DeleteUser(uint(userID)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *AdminHandler) UpdateUserAdminType(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	var req model.UpdateUserAdminTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.adminService.UpdateUserAdminType(uint(userID), req.AdminType); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	userID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	var req model.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.adminService.UpdateUserStatus(uint(userID), req.Status); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *AdminHandler) GetPostList(c *gin.Context) {
	var req model.AdminPostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	result, err := h.adminService.GetPostList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *AdminHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	if err := h.adminService.DeletePost(uint(postID)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}
