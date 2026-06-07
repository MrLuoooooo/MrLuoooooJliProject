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

type NotificationService struct {
	notifRepo repository.NotificationRepository
	userRepo  repository.UserRepository
	imClient  im.IMClient
	wsManager *ws.Manager
}

func NewNotificationService(
	notifRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
	imClient im.IMClient,
	wsManager *ws.Manager,
) *NotificationService {
	return &NotificationService{notifRepo: notifRepo, userRepo: userRepo, imClient: imClient, wsManager: wsManager}
}

func (s *NotificationService) Create(userID, fromID uint, ntype int, targetID uint, content string) error {
	notif := mysql.Notification{
		UserID: userID, FromID: fromID, Type: ntype,
		TargetID: targetID, Content: truncateStr(content, 500), IsRead: false,
	}
	if err := s.notifRepo.Create(&notif); err != nil {
		zap.S().Error("创建通知失败", "userID", userID, "error", err)
		return err
	}
	return nil
}

func (s *NotificationService) GetList(userID uint, page, pageSize int) (*model.NotificationListResponse, error) {
	notifs, total, err := s.notifRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}
	unread, _ := s.notifRepo.CountUnread(userID)

	items := make([]model.NotificationResponse, 0, len(notifs))
	if len(notifs) > 0 {
		fromIDs := make([]uint, 0)
		for _, n := range notifs {
			if n.FromID > 0 {
				fromIDs = append(fromIDs, n.FromID)
			}
		}
		var fromUsers []mysql.User
		if len(fromIDs) > 0 {
			fromUsers, _ = s.userRepo.FindByIDs(fromIDs)
		}
		fromUserMap := make(map[uint]mysql.User)
		for _, u := range fromUsers {
			fromUserMap[u.ID] = u
		}
		for _, n := range notifs {
			item := model.NotificationResponse{
				ID: n.ID, UserID: n.UserID, FromID: n.FromID, Type: n.Type,
				TargetID: n.TargetID, Content: n.Content, IsRead: n.IsRead,
				CreatedAt: n.CreatedAt.Format("2006-01-02 15:04:05"),
			}
			if n.FromID > 0 {
				if fromUser, ok := fromUserMap[n.FromID]; ok {
					item.FromName = fromUser.Nickname
					if item.FromName == "" {
						item.FromName = fromUser.Username
					}
					item.FromAvatar = fromUser.Avatar
				}
			}
			items = append(items, item)
		}
	}
	return &model.NotificationListResponse{Total: total, Items: items, Unread: unread}, nil
}

func (s *NotificationService) MarkRead(notifID, userID uint) error {
	if err := s.notifRepo.MarkRead(notifID, userID); err != nil {
		return errors.New("通知不存在")
	}
	return nil
}

func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.notifRepo.MarkAllRead(userID)
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	return s.notifRepo.CountUnread(userID)
}

// CreateAndPush 写入通知并通过 IM 异步推送给用户，推送失败不影响主流程
func (s *NotificationService) CreateAndPush(userID, fromID uint, ntype int, targetID uint, content string) {
	if err := s.Create(userID, fromID, ntype, targetID, content); err != nil {
		zap.S().Warn("写入通知失败", "error", err)
		return
	}
	// IM 异步推送，不阻塞调用方
	go func() {
		senderID := im.UserIDToStr(0) // system sender
		targetIDStr := im.UserIDToStr(userID)
		if err := s.imClient.SendSystemMsg(senderID, targetIDStr, content); err != nil {
			zap.S().Warn("IM推送通知失败", "userID", userID, "error", err)
		}
	}()
	// WebSocket 实时推送
	s.wsManager.SendToUser(userID, ws.PushMessage{
		Type: "notification",
		Data: map[string]interface{}{
			"from_id":  fromID,
			"type":     ntype,
			"target": targetID,
			"content":  content,
		},
	})
}

func truncateStr(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen])
	}
	return s
}
