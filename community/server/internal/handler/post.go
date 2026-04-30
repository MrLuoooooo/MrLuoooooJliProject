package handler

import (
	"strconv"

	"community-server/DB/mysql"
	"community-server/internal/model"
	"community-server/internal/service"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
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

func (h *PostHandler) GetPost(c *gin.Context) {
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.ErrorWithMsg(c, response.CodeInvalidParam, "无效的帖子ID")
		return
	}

	post, err := h.postService.GetPostByID(uint(postID))
	if err != nil {
		response.ErrorWithMsg(c, response.CodeNotFound, err.Error())
		return
	}

	var user mysql.User
	mysql.DB.Where("id = ?", post.UserID).First(&user)

	response.Success(c, model.PostResponse{
		ID:           post.ID,
		UserID:       post.UserID,
		Username:     user.Username,
		Nickname:     user.Nickname,
		Title:        post.Title,
		Content:      post.Content,
		Summary:      post.Summary,
		CoverImage:   post.CoverImage,
		CategoryID:   post.CategoryID,
		Status:       post.Status,
		ViewCount:    post.ViewCount,
		LikeCount:    post.LikeCount,
		CommentCount: post.CommentCount,
		IsTop:        post.IsTop,
		IsEssence:    post.IsEssence,
		CreatedAt:    post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    post.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

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

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
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

func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
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

func (h *PostHandler) LikePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
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

func (h *PostHandler) UnlikePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
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

func (h *PostHandler) FavoritePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
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

func (h *PostHandler) UnfavoritePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorWithMsg(c, response.CodeUnauthorized, "未登录")
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
