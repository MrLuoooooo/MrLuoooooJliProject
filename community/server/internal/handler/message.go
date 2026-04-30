package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("user_id")

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

func (h *MessageHandler) GetMessageList(c *gin.Context) {
	userID, _ := c.Get("user_id")

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

func (h *MessageHandler) GetConversationList(c *gin.Context) {
	userID, _ := c.Get("user_id")

	result, err := h.messageService.GetConversationList(userID.(uint))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}

func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	count, err := h.messageService.GetUnreadCount(userID.(uint))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, model.UnreadCountResponse{Count: count})
}

func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
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
