package service

import (
	"errors"

	"community-server/DB/mysql"
	"community-server/internal/model"

	"go.uber.org/zap"
)

type AdminService struct{}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) GetUserList(req *model.AdminUserListRequest) (*model.AdminUserListResponse, error) {
	var users []mysql.User
	var total int64

	mysql.DB.Model(&mysql.User{}).Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := mysql.DB.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&users)

	if result.Error != nil {
		return nil, errors.New("获取用户列表失败")
	}

	items := make([]model.AdminUserInfo, 0, len(users))
	for _, user := range users {
		createdAt := ""
		if !user.CreatedAt.IsZero() {
			createdAt = user.CreatedAt.Format("2006-01-02 15:04:05")
		}
		items = append(items, model.AdminUserInfo{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Phone:     user.Phone,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			AdminType: user.AdminType,
			Status:    user.Status,
			CreatedAt: createdAt,
		})
	}

	return &model.AdminUserListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *AdminService) DeleteUser(userID uint) error {
	var user mysql.User
	if err := mysql.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if user.AdminType == 1 {
		zap.S().Warn("尝试删除管理员账号", "targetUserId", userID)
		return errors.New("不能删除管理员账号")
	}

	tx := mysql.DB.Begin()
	defer tx.Rollback()

	if err := tx.Where("user_id = ?", userID).Delete(&mysql.Post{}).Error; err != nil {
		zap.S().Error("删除用户帖子失败", "userId", userID, "error", err)
		return errors.New("删除用户帖子失败")
	}

	if err := tx.Where("user_id = ?", userID).Delete(&mysql.Comment{}).Error; err != nil {
		zap.S().Error("删除用户评论失败", "userId", userID, "error", err)
		return errors.New("删除用户评论失败")
	}

	if err := tx.Where("user_id = ?", userID).Delete(&mysql.PostLike{}).Error; err != nil {
		zap.S().Error("删除用户点赞记录失败", "userId", userID, "error", err)
		return errors.New("删除用户点赞记录失败")
	}

	if err := tx.Where("user_id = ?", userID).Delete(&mysql.PostFavorite{}).Error; err != nil {
		zap.S().Error("删除用户收藏记录失败", "userId", userID, "error", err)
		return errors.New("删除用户收藏记录失败")
	}

	if err := tx.Where("user_id = ?", userID).Delete(&mysql.CommentLike{}).Error; err != nil {
		zap.S().Error("删除用户评论点赞记录失败", "userId", userID, "error", err)
		return errors.New("删除用户评论点赞记录失败")
	}

	if err := tx.Delete(&user).Error; err != nil {
		zap.S().Error("删除用户失败", "userId", userID, "error", err)
		return errors.New("删除用户失败")
	}

	if err := tx.Commit().Error; err != nil {
		zap.S().Error("提交事务失败", "userId", userID, "error", err)
		return errors.New("删除用户失败")
	}

	zap.S().Info("管理员删除用户成功", "userId", userID, "username", user.Username)
	return nil
}

func (s *AdminService) UpdateUserAdminType(userID uint, adminType int) error {
	var user mysql.User
	if err := mysql.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if err := mysql.DB.Model(&user).Update("admin_type", adminType).Error; err != nil {
		return errors.New("更新用户角色失败")
	}

	return nil
}

func (s *AdminService) UpdateUserStatus(userID uint, status int) error {
	var user mysql.User
	if err := mysql.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if err := mysql.DB.Model(&user).Update("status", status).Error; err != nil {
		return errors.New("更新用户状态失败")
	}

	return nil
}

func (s *AdminService) GetPostList(req *model.AdminPostListRequest) (*model.AdminPostListResponse, error) {
	var posts []mysql.Post
	var total int64

	mysql.DB.Model(&mysql.Post{}).Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := mysql.DB.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&posts)

	if result.Error != nil {
		return nil, errors.New("获取帖子列表失败")
	}

	items := make([]model.AdminPostInfo, 0, len(posts))
	for _, post := range posts {
		var user mysql.User
		mysql.DB.Where("id = ?", post.UserID).First(&user)

		createdAt := ""
		if !post.CreatedAt.IsZero() {
			createdAt = post.CreatedAt.Format("2006-01-02 15:04:05")
		}

		items = append(items, model.AdminPostInfo{
			ID:           post.ID,
			UserID:       post.UserID,
			Username:     user.Username,
			Title:        post.Title,
			Summary:      post.Summary,
			Status:       post.Status,
			ViewCount:    post.ViewCount,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			CreatedAt:    createdAt,
		})
	}

	return &model.AdminPostListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *AdminService) DeletePost(postID uint) error {
	var post mysql.Post
	if err := mysql.DB.First(&post, postID).Error; err != nil {
		return errors.New("帖子不存在")
	}

	tx := mysql.DB.Begin()
	defer tx.Rollback()

	if err := tx.Where("post_id = ?", postID).Delete(&mysql.Comment{}).Error; err != nil {
		return errors.New("删除帖子评论失败")
	}

	if err := tx.Where("post_id = ?", postID).Delete(&mysql.PostLike{}).Error; err != nil {
		return errors.New("删除帖子点赞记录失败")
	}

	if err := tx.Where("post_id = ?", postID).Delete(&mysql.PostFavorite{}).Error; err != nil {
		return errors.New("删除帖子收藏记录失败")
	}

	if err := tx.Where("post_id = ?", postID).Delete(&mysql.PostTag{}).Error; err != nil {
		return errors.New("删除帖子标签关联失败")
	}

	if err := tx.Delete(&post).Error; err != nil {
		return errors.New("删除帖子失败")
	}

	return tx.Commit().Error
}
