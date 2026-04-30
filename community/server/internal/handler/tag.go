package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService *service.TagService
}

func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

func (h *TagHandler) CreateTag(c *gin.Context) {
	var req model.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	tagID, err := h.tagService.CreateTag(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, gin.H{"tag_id": tagID})
}

func (h *TagHandler) GetTagList(c *gin.Context) {
	var req model.TagListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.tagService.GetTagList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}

func (h *TagHandler) UpdateTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的标签ID")
		return
	}

	var req model.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	err = h.tagService.UpdateTag(uint(tagID), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TagHandler) DeleteTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的标签ID")
		return
	}

	err = h.tagService.DeleteTag(uint(tagID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TagHandler) AddPostTags(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	var req model.AddPostTagsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	err = h.tagService.AddPostTags(uint(postID), req.TagIDs)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TagHandler) RemovePostTag(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	tagIDStr := c.Param("tag_id")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的标签ID")
		return
	}

	err = h.tagService.RemovePostTag(uint(postID), uint(tagID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *TagHandler) GetPostTags(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	tags, err := h.tagService.GetPostTags(uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, tags)
}

func (h *TagHandler) GetPostsByTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := strconv.ParseUint(tagIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的标签ID")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, err := h.tagService.GetPostsByTag(uint(tagID), page, pageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}
