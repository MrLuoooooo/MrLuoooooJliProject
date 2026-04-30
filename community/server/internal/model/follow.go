package model

type FollowRequest struct {
	FollowID uint `json:"follow_id" binding:"required"`
}

type FollowListRequest struct {
	UserID   uint `form:"user_id" binding:"required"`
	Page     int  `form:"page,default=1"`
	PageSize int  `form:"page_size,default=20"`
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

type SendMessageRequest struct {
	ReceiverID uint   `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required,max=1000"`
}

type MessageListRequest struct {
	SenderID   uint `form:"sender_id"`
	ReceiverID uint `form:"receiver_id"`
	Page       int  `form:"page,default=1"`
	PageSize   int  `form:"page_size,default=20"`
}

type MessageListResponse struct {
	Total int64             `json:"total"`
	Items []MessageInfo     `json:"items"`
}

type MessageInfo struct {
	ID         uint   `json:"id"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
	IsRead     bool   `json:"is_read"`
	CreatedAt  string `json:"created_at"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

type ConversationListResponse struct {
	Items []ConversationInfo `json:"items"`
}

type ConversationInfo struct {
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Avatar       string `json:"avatar"`
	LastMessage  string `json:"last_message"`
	LastTime     string `json:"last_time"`
	UnreadCount  int64  `json:"unread_count"`
}
