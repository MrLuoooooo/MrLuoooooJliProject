package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"community-server/internal/ai"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AIHandler struct {
	engine ai.Engine
}

func NewAIHandler(engine ai.Engine) *AIHandler {
	return &AIHandler{
		engine: engine,
	}
}

func (h *AIHandler) Chat(c *gin.Context) {
	var req ai.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if len(req.Messages) == 0 {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "消息列表不能为空")
		return
	}

	if req.Stream {
		h.streamChat(c, req.Messages)
	} else {
		h.normalChat(c, req.Messages)
	}
}

func (h *AIHandler) normalChat(c *gin.Context, messages []ai.Message) {
	content, err := h.engine.Chat(c.Request.Context(), messages)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, fmt.Sprintf("AI服务调用失败: %v", err))
		return
	}

	response.Success(c, ai.ChatResponse{
		Content: content,
	})
}

func (h *AIHandler) streamChat(c *gin.Context, messages []ai.Message) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		http.Error(c.Writer, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	err := h.engine.StreamChat(c.Request.Context(), messages, func(chunk string, isFinish bool) {
		if isFinish {
			fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
			flusher.Flush()
			return
		}

		data := fmt.Sprintf("data: %s\n\n", chunk)
		fmt.Fprint(c.Writer, data)
		flusher.Flush()
	})

	if err != nil {
		errorMsg := fmt.Sprintf("data: [ERROR] %s\n\n", err.Error())
		fmt.Fprint(c.Writer, errorMsg)
		flusher.Flush()
	}
}

func (h *AIHandler) ChatSSE(c *gin.Context) {
	var req ai.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	if len(req.Messages) == 0 {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "消息列表不能为空")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")
	c.Header("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		http.Error(c.Writer, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	err := h.engine.StreamChat(c.Request.Context(), req.Messages, func(chunk string, isFinish bool) {
		if isFinish {
			eventData := SSEvent{
				Event: "done",
				Data:  "",
			}
			fmt.Fprintf(c.Writer, "event: done\ndata: %s\n\n", eventData.Data)
			flusher.Flush()
			return
		}

		eventData := SSEvent{
			Event: "message",
			Data:  chunk,
		}
		fmt.Fprintf(c.Writer, "event: message\ndata: %s\n\n", eventData.Data)
		flusher.Flush()
	})

	if err != nil {
		eventData := SSEvent{
			Event: "error",
			Data:  err.Error(),
		}
		fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", eventData.Data)
		flusher.Flush()
	}
}

type SSEvent struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

func FormatSSE(data string) string {
	lines := strings.Split(data, "\n")
	result := ""
	for _, line := range lines {
		result += fmt.Sprintf("data: %s\n", line)
	}
	result += "\n"
	return result
}

func SendSSEvent(w io.Writer, event, data string) {
	fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
}

func SendSSData(w io.Writer, data string) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		fmt.Fprintf(w, "data: %s\n", line)
	}
	fmt.Fprint(w, "\n")
}
