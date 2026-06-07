package service

import (
	"errors"

	"community-server/internal/db/mysql"
	"community-server/internal/model"
	"community-server/internal/repository"
	"community-server/internal/ws"

	"go.uber.org/zap"
)

type AdminService struct {
	userRepo    repository.UserRepository
	postRepo    repository.PostRepository
	postTagRepo repository.PostTagRepository
	commentRepo repository.CommentRepository
	wsManager   *ws.Manager
}

func NewAdminService(userRepo repository.UserRepository, postRepo repository.PostRepository, postTagRepo repository.PostTagRepository, commentRepo repository.CommentRepository, wsManager *ws.Manager) *AdminService {
	return &AdminService{userRepo: userRepo, postRepo: postRepo, postTagRepo: postTagRepo, commentRepo: commentRepo, wsManager: wsManager}
}

func (s *AdminService) GetUserList(req *model.AdminUserListRequest) (*model.AdminUserListResponse, error) {
	users, total, err := s.userRepo.List(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("获取用户列表失败")
	}
	items := make([]model.AdminUserInfo, 0, len(users))
	for _, user := range users {
		createdAt := ""
		if !user.CreatedAt.IsZero() {
			createdAt = user.CreatedAt.Format("2006-01-02 15:04:05")
		}
		items = append(items, model.AdminUserInfo{
			ID: user.ID, Username: user.Username, Email: user.Email, Phone: user.Phone,
			Nickname: user.Nickname, Avatar: user.Avatar, AdminType: user.AdminType,
			Status: user.Status, CreatedAt: createdAt,
		})
	}
	return &model.AdminUserListResponse{Total: total, Items: items}, nil
}

func (s *AdminService) DeleteUser(userID uint) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if user.AdminType == 1 {
		zap.S().Warn("尝试删除管理员账号", "targetUserId", userID)
		return errors.New("不能删除管理员账号")
	}
	return s.userRepo.Delete(userID)
}

func (s *AdminService) UpdateUserAdminType(userID uint, adminType int) error {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return errors.New("用户不存在")
	}
	return s.userRepo.Update(userID, map[string]interface{}{"admin_type": adminType})
}

func (s *AdminService) UpdateUserStatus(userID uint, status int) error {
	if _, err := s.userRepo.FindByID(userID); err != nil {
		return errors.New("用户不存在")
	}
	return s.userRepo.Update(userID, map[string]interface{}{"status": status})
}

func (s *AdminService) GetPostList(req *model.AdminPostListRequest) (*model.AdminPostListResponse, error) {
	posts, total, err := s.postRepo.List(repository.PostListQuery{Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, errors.New("获取帖子列表失败")
	}
	items := make([]model.AdminPostInfo, 0, len(posts))
	if len(posts) > 0 {
		userIDs := make([]uint, len(posts))
		for i, p := range posts {
			userIDs[i] = p.UserID
		}
		users, _ := s.userRepo.FindByIDs(userIDs)
		userMap := make(map[uint]mysql.User)
		for _, u := range users {
			userMap[u.ID] = u
		}
		for _, post := range posts {
			u := userMap[post.UserID]
			createdAt := ""
			if !post.CreatedAt.IsZero() {
				createdAt = post.CreatedAt.Format("2006-01-02 15:04:05")
			}
			items = append(items, model.AdminPostInfo{
				ID: post.ID, UserID: post.UserID, Username: u.Username,
				Title: post.Title, Summary: post.Summary, Status: post.Status,
				ViewCount: post.ViewCount, LikeCount: post.LikeCount,
				CommentCount: post.CommentCount, CreatedAt: createdAt,
			})
		}
	}
	return &model.AdminPostListResponse{Total: total, Items: items}, nil
}

func (s *AdminService) DeletePost(postID uint) error {
	if _, err := s.postRepo.FindByID(postID); err != nil {
		return errors.New("帖子不存在")
	}
	s.postTagRepo.DeleteByPostID(postID)
	return s.postRepo.Delete(postID)
}

func (s *AdminService) SetPostTop(postID uint, isTop bool) error {
	if _, err := s.postRepo.FindByID(postID); err != nil {
		return errors.New("帖子不存在")
	}
	return s.postRepo.Update(postID, map[string]interface{}{"is_top": isTop})
}

func (s *AdminService) SetPostEssence(postID uint, isEssence bool) error {
	if _, err := s.postRepo.FindByID(postID); err != nil {
		return errors.New("帖子不存在")
	}
	return s.postRepo.Update(postID, map[string]interface{}{"is_essence": isEssence})
}

func (s *AdminService) GetStats() (*model.AdminStatsResponse, error) {
	totalUsers, _ := s.userRepo.Count()
	newUsersToday, _ := s.userRepo.CountToday()
	totalPosts, _ := s.postRepo.Count()
	postsToday, _ := s.postRepo.CountToday()
	totalComments, _ := s.commentRepo.Count()
	onlineCount := s.wsManager.ConnCount()
	return &model.AdminStatsResponse{
		TotalUsers:    totalUsers,
		NewUsersToday: newUsersToday,
		TotalPosts:    totalPosts,
		PostsToday:    postsToday,
		TotalComments: totalComments,
		OnlineCount:   onlineCount,
	}, nil
}
