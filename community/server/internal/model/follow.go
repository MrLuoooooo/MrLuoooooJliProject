package model

type FollowRequest struct {
	FollowID uint `json:"follow_id" binding:"required"`
}

type FollowListRequest struct {
	UserID   uint `form:"user_id" binding:"required"`
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=20" binding:"max=100"`
}

type FollowListResponse struct {
	Total int64            `json:"total"`
	Items []FollowUserInfo `json:"items"`
}

type FollowUserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Bio      string `json:"bio"`
}
