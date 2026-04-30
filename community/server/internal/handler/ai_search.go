package handler

import (
	"context"
	"strings"

	"community-server/internal/ai"
	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AISearchHandler struct {
	aiEngine      ai.Engine
	searchService *service.SearchService
}

func NewAISearchHandler(aiEngine ai.Engine, searchService *service.SearchService) *AISearchHandler {
	return &AISearchHandler{
		aiEngine:      aiEngine,
		searchService: searchService,
	}
}

func (h *AISearchHandler) AISearch(c *gin.Context) {
	var req model.AISearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	intent, err := h.aiEngine.ExtractSearchIntent(context.Background(), req.Question)
	if err != nil {
		intent = &ai.SearchIntent{
			Keywords: []string{req.Question},
			Intent:   req.Question,
		}
	}

	var allPosts []model.PostListItem
	var total int64

	for _, keyword := range intent.Keywords {
		searchReq := &model.SearchRequest{
			Keyword:  keyword,
			Page:     1,
			PageSize: 10,
		}
		result, err := h.searchService.SearchPosts(searchReq)
		if err == nil && result != nil {
			total += result.Total
			allPosts = append(allPosts, result.Posts...)
		}
	}

	seen := make(map[uint]bool)
	uniquePosts := make([]model.PostListItem, 0)
	for _, post := range allPosts {
		if !seen[post.ID] {
			seen[post.ID] = true
			uniquePosts = append(uniquePosts, post)
		}
	}

	if len(uniquePosts) > 10 {
		uniquePosts = uniquePosts[:10]
	}

	aiSummary := ""
	if len(uniquePosts) > 0 {
		postContext := ""
		for i, post := range uniquePosts {
			postContext += post.Title + ": " + post.Summary + "\n"
			if i >= 4 {
				break
			}
		}

		summaryPrompt := []ai.Message{
			{Role: "system", Content: "你是一个论坛助手，请根据搜索结果为用户的问题提供简要总结。"},
			{Role: "user", Content: "问题: " + req.Question + "\n相关帖子:\n" + postContext},
		}

		summary, err := h.aiEngine.Chat(context.Background(), summaryPrompt)
		if err == nil {
			aiSummary = summary
		}
	}

	response.Success(c, model.AISearchResponse{
		Keywords:  intent.Keywords,
		Intent:    intent.Intent,
		Posts:     uniquePosts,
		AISummary: aiSummary,
		Total:     total,
	})
}

func (h *AISearchHandler) AISearchStream(c *gin.Context) {
	var req model.AISearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	intent, err := h.aiEngine.ExtractSearchIntent(context.Background(), req.Question)
	if err != nil {
		intent = &ai.SearchIntent{
			Keywords: []string{req.Question},
			Intent:   req.Question,
		}
	}

	c.Writer.WriteString("event: intent\ndata: " + intent.Intent + "\n\n")
	c.Writer.Flush()

	var allPosts []model.PostListItem
	for _, keyword := range intent.Keywords {
		searchReq := &model.SearchRequest{
			Keyword:  keyword,
			Page:     1,
			PageSize: 10,
		}
		result, err := h.searchService.SearchPosts(searchReq)
		if err == nil && result != nil {
			allPosts = append(allPosts, result.Posts...)
		}
	}

	seen := make(map[uint]bool)
	uniquePosts := make([]model.PostListItem, 0)
	for _, post := range allPosts {
		if !seen[post.ID] {
			seen[post.ID] = true
			uniquePosts = append(uniquePosts, post)
		}
	}

	if len(uniquePosts) > 10 {
		uniquePosts = uniquePosts[:10]
	}

	postsJSON := ""
	for _, post := range uniquePosts {
		postsJSON += post.Title + "\n"
	}
	c.Writer.WriteString("event: posts\ndata: " + postsJSON + "\n\n")
	c.Writer.Flush()

	if len(uniquePosts) > 0 {
		postContext := ""
		for i, post := range uniquePosts {
			postContext += post.Title + ": " + post.Summary + "\n"
			if i >= 4 {
				break
			}
		}

		summaryPrompt := []ai.Message{
			{Role: "system", Content: "你是一个论坛助手，请根据搜索结果为用户的问题提供简要总结。"},
			{Role: "user", Content: "问题: " + req.Question + "\n相关帖子:\n" + postContext},
		}

		c.Writer.WriteString("event: summary_start\ndata: \n\n")
		c.Writer.Flush()

		h.aiEngine.StreamChat(context.Background(), summaryPrompt, func(chunk string, isFinish bool) {
			if isFinish {
				c.Writer.WriteString("event: summary_end\ndata: \n\n")
				c.Writer.Flush()
				return
			}
			chunk = strings.ReplaceAll(chunk, "\n", "\\n")
			c.Writer.WriteString("event: summary_chunk\ndata: " + chunk + "\n\n")
			c.Writer.Flush()
		})
	}
}
