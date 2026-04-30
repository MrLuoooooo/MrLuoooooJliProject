package model

type CreateCommentRequest struct {
	PostID   uint   `json:"post_id" binding:"required"`
	ParentID uint   `json:"parent_id"`
	Content  string `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentListRequest struct {
	PostID   uint `form:"post_id" binding:"required"`
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=20"`
}

type CommentResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	PostID    uint   `json:"post_id"`
	ParentID  uint   `json:"parent_id"`
	Content   string `json:"content"`
	Status    int    `json:"status"`
	LikeCount int    `json:"like_count"`
	CreatedAt string `json:"created_at"`
}

type CommentListItem struct {
	ID         uint            `json:"id"`
	UserID     uint            `json:"user_id"`
	Username   string          `json:"username"`
	Nickname   string          `json:"nickname"`
	Avatar     string          `json:"avatar"`
	PostID     uint            `json:"post_id"`
	ParentID   uint            `json:"parent_id"`
	Content    string          `json:"content"`
	Status     int             `json:"status"`
	LikeCount  int             `json:"like_count"`
	CreatedAt  string          `json:"created_at"`
	Replies    []CommentListItem `json:"replies,omitempty"`
}

type CommentListResponse struct {
	Total int64             `json:"total"`
	Items []CommentListItem `json:"items"`
}
