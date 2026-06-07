package handler

import (
	"mime/multipart"

	"community-server/internal/model"
	"community-server/internal/service"
)

// UploadService 上传服务接口
type UploadService interface {
	Upload(file *multipart.FileHeader) (*service.UploadResult, error)
}

// UserService 用户服务，定义 handler 依赖的接口供测试 mock 用
type UserService interface {
	Register(req *model.RegisterRequest) (uint, error)
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
	GetUserByID(userID uint) (*model.UserProfileResponse, error)
	UpdateProfile(userID uint, req *model.UpdateProfileRequest) error
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
}

// PostService 帖子服务接口
type PostService interface {
	CreatePost(userID uint, req *model.CreatePostRequest) (uint, error)
	GetPostList(req *model.PostListRequest) (*model.PostListResponse, error)
	GetPost(postID uint) (*model.PostDetailResponse, error)
	UpdatePost(userID, postID uint, req *model.UpdatePostRequest) error
	DeletePost(userID, postID uint) error
	LikePost(userID, postID uint) error
	UnlikePost(userID, postID uint) error
	FavoritePost(userID, postID uint) error
	UnfavoritePost(userID, postID uint) error
	GetUserPosts(userID uint, req *model.PostListRequest) (*model.PostListResponse, error)
	GetUserFavorites(userID uint, page, pageSize int) (*model.PostListResponse, error)
	GetUserLikedPosts(userID uint, page, pageSize int) (*model.PostListResponse, error)
	GetFollowFeed(userID uint, page, pageSize int) (*model.PostListResponse, error)
}

// CommentService 评论服务接口
type CommentService interface {
	CreateComment(userID uint, req *model.CreateCommentRequest) (uint, error)
	GetCommentList(req *model.CommentListRequest) (*model.CommentListResponse, error)
	UpdateComment(userID, commentID uint, req *model.UpdateCommentRequest) error
	DeleteComment(userID, commentID uint) error
	LikeComment(userID, commentID uint) error
	UnlikeComment(userID, commentID uint) error
}

// MessageService 消息服务接口
type MessageService interface {
	SendMessage(senderID uint, req *model.SendMessageRequest) error
	GetMessageList(req *model.MessageListRequest) (*model.MessageListResponse, error)
	GetUnreadCount(userID uint) (int64, error)
	MarkAsRead(messageID uint, userID uint) error
	GetConversationList(userID uint) (*model.ConversationListResponse, error)
}

// FollowService 关注服务接口
type FollowService interface {
	FollowUser(userID uint, req *model.FollowRequest) error
	UnfollowUser(userID, followID uint) error
	GetFollowers(userID uint, page, pageSize int) (*model.FollowListResponse, error)
	GetFollowing(userID uint, page, pageSize int) (*model.FollowListResponse, error)
	IsFollowing(userID, targetID uint) (bool, error)
	GetFollowCounts(userID uint) (followers int64, following int64)
}

// NotificationService 通知服务接口
type NotificationService interface {
	GetList(userID uint, page, pageSize int) (*model.NotificationListResponse, error)
	MarkRead(notifID, userID uint) error
	MarkAllRead(userID uint) error
	GetUnreadCount(userID uint) (int64, error)
}

// SearchService 搜索服务接口
type SearchService interface {
	SearchPosts(req *model.SearchRequest) (*model.SearchResponse, error)
	SearchUsers(req *model.SearchRequest) (*model.SearchResponse, error)
}

// AdminService 管理服务接口
type AdminService interface {
	GetUserList(req *model.AdminUserListRequest) (*model.AdminUserListResponse, error)
	DeleteUser(userID uint) error
	UpdateUserAdminType(userID uint, adminType int) error
	UpdateUserStatus(userID uint, status int) error
	GetPostList(req *model.AdminPostListRequest) (*model.AdminPostListResponse, error)
	DeletePost(postID uint) error
	SetPostTop(postID uint, isTop bool) error
	SetPostEssence(postID uint, isEssence bool) error
	GetStats() (*model.AdminStatsResponse, error)
}

// TagService 标签服务接口
type TagService interface {
	CreateTag(req *model.CreateTagRequest) (uint, error)
	GetTagList(req *model.TagListRequest) (*model.TagListResponse, error)
	UpdateTag(tagID uint, req *model.UpdateTagRequest) error
	DeleteTag(tagID uint) error
	AddPostTags(postID uint, tagIDs []uint) error
	RemovePostTag(postID, tagID uint) error
	GetPostTags(postID uint) ([]model.TagResponse, error)
	GetPostsByTag(tagID uint, page, pageSize int) (*model.PostListResponse, error)
}

// CategoryService 分类服务接口
type CategoryService interface {
	Create(req *model.CreateCategoryRequest) (uint, error)
	Update(catID uint, req *model.UpdateCategoryRequest) error
	Delete(catID uint) error
	GetList() ([]model.CategoryResponse, error)
	GetByID(catID uint) (*model.CategoryResponse, error)
}
