package model

// 通知类型常量
const (
	NotifyComment = 1 // 评论/回复
	NotifyLike    = 2 // 点赞
	NotifyFollow  = 3 // 关注
	NotifySystem  = 4 // 系统通知
)

type NotificationListRequest struct {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=20" binding:"max=100"`
}

type NotificationResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	FromID    uint   `json:"from_id"`
	FromName  string `json:"from_name"`
	FromAvatar string `json:"from_avatar"`
	Type      int    `json:"type"`
	TargetID  uint   `json:"target_id"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type NotificationListResponse struct {
	Total int64                 `json:"total"`
	Items []NotificationResponse `json:"items"`
	Unread int64                `json:"unread"`
}
