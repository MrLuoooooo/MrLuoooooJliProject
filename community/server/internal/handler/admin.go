package handler

import (
	"strconv"

	"community-server/internal/im"
	"community-server/internal/model"
	"community-server/internal/ws"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService AdminService
	imClient     im.IMClient
	wsManager    *ws.Manager
}

func NewAdminHandler(adminService AdminService, imClient im.IMClient, wsManager *ws.Manager) *AdminHandler {
	return &AdminHandler{
		adminService: adminService,
		imClient:     imClient,
		wsManager:    wsManager,
	}
}

// GetUserList 用户列表
// @Summary 用户列表
// @Tags 管理
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /admin/users [get]
func (h *AdminHandler) GetUserList(c *gin.Context) {
	var req model.AdminUserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req = model.AdminUserListRequest{Page: 1, PageSize: 20}
	}
	result, err := h.adminService.GetUserList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Tags 管理
// @Security BearerAuth
// @Produce json
// @Param id path uint true "用户ID"
// @Success 200 {object} response.Response
// @Router /admin/users/{id} [delete]
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	if err := h.adminService.DeleteUser(uint(id)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateUserAdminType 设置管理员
// @Summary 设置管理员
// @Tags 管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Param body body model.UpdateUserAdminTypeRequest true "管理员类型"
// @Success 200 {object} response.Response
// @Router /admin/users/{id}/admin_type [put]
func (h *AdminHandler) UpdateUserAdminType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	var req model.UpdateUserAdminTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.adminService.UpdateUserAdminType(uint(id), req.AdminType); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// UpdateUserStatus 封禁/解封用户
// @Summary 封禁/解封
// @Tags 管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "用户ID"
// @Param body body model.UpdateUserStatusRequest true "状态"
// @Success 200 {object} response.Response
// @Router /admin/users/{id}/status [put]
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	var req model.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.adminService.UpdateUserStatus(uint(id), req.Status); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetPostList 帖子管理列表
// @Summary 帖子管理
// @Tags 管理
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /admin/posts [get]
func (h *AdminHandler) GetPostList(c *gin.Context) {
	var req model.AdminPostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req = model.AdminPostListRequest{Page: 1, PageSize: 20}
	}
	result, err := h.adminService.GetPostList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

// DeletePost 管理员删帖
// @Summary 管理员删帖
// @Tags 管理
// @Security BearerAuth
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /admin/posts/{id} [delete]
func (h *AdminHandler) DeletePost(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.ParseUint(idStr, 10, 64)
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

// BroadcastRequest 全站广播请求
type BroadcastRequest struct {
	Content string `json:"content" binding:"required,max=2000"`
}

// SendBroadcast 发送全站广播（管理员专用）
// @Summary 全站广播
// @Tags 管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body BroadcastRequest true "广播内容"
// @Success 200 {object} response.Response
// @Router /admin/broadcast [post]
func (h *AdminHandler) SendBroadcast(c *gin.Context) {
	var req BroadcastRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
		return
	}
	if err := h.imClient.SendBroadcastMsg(im.UserIDToStr(userID.(uint)), req.Content); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	h.wsManager.Broadcast(ws.PushMessage{
		Type: "broadcast",
		Data: map[string]interface{}{"content": req.Content},
	})
	response.Success(c, nil)
}

// SetPostTop 设置/取消帖子置顶
// @Summary 置顶/取消置顶
// @Tags 管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "帖子ID"
// @Param body body object true "is_top: true/false"
// @Success 200 {object} response.Response
// @Router /admin/posts/{id}/top [put]
func (h *AdminHandler) SetPostTop(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	var req struct {
		IsTop bool `json:"is_top"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.adminService.SetPostTop(uint(postID), req.IsTop); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}

// SetPostEssence 设置/取消帖子精华
// @Summary 精华/取消精华
// @Tags 管理
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "帖子ID"
// @Param body body object true "is_essence: true/false"
// @Success 200 {object} response.Response
// @Router /admin/posts/{id}/essence [put]
func (h *AdminHandler) SetPostEssence(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	var req struct {
		IsEssence bool `json:"is_essence"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.adminService.SetPostEssence(uint(postID), req.IsEssence); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetStats 获取统计数据
// @Summary 仪表盘统计
// @Tags 管理
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /admin/stats [get]
func (h *AdminHandler) GetStats(c *gin.Context) {
	result, err := h.adminService.GetStats()
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, result)
}
