package handler

import (
	"context"
	"io"

	"community-server/internal/ai"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// BotWebhookHandler 处理 JuggleIM Bot Engine 的回调
type BotWebhookHandler struct {
	aiEngine ai.Engine
}

func NewBotWebhookHandler(aiEngine ai.Engine) *BotWebhookHandler {
	return &BotWebhookHandler{
		aiEngine: aiEngine,
	}
}

// Webhook JuggleIM Bot DefaultEngine 回调
func (h *BotWebhookHandler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		zap.S().Warn("Bot webhook 读取请求体失败", "error", err)
		response.ErrorWithMsg(c, response.CodeInvalidParam, "读取请求失败")
		return
	}

	reply, err := h.aiEngine.Chat(context.Background(), []ai.Message{
		{Role: "user", Content: string(body)},
	})
	if err != nil {
		zap.S().Warn("Bot AI 对话失败", "error", err)
		response.ErrorWithMsg(c, response.CodeServerBusy, "AI 响应失败")
		return
	}

	c.String(200, reply)
}
