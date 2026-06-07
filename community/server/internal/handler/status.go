package handler

import (
	"strconv"

	"community-server/internal/im"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	imClient im.IMClient
}

func NewStatusHandler(imClient im.IMClient) *StatusHandler {
	return &StatusHandler{
		imClient: imClient,
	}
}

// GetOnlineStatus 查询用户在线状态
// @Summary 在线状态
// @Tags 用户
// @Produce json
// @Param user_id path uint true "用户ID"
// @Success 200 {object} response.Response
// @Router /users/{user_id}/online [get]
func (h *StatusHandler) GetOnlineStatus(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	statuses, err := h.imClient.QueryOnlineStatus([]string{im.UserIDToStr(uint(userID))})
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	isOnline := false
	if s, ok := statuses[im.UserIDToStr(uint(userID))]; ok {
		isOnline = s
	}

	response.Success(c, gin.H{"user_id": userID, "is_online": isOnline})
}

// BatchOnlineStatus 批量查询在线状态
// @Summary 批量在线状态
// @Tags 用户
// @Accept json
// @Produce json
// @Param body body object true "user_ids: [1,2,3]"
// @Success 200 {object} response.Response
// @Router /users/online/batch [post]
func (h *StatusHandler) BatchOnlineStatus(c *gin.Context) {
	var req struct {
		UserIDs []uint `json:"user_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	strIDs := make([]string, 0, len(req.UserIDs))
	for _, id := range req.UserIDs {
		strIDs = append(strIDs, im.UserIDToStr(id))
	}

	statuses, err := h.imClient.QueryOnlineStatus(strIDs)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, statuses)
}
