package repository

import (
	"community-server/internal/db/mysql"
)

// ============================================
// UserRepository
// ============================================
type UserRepository interface {
	Create(user *mysql.User) error
	FindByUsername(username string) (*mysql.User, error)
	FindByEmail(email string) (*mysql.User, error)
	FindByID(id uint) (*mysql.User, error)
	Update(userID uint, updates map[string]interface{}) error
	Delete(userID uint) error
	FindByIDs(ids []uint) ([]mysql.User, error)
	Search(keyword string, page, pageSize int) ([]mysql.User, int64, error)
	List(page, pageSize int) ([]mysql.User, int64, error)
	Count() (int64, error)
	CountToday() (int64, error)
}

// ============================================
// PostRepository
// ============================================
type PostRepository interface {
	Create(post *mysql.Post) error
	FindByID(id uint) (*mysql.Post, error)
	Update(postID uint, updates map[string]interface{}) error
	UpdateColumn(postID uint, column string, value interface{}) error
	List(req PostListQuery) ([]mysql.Post, int64, error)
	Search(keyword string, page, pageSize int) ([]mysql.Post, int64, error)
	Delete(postID uint) error
	SoftDelete(postID uint) error
	Count() (int64, error)
	CountToday() (int64, error)
}

type PostListQuery struct {
	CategoryID uint
	Keyword    string
	Sort       string
	Page       int
	PageSize   int
	UserID     uint // filter by user ID
	UserIDs    []uint // filter by multiple user IDs (for feed)
	StatusExclude int // exclude status value, 0 = no filter
}

// ============================================
// PostLikeRepository
// ============================================
type PostLikeRepository interface {
	Create(like *mysql.PostLike) error
	Delete(userID, postID uint) error
	Exists(userID, postID uint) (bool, error)
	FindByUserID(userID uint, page, pageSize int) ([]mysql.PostLike, int64, error)
}

// ============================================
// PostFavoriteRepository
// ============================================
type PostFavoriteRepository interface {
	Create(fav *mysql.PostFavorite) error
	Delete(userID, postID uint) error
	Exists(userID, postID uint) (bool, error)
	FindByUserID(userID uint, page, pageSize int) ([]mysql.PostFavorite, int64, error)
}

// ============================================
// CommentRepository
// ============================================
type CommentRepository interface {
	Create(comment *mysql.Comment) error
	FindByID(id uint) (*mysql.Comment, error)
	FindRootByPostID(postID uint, page, pageSize int) ([]mysql.Comment, int64, error)
	FindRepliesByPostAndParents(postID uint, parentIDs []uint) ([]mysql.Comment, error)
	Update(commentID uint, updates map[string]interface{}) error
	UpdateColumn(commentID uint, column string, value interface{}) error
	SoftDelete(commentID uint) error
	SoftDeleteByPostID(postID uint) error
	Count() (int64, error)
}

// ============================================
// CommentLikeRepository
// ============================================
type CommentLikeRepository interface {
	Create(like *mysql.CommentLike) error
	Delete(userID, commentID uint) error
	Exists(userID, commentID uint) (bool, error)
}

// ============================================
// TagRepository
// ============================================
type TagRepository interface {
	Create(tag *mysql.Tag) error
	FindByID(id uint) (*mysql.Tag, error)
	FindByName(name string) (*mysql.Tag, error)
	FindByIDs(ids []uint) ([]mysql.Tag, error)
	Update(tagID uint, updates map[string]interface{}) error
	Delete(tagID uint) error
	List(page, pageSize int) ([]mysql.Tag, int64, error)
	IncrementPostCount(tagID uint, delta int) error
}

// ============================================
// PostTagRepository
// ============================================
type PostTagRepository interface {
	Create(pt *mysql.PostTag) error
	Delete(postID, tagID uint) error
	FindByPostID(postID uint) ([]mysql.PostTag, error)
	FindByTagID(tagID uint, page, pageSize int) ([]mysql.PostTag, int64, error)
	DeleteByPostID(postID uint) error
}

// ============================================
// CategoryRepository
// ============================================
type CategoryRepository interface {
	Create(cat *mysql.Category) error
	FindByID(id uint) (*mysql.Category, error)
	Update(catID uint, updates map[string]interface{}) error
	Delete(catID uint) error
	List() ([]mysql.Category, error)
	IncrementPostCount(catID uint, delta int) error
}

// ============================================
// FollowRepository
// ============================================
type FollowRepository interface {
	Create(follow *mysql.UserFollow) error
	Delete(userID, followID uint) error
	Exists(userID, followID uint) (bool, error)
	FindFollowers(userID uint, page, pageSize int) ([]mysql.UserFollow, int64, error)
	FindFollowing(userID uint, page, pageSize int) ([]mysql.UserFollow, int64, error)
	CountFollowers(userID uint) int64
	CountFollowing(userID uint) int64
	FindAllFollowing(userID uint) ([]mysql.UserFollow, error)
}

// ============================================
// MessageRepository
// ============================================
type MessageRepository interface {
	Create(msg *mysql.Message) error
	FindByConversation(senderID, receiverID uint, page, pageSize int) ([]mysql.Message, int64, error)
	FindReceived(userID uint, page, pageSize int) ([]mysql.Message, int64, error)
	CountUnread(userID uint) (int64, error)
	MarkAsRead(messageID, userID uint) error
	FindConversationMessages(userID uint) ([]mysql.Message, error)
	CountUnreadBySenders(receiverID uint, senderIDs []uint) (map[uint]int64, error)
}

// ============================================
// NotificationRepository
// ============================================
type NotificationRepository interface {
	Create(notif *mysql.Notification) error
	FindByUserID(userID uint, page, pageSize int) ([]mysql.Notification, int64, error)
	CountUnread(userID uint) (int64, error)
	MarkRead(notifID, userID uint) error
	MarkAllRead(userID uint) error
}

// ============================================
// PasswordResetRepository
// ============================================
type PasswordResetRepository interface {
	Create(reset *mysql.PasswordReset) error
	FindByToken(token string) (*mysql.PasswordReset, error)
	MarkUsed(id uint) error
}
