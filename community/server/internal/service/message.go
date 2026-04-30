package service

import (
	"errors"

	"community-server/DB/mysql"
	"community-server/internal/model"

	"go.uber.org/zap"
)

type MessageService struct{}

func NewMessageService() *MessageService {
	return &MessageService{}
}

func (s *MessageService) SendMessage(senderID uint, req *model.SendMessageRequest) error {
	if senderID == req.ReceiverID {
		return errors.New("不能给自己发消息")
	}

	var receiver mysql.User
	if err := mysql.DB.First(&receiver, req.ReceiverID).Error; err != nil {
		return errors.New("用户不存在")
	}

	message := mysql.Message{
		SenderID:   senderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
		IsRead:     false,
	}

	if err := mysql.DB.Create(&message).Error; err != nil {
		zap.S().Error("发送消息失败", "senderId", senderID, "receiverId", req.ReceiverID, "error", err)
		return errors.New("发送消息失败")
	}

	zap.S().Info("发送消息成功", "senderId", senderID, "receiverId", req.ReceiverID, "messageId", message.ID)
	return nil
}

func (s *MessageService) GetMessageList(req *model.MessageListRequest) (*model.MessageListResponse, error) {
	var messages []mysql.Message
	var total int64

	query := mysql.DB.Model(&mysql.Message{})
	if req.SenderID > 0 && req.ReceiverID > 0 {
		query = query.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			req.SenderID, req.ReceiverID, req.ReceiverID, req.SenderID)
	} else if req.ReceiverID > 0 {
		query = query.Where("receiver_id = ?", req.ReceiverID)
	}

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&messages)

	if result.Error != nil {
		return nil, errors.New("获取消息列表失败")
	}

	items := make([]model.MessageInfo, 0, len(messages))
	for _, msg := range messages {
		var sender mysql.User
		mysql.DB.Where("id = ?", msg.SenderID).First(&sender)

		createdAt := ""
		if !msg.CreatedAt.IsZero() {
			createdAt = msg.CreatedAt.Format("2006-01-02 15:04:05")
		}

		items = append(items, model.MessageInfo{
			ID:         msg.ID,
			SenderID:   msg.SenderID,
			ReceiverID: msg.ReceiverID,
			SenderName: sender.Username,
			Content:    msg.Content,
			IsRead:     msg.IsRead,
			CreatedAt:  createdAt,
		})
	}

	return &model.MessageListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *MessageService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	result := mysql.DB.Model(&mysql.Message{}).
		Where("receiver_id = ? AND is_read = ?", userID, false).
		Count(&count)

	if result.Error != nil {
		return 0, errors.New("获取未读消息数失败")
	}

	return count, nil
}

func (s *MessageService) MarkAsRead(messageID uint, userID uint) error {
	result := mysql.DB.Model(&mysql.Message{}).
		Where("id = ? AND receiver_id = ?", messageID, userID).
		Update("is_read", true)

	if result.Error != nil {
		return errors.New("标记已读失败")
	}

	return nil
}

func (s *MessageService) GetConversationList(userID uint) (*model.ConversationListResponse, error) {
	var messages []mysql.Message
	mysql.DB.Where("sender_id = ? OR receiver_id = ?", userID, userID).
		Order("created_at DESC").
		Find(&messages)

	conversations := make(map[uint]*model.ConversationInfo)
	for _, msg := range messages {
		var otherID uint
		if msg.SenderID == userID {
			otherID = msg.ReceiverID
		} else {
			otherID = msg.SenderID
		}

		if _, exists := conversations[otherID]; !exists {
			var user mysql.User
			mysql.DB.Where("id = ?", otherID).First(&user)

			createdAt := ""
			if !msg.CreatedAt.IsZero() {
				createdAt = msg.CreatedAt.Format("2006-01-02 15:04:05")
			}

			var unreadCount int64
			mysql.DB.Model(&mysql.Message{}).
				Where("sender_id = ? AND receiver_id = ? AND is_read = ?", otherID, userID, false).
				Count(&unreadCount)

			conversations[otherID] = &model.ConversationInfo{
				UserID:      otherID,
				Username:    user.Username,
				Nickname:    user.Nickname,
				Avatar:      user.Avatar,
				LastMessage: msg.Content,
				LastTime:    createdAt,
				UnreadCount: unreadCount,
			}
		}
	}

	items := make([]model.ConversationInfo, 0, len(conversations))
	for _, conv := range conversations {
		items = append(items, *conv)
	}

	return &model.ConversationListResponse{
		Items: items,
	}, nil
}
