package model

type AdminUserListRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type AdminUserListResponse struct {
	Total int64           `json:"total"`
	Items []AdminUserInfo `json:"items"`
}

type AdminUserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	AdminType int    `json:"admin_type"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
}

type AdminPostListRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20"`
}

type AdminPostListResponse struct {
	Total int64           `json:"total"`
	Items []AdminPostInfo `json:"items"`
}

type AdminPostInfo struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	Summary      string `json:"summary"`
	Status       int    `json:"status"`
	ViewCount    int    `json:"view_count"`
	LikeCount    int    `json:"like_count"`
	CommentCount int    `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
}

type UpdateUserAdminTypeRequest struct {
	AdminType int `json:"admin_type" binding:"required,oneof=0 1"`
}

type UpdateUserStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=0 1"`
}
