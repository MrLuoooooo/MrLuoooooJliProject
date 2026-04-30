package mysql

import (
	"time"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	BaseModel
	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone     string     `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Nickname  string     `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string     `gorm:"type:varchar(255)" json:"avatar"`
	Bio       string     `gorm:"type:varchar(500)" json:"bio"`
	Status    int        `gorm:"default:1" json:"status"`
	LastLogin *time.Time `json:"last_login"`
}

func (User) TableName() string {
	return "users"
}

type Post struct {
	BaseModel
	UserID       uint   `gorm:"index;not null" json:"user_id"`
	Title        string `gorm:"type:varchar(200);not null" json:"title"`
	Content      string `gorm:"type:text;not null" json:"content"`
	Summary      string `gorm:"type:varchar(500)" json:"summary"`
	CoverImage   string `gorm:"type:varchar(255)" json:"cover_image"`
	CategoryID   uint   `gorm:"index" json:"category_id"`
	Status       int    `gorm:"default:1" json:"status"`
	ViewCount    int    `gorm:"default:0" json:"view_count"`
	LikeCount    int    `gorm:"default:0" json:"like_count"`
	CommentCount int    `gorm:"default:0" json:"comment_count"`
	IsTop        bool   `gorm:"default:false" json:"is_top"`
	IsEssence    bool   `gorm:"default:false" json:"is_essence"`
}

func (Post) TableName() string {
	return "posts"
}

type PostLike struct {
	BaseModel
	UserID uint `gorm:"uniqueIndex:idx_user_post_like;not null" json:"user_id"`
	PostID uint `gorm:"uniqueIndex:idx_user_post_like;not null" json:"post_id"`
}

func (PostLike) TableName() string {
	return "post_likes"
}

type PostFavorite struct {
	BaseModel
	UserID uint `gorm:"uniqueIndex:idx_user_post_favorite;not null" json:"user_id"`
	PostID uint `gorm:"uniqueIndex:idx_user_post_favorite;not null" json:"post_id"`
}

func (PostFavorite) TableName() string {
	return "post_favorites"
}

type Comment struct {
	BaseModel
	UserID    uint   `gorm:"index;not null" json:"user_id"`
	PostID    uint   `gorm:"index;not null" json:"post_id"`
	ParentID  uint   `gorm:"default:0" json:"parent_id"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Status    int    `gorm:"default:1" json:"status"`
	LikeCount int    `gorm:"default:0" json:"like_count"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentLike struct {
	BaseModel
	UserID    uint `gorm:"uniqueIndex:idx_user_comment_like;not null" json:"user_id"`
	CommentID uint `gorm:"uniqueIndex:idx_user_comment_like;not null" json:"comment_id"`
}

func (CommentLike) TableName() string {
	return "comment_likes"
}

type Tag struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:varchar(200)" json:"description"`
	PostCount   int    `gorm:"default:0" json:"post_count"`
	Status      int    `gorm:"default:1" json:"status"`
}

func (Tag) TableName() string {
	return "tags"
}

type PostTag struct {
	BaseModel
	PostID uint `gorm:"uniqueIndex:idx_post_tag;not null" json:"post_id"`
	TagID  uint `gorm:"uniqueIndex:idx_post_tag;not null" json:"tag_id"`
}

func (PostTag) TableName() string {
	return "post_tags"
}
