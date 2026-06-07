package service

import (
	"errors"

	"community-server/internal/db/mysql"
	"community-server/internal/im"
	"community-server/internal/model"
	"community-server/internal/repository"
	"community-server/internal/ws"

	"go.uber.org/zap"
)

type MessageService struct {
	imClient    im.IMClient
	msgRepo     repository.MessageRepository
	userRepo    repository.UserRepository
	wsManager   *ws.Manager
}

func NewMessageService(imClient im.IMClient, msgRepo repository.MessageRepository, userRepo repository.UserRepository, wsManager *ws.Manager) *MessageService {
	return &MessageService{imClient: imClient, msgRepo: msgRepo, userRepo: userRepo, wsManager: wsManager}
}

func (s *MessageService) SendMessage(senderID uint, req *model.SendMessageRequest) error {
	if senderID == req.ReceiverID {
		return errors.New("不能给自己发消息")
	}
	if _, err := s.userRepo.FindByID(req.ReceiverID); err != nil {
		return errors.New("用户不存在")
	}
	message := mysql.Message{
		SenderID: senderID, ReceiverID: req.ReceiverID, Content: req.Content, IsRead: false,
	}
	if err := s.msgRepo.Create(&message); err != nil {
		zap.S().Error("发送消息失败", "senderId", senderID, "receiverId", req.ReceiverID, "error", err)
		return errors.New("发送消息失败")
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				zap.S().Error("IM推送消息 panic", "error", r)
			}
		}()
		if err := s.imClient.SendPrivateMsg(im.UserIDToStr(senderID), im.UserIDToStr(req.ReceiverID), req.Content); err != nil {
			zap.S().Warn("IM推送消息失败（不影响本地存储）", "error", err)
		}
	}()
	zap.S().Info("发送消息成功", "senderId", senderID, "receiverId", req.ReceiverID, "messageId", message.ID)
	s.wsManager.SendToUser(req.ReceiverID, ws.PushMessage{
		Type: "message",
		Data: map[string]interface{}{
			"message_id":  message.ID,
			"sender_id":   senderID,
			"receiver_id": req.ReceiverID,
			"content":     req.Content,
		},
	})
	return nil
}

func (s *MessageService) GetMessageList(req *model.MessageListRequest) (*model.MessageListResponse, error) {
	var messages []mysql.Message
	var total int64
	var err error
	if req.SenderID > 0 && req.ReceiverID > 0 {
		messages, total, err = s.msgRepo.FindByConversation(req.SenderID, req.ReceiverID, req.Page, req.PageSize)
	} else {
		messages, total, err = s.msgRepo.FindReceived(req.ReceiverID, req.Page, req.PageSize)
	}
	if err != nil {
		return nil, errors.New("获取消息列表失败")
	}
	items := make([]model.MessageInfo, 0, len(messages))
	if len(messages) > 0 {
		userIDs := make([]uint, len(messages))
		for i, msg := range messages {
			userIDs[i] = msg.SenderID
		}
		users, _ := s.userRepo.FindByIDs(userIDs)
		userMap := make(map[uint]mysql.User, len(users))
		for _, u := range users {
			userMap[u.ID] = u
		}
		for _, msg := range messages {
			sender := userMap[msg.SenderID]
			createdAt := ""
			if !msg.CreatedAt.IsZero() {
				createdAt = msg.CreatedAt.Format("2006-01-02 15:04:05")
			}
			items = append(items, model.MessageInfo{
				ID: msg.ID, SenderID: msg.SenderID, ReceiverID: msg.ReceiverID,
				SenderName: sender.Username, Content: msg.Content, IsRead: msg.IsRead, CreatedAt: createdAt,
			})
		}
	}
	return &model.MessageListResponse{Total: total, Items: items}, nil
}

func (s *MessageService) GetUnreadCount(userID uint) (int64, error) {
	return s.msgRepo.CountUnread(userID)
}

func (s *MessageService) MarkAsRead(messageID, userID uint) error {
	return s.msgRepo.MarkAsRead(messageID, userID)
}

func (s *MessageService) GetConversationList(userID uint) (*model.ConversationListResponse, error) {
	messages, _ := s.msgRepo.FindConversationMessages(userID)
	conversations := make(map[uint]*model.ConversationInfo)
	otherIDs := make([]uint, 0)

	for _, msg := range messages {
		var otherID uint
		if msg.SenderID == userID {
			otherID = msg.ReceiverID
		} else {
			otherID = msg.SenderID
		}
		if _, exists := conversations[otherID]; !exists {
			conversations[otherID] = &model.ConversationInfo{UserID: otherID}
			otherIDs = append(otherIDs, otherID)
		}
	}

	if len(otherIDs) > 0 {
		users, _ := s.userRepo.FindByIDs(otherIDs)
		userMap := make(map[uint]mysql.User, len(users))
		for _, u := range users {
			userMap[u.ID] = u
		}
		unreadMap, _ := s.msgRepo.CountUnreadBySenders(userID, otherIDs)

		for _, msg := range messages {
			var otherID uint
			if msg.SenderID == userID {
				otherID = msg.ReceiverID
			} else {
				otherID = msg.SenderID
			}
			conv := conversations[otherID]
			if u, ok := userMap[otherID]; ok && conv.Username == "" {
				conv.Username = u.Username; conv.Nickname = u.Nickname; conv.Avatar = u.Avatar
			}
			conv.UnreadCount = unreadMap[otherID]
			conv.LastMessage = msg.Content
			if !msg.CreatedAt.IsZero() {
				conv.LastTime = msg.CreatedAt.Format("2006-01-02 15:04:05")
			}
		}
	}

	items := make([]model.ConversationInfo, 0, len(conversations))
	for _, conv := range conversations {
		items = append(items, *conv)
	}
	return &model.ConversationListResponse{Items: items}, nil
}
