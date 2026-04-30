package handler

import (
	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

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
