package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService NotificationService
}

func NewNotificationHandler(notificationService NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GetList 通知列表
// @Summary 通知列表
// @Tags 通知
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /notifications [get]
func (h *NotificationHandler) GetList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	var req model.NotificationListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.notificationService.GetList(userID.(uint), req.Page, req.PageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, list)
}

// MarkRead 标记已读
// @Summary 标记通知已读
// @Tags 通知
// @Security BearerAuth
// @Produce json
// @Param id path uint true "通知ID"
// @Success 200 {object} response.Response
// @Router /notifications/{id}/read [put]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	notifIDStr := c.Param("id")
	notifID, err := strconv.ParseUint(notifIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的通知ID")
		return
	}

	if err := h.notificationService.MarkRead(uint(notifID), userID.(uint)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}

// MarkAllRead 全部已读
// @Summary 全部标记已读
// @Tags 通知
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /notifications/read-all [put]
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	if err := h.notificationService.MarkAllRead(userID.(uint)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}

// GetUnreadCount 未读计数
// @Summary 未读通知数
// @Tags 通知
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /notifications/unread [get]
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	count, err := h.notificationService.GetUnreadCount(userID.(uint))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, model.UnreadCountResponse{Count: count})
}
