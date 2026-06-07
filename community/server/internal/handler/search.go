package handler

import (
	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService SearchService
}

func NewSearchHandler(searchService SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// SearchPosts 搜索
// @Summary 搜索
// @Tags 搜索
// @Produce json
// @Param keyword query string true "关键词"
// @Param type query string false "类型: post, user" default(post)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /search [get]
func (h *SearchHandler) SearchPosts(c *gin.Context) {
	var req model.SearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	var result *model.SearchResponse
	var err error

	if req.Type == "user" {
		result, err = h.searchService.SearchUsers(&req)
	} else {
		result, err = h.searchService.SearchPosts(&req)
	}

	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, result)
}
