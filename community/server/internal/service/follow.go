package service

import (
	"errors"

	"community-server/internal/db/mysql"
	"community-server/internal/model"
	"community-server/internal/repository"

	"go.uber.org/zap"
)

type FollowService struct {
	notifSvc   *NotificationService
	followRepo repository.FollowRepository
	userRepo   repository.UserRepository
}

func NewFollowService(notifSvc *NotificationService, followRepo repository.FollowRepository, userRepo repository.UserRepository) *FollowService {
	return &FollowService{notifSvc: notifSvc, followRepo: followRepo, userRepo: userRepo}
}

func (s *FollowService) FollowUser(userID uint, req *model.FollowRequest) error {
	if userID == req.FollowID {
		return errors.New("不能关注自己")
	}
	if _, err := s.userRepo.FindByID(req.FollowID); err != nil {
		return errors.New("用户不存在")
	}
	exists, _ := s.followRepo.Exists(userID, req.FollowID)
	if exists {
		return errors.New("已关注该用户")
	}
	if err := s.followRepo.Create(&mysql.UserFollow{UserID: userID, FollowID: req.FollowID}); err != nil {
		zap.S().Error("关注失败", "userId", userID, "followId", req.FollowID, "error", err)
		return errors.New("关注失败")
	}
	zap.S().Info("关注成功", "userId", userID, "followId", req.FollowID)
	if s.notifSvc != nil {
		s.notifSvc.CreateAndPush(req.FollowID, userID, model.NotifyFollow, 0, "关注了你")
	}
	return nil
}

func (s *FollowService) UnfollowUser(userID, followID uint) error {
	exists, _ := s.followRepo.Exists(userID, followID)
	if !exists {
		return errors.New("未关注该用户")
	}
	if err := s.followRepo.Delete(userID, followID); err != nil {
		zap.S().Error("取消关注失败", "userId", userID, "followId", followID, "error", err)
		return errors.New("取消关注失败")
	}
	zap.S().Info("取消关注成功", "userId", userID, "followId", followID)
	return nil
}

func (s *FollowService) GetFollowers(userID uint, page, pageSize int) (*model.FollowListResponse, error) {
	follows, total, err := s.followRepo.FindFollowers(userID, page, pageSize)
	if err != nil {
		return nil, errors.New("获取粉丝列表失败")
	}
	return s.buildFollowList(follows, total, true), nil
}

func (s *FollowService) GetFollowing(userID uint, page, pageSize int) (*model.FollowListResponse, error) {
	follows, total, err := s.followRepo.FindFollowing(userID, page, pageSize)
	if err != nil {
		return nil, errors.New("获取关注列表失败")
	}
	return s.buildFollowList(follows, total, false), nil
}

func (s *FollowService) IsFollowing(userID, followID uint) (bool, error) {
	return s.followRepo.Exists(userID, followID)
}

func (s *FollowService) GetFollowCounts(userID uint) (int64, int64) {
	return s.followRepo.CountFollowers(userID), s.followRepo.CountFollowing(userID)
}

func (s *FollowService) buildFollowList(follows []mysql.UserFollow, total int64, isFollowers bool) *model.FollowListResponse {
	items := make([]model.FollowUserInfo, 0, len(follows))
	if len(follows) == 0 {
		return &model.FollowListResponse{Total: total, Items: items}
	}
	userIDs := make([]uint, len(follows))
	for i, f := range follows {
		if isFollowers {
			userIDs[i] = f.UserID
		} else {
			userIDs[i] = f.FollowID
		}
	}
	users, _ := s.userRepo.FindByIDs(userIDs)
	userMap := make(map[uint]mysql.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}
	for _, f := range follows {
		var uid uint
		if isFollowers {
			uid = f.UserID
		} else {
			uid = f.FollowID
		}
		u := userMap[uid]
		items = append(items, model.FollowUserInfo{
			ID: u.ID, Username: u.Username, Nickname: u.Nickname, Avatar: u.Avatar, Bio: u.Bio,
		})
	}
	return &model.FollowListResponse{Total: total, Items: items}
}
