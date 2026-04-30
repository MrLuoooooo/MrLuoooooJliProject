package model

type CreatePostRequest struct {
	Title      string `json:"title" binding:"required,max=200"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"max=255"`
	CategoryID uint   `json:"category_id"`
	Status     int    `json:"status"`
}

type UpdatePostRequest struct {
	Title      string `json:"title" binding:"max=200"`
	Content    string `json:"content"`
	Summary    string `json:"summary" binding:"max=500"`
	CoverImage string `json:"cover_image" binding:"max=255"`
	CategoryID uint   `json:"category_id"`
	Status     int    `json:"status"`
}

type PostListRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type PostResponse struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	Summary      string `json:"summary"`
	CoverImage   string `json:"cover_image"`
	CategoryID   uint   `json:"category_id"`
	Status       int    `json:"status"`
	ViewCount    int    `json:"view_count"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	IsTop        bool   `json:"is_top"`
	IsEssence    bool   `json:"is_essence"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type PostListItem struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Title        string `json:"title"`
	Summary      string `json:"summary"`
	CoverImage   string `json:"cover_image"`
	CategoryID   uint   `json:"category_id"`
	Status       int    `json:"status"`
	ViewCount    int    `json:"view_count"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	IsTop        bool   `json:"is_top"`
	IsEssence    bool   `json:"is_essence"`
	CreatedAt    string `json:"created_at"`
}

type PostListResponse struct {
	Total int64          `json:"total"`
	Items []PostListItem `json:"items"`
}
