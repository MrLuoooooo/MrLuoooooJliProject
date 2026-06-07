package handler

import (
	"strconv"

	"community-server/internal/model"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService PostService
}

func NewPostHandler(postService PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

// CreatePost 发布帖子
// @Summary 发布帖子
// @Tags 帖子
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body model.CreatePostRequest true "帖子内容"
// @Success 200 {object} response.Response
// @Router /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	var req model.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	postID, err := h.postService.CreatePost(userID.(uint), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, gin.H{"post_id": postID})
}

// GetPost 帖子详情
// @Summary 帖子详情
// @Tags 帖子
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /posts/{id} [get]
func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	detail, err := h.postService.GetPost(uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeNotFound, err.Error())
		return
	}

	response.Success(c, detail)
}

// GetPostList 帖子列表
// @Summary 帖子列表
// @Tags 帖子
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category_id query uint false "分类ID"
// @Param sort query string false "排序: hot, essence, newest"
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response
// @Router /posts [get]
func (h *PostHandler) GetPostList(c *gin.Context) {
	var req model.PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.postService.GetPostList(&req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}

// UpdatePost 编辑帖子
// @Summary 编辑帖子
// @Tags 帖子
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path uint true "帖子ID"
// @Param body body model.UpdatePostRequest true "更新内容"
// @Success 200 {object} response.Response
// @Router /posts/{id} [put]
func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	var req model.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	err = h.postService.UpdatePost(userID.(uint), uint(postID), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeletePost 删除帖子
// @Summary 删除帖子
// @Tags 帖子
// @Security BearerAuth
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	err = h.postService.DeletePost(userID.(uint), uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// LikePost 点赞帖子
// @Summary 点赞
// @Tags 帖子
// @Security BearerAuth
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /posts/{id}/like [post]
func (h *PostHandler) LikePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	err = h.postService.LikePost(userID.(uint), uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnlikePost 取消点赞
// @Summary 取消点赞
// @Tags 帖子
// @Security BearerAuth
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /posts/{id}/like [delete]
func (h *PostHandler) UnlikePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	err = h.postService.UnlikePost(userID.(uint), uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// FavoritePost 收藏帖子
// @Summary 收藏
// @Tags 帖子
// @Security BearerAuth
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /posts/{id}/favorite [post]
func (h *PostHandler) FavoritePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	err = h.postService.FavoritePost(userID.(uint), uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// UnfavoritePost 取消收藏
// @Summary 取消收藏
// @Tags 帖子
// @Security BearerAuth
// @Produce json
// @Param id path uint true "帖子ID"
// @Success 200 {object} response.Response
// @Router /posts/{id}/favorite [delete]
func (h *PostHandler) UnfavoritePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	err = h.postService.UnfavoritePost(userID.(uint), uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetUserPosts 获取用户帖子
// @Summary 用户帖子列表
// @Tags 帖子
// @Produce json
// @Param user_id path uint true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /users/{user_id}/posts [get]
func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	var req model.PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.postService.GetUserPosts(uint(userID), &req)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}

// GetUserFavorites 获取用户收藏的帖子
// @Summary 用户收藏列表
// @Tags 帖子
// @Produce json
// @Param user_id path uint true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /users/{user_id}/favorites [get]
func (h *PostHandler) GetUserFavorites(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	var req model.PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.postService.GetUserFavorites(uint(userID), req.Page, req.PageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}

// GetUserLikedPosts 获取用户点赞的帖子
// @Summary 用户点赞列表
// @Tags 帖子
// @Produce json
// @Param user_id path uint true "用户ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /users/{user_id}/likes [get]
func (h *PostHandler) GetUserLikedPosts(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的用户ID")
		return
	}

	var req model.PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.postService.GetUserLikedPosts(uint(userID), req.Page, req.PageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}

// GetFollowFeed 获取关注用户的动态
// @Summary 关注动态
// @Tags 帖子
// @Security BearerAuth
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response
// @Router /posts/feed [get]
func (h *PostHandler) GetFollowFeed(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, response.CodeUnauthorized)
		return
	}

	var req model.PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "参数错误")
		return
	}

	list, err := h.postService.GetFollowFeed(userID.(uint), req.Page, req.PageSize)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, list)
}
