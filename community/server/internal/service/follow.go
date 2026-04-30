package service

import (
	"errors"

	"community-server/DB/mysql"
	"community-server/internal/model"

	"go.uber.org/zap"
)

type FollowService struct{}

func NewFollowService() *FollowService {
	return &FollowService{}
}

func (s *FollowService) FollowUser(userID uint, req *model.FollowRequest) error {
	if userID == req.FollowID {
		return errors.New("不能关注自己")
	}

	var target mysql.User
	if err := mysql.DB.First(&target, req.FollowID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var existing mysql.UserFollow
	result := mysql.DB.Where("user_id = ? AND follow_id = ?", userID, req.FollowID).First(&existing)
	if result.Error == nil {
		return errors.New("已关注该用户")
	}

	follow := mysql.UserFollow{
		UserID:   userID,
		FollowID: req.FollowID,
	}

	if err := mysql.DB.Create(&follow).Error; err != nil {
		zap.S().Error("关注失败", "userId", userID, "followId", req.FollowID, "error", err)
		return errors.New("关注失败")
	}

	zap.S().Info("关注成功", "userId", userID, "followId", req.FollowID)
	return nil
}

func (s *FollowService) UnfollowUser(userID uint, followID uint) error {
	result := mysql.DB.Where("user_id = ? AND follow_id = ?", userID, followID).Delete(&mysql.UserFollow{})
	if result.Error != nil {
		zap.S().Error("取消关注失败", "userId", userID, "followId", followID, "error", result.Error)
		return errors.New("取消关注失败")
	}

	if result.RowsAffected == 0 {
		return errors.New("未关注该用户")
	}

	zap.S().Info("取消关注成功", "userId", userID, "followId", followID)
	return nil
}

func (s *FollowService) GetFollowers(req *model.FollowListRequest) (*model.FollowListResponse, error) {
	var follows []mysql.UserFollow
	var total int64

	mysql.DB.Model(&mysql.UserFollow{}).Where("follow_id = ?", req.UserID).Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := mysql.DB.Where("follow_id = ?", req.UserID).
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&follows)

	if result.Error != nil {
		return nil, errors.New("获取粉丝列表失败")
	}

	items := make([]model.FollowUserInfo, 0, len(follows))
	for _, follow := range follows {
		var user mysql.User
		mysql.DB.Where("id = ?", follow.UserID).First(&user)

		items = append(items, model.FollowUserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Bio:      user.Bio,
		})
	}

	return &model.FollowListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *FollowService) GetFollowing(req *model.FollowListRequest) (*model.FollowListResponse, error) {
	var follows []mysql.UserFollow
	var total int64

	mysql.DB.Model(&mysql.UserFollow{}).Where("user_id = ?", req.UserID).Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := mysql.DB.Where("user_id = ?", req.UserID).
		Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&follows)

	if result.Error != nil {
		return nil, errors.New("获取关注列表失败")
	}

	items := make([]model.FollowUserInfo, 0, len(follows))
	for _, follow := range follows {
		var user mysql.User
		mysql.DB.Where("id = ?", follow.FollowID).First(&user)

		items = append(items, model.FollowUserInfo{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Bio:      user.Bio,
		})
	}

	return &model.FollowListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *FollowService) IsFollowing(userID uint, followID uint) (bool, error) {
	var follow mysql.UserFollow
	result := mysql.DB.Where("user_id = ? AND follow_id = ?", userID, followID).First(&follow)
	if result.Error != nil {
		return false, nil
	}
	return true, nil
}

func (s *FollowService) GetFollowCounts(userID uint) (followers int64, following int64) {
	mysql.DB.Model(&mysql.UserFollow{}).Where("follow_id = ?", userID).Count(&followers)
	mysql.DB.Model(&mysql.UserFollow{}).Where("user_id = ?", userID).Count(&following)
	return
}
