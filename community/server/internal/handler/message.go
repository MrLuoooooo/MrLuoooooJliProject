package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService MessageService
}

func NewMessageHandler(messageService MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// SendMessage 发送私信
// @Summary 发送私信
// @Tags 私信
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.SendMessageRequest true "消息内容"
// @Success 200 {object} response.Response
// @Router /messages [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	var req model.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if err := h.messageService.SendMessage(userID.(uint), &req); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetMessageList 消息列表
// @Summary 私信列表
// @Tags 私信
// @Security BearerAuth
// @Produce json
// @Param sender_id query uint false "对方用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /messages [get]
func (h *MessageHandler) GetMessageList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	var req model.MessageListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if req.ReceiverID == 0 {
		req.ReceiverID = userID.(uint)
	}

	result, err := h.messageService.GetMessageList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

// GetConversationList 会话列表
// @Summary 会话列表
// @Tags 私信
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /messages/conversations [get]
func (h *MessageHandler) GetConversationList(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	result, err := h.messageService.GetConversationList(userID.(uint))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

// GetUnreadCount 未读私信数
// @Summary 未读私信数
// @Tags 私信
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Router /messages/unread [get]
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	count, err := h.messageService.GetUnreadCount(userID.(uint))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, model.UnreadCountResponse{Count: count})
}

// MarkAsRead 标记已读
// @Summary 标记私信已读
// @Tags 私信
// @Security BearerAuth
// @Produce json
// @Param id path uint true "消息ID"
// @Success 200 {object} response.Response
// @Router /messages/{id}/read [put]
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}
	messageIDStr := c.Param("id")

	messageID, err := strconv.ParseUint(messageIDStr, 10, 32)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的消息ID")
		return
	}

	if err := h.messageService.MarkAsRead(uint(messageID), userID.(uint)); err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}
