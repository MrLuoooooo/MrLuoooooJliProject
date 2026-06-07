package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService TagService
}

func NewTagHandler(tagService TagService) *TagHandler {
	return &TagHandler{
		tagService: tagService,
	}
}

// CreateTag 创建标签
// @Summary 创建标签
// @Tags 标签
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.CreateTagRequest true "标签信息"
// @Success 200 {object} response.Response
// @Router /tags [post]
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

// GetTagList 标签列表
// @Summary 标签列表
// @Tags 标签
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /tags [get]
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

// UpdateTag 编辑标签
// @Summary 编辑标签
// @Tags 标签
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "标签ID"
// @Param body body model.UpdateTagRequest true "更新内容"
// @Success 200 {object} response.Response
// @Router /tags/{id} [put]
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

// DeleteTag 删除标签
// @Summary 删除标签
// @Tags 标签
// @Security BearerAuth
// @Produce json
// @Param id path uint true "标签ID"
// @Success 200 {object} response.Response
// @Router /tags/{id} [delete]
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

// AddPostTags 给帖子打标签
// @Summary 帖子打标签
// @Tags 标签
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "帖子ID"
// @Param body body model.AddPostTagsRequest true "标签ID列表"
// @Success 200 {object} response.Response
// @Router /posts/{id}/tags [post]
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

// RemovePostTag 移除帖子标签
// @Summary 移除标签
// @Tags 标签
// @Security BearerAuth
// @Produce json
// @Param id path uint true "帖子ID"
// @Param tag_id path uint true "标签ID"
// @Success 200 {object} response.Response
// @Router /posts/{id}/tags/{tag_id} [delete]
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

// GetPostTags 帖子的标签
// @Summary 获取帖子标签
// @Tags 标签
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /posts/{id}/tags [get]
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

// GetPostsByTag 标签下的帖子
// @Summary 标签下的帖子
// @Tags 标签
// @Produce json
// @Param id path uint true "标签ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /tags/{id}/posts [get]
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
