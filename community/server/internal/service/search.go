package service

import (
	"errors"

	"community-server/internal/db/mysql"
	"community-server/internal/model"
	"community-server/internal/repository"
)

type SearchService struct {
	postRepo repository.PostRepository
	userRepo repository.UserRepository
}

func NewSearchService(postRepo repository.PostRepository, userRepo repository.UserRepository) *SearchService {
	return &SearchService{postRepo: postRepo, userRepo: userRepo}
}

func (s *SearchService) SearchPosts(req *model.SearchRequest) (*model.SearchResponse, error) {
	posts, total, err := s.postRepo.Search(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("搜索失败")
	}
	items := make([]model.PostListItem, 0, len(posts))
	if len(posts) > 0 {
		userIDs := make([]uint, len(posts))
		for i, p := range posts {
			userIDs[i] = p.UserID
		}
		users, _ := s.userRepo.FindByIDs(userIDs)
		userMap := make(map[uint]mysql.User, len(users))
		for _, u := range users {
			userMap[u.ID] = u
		}
		for _, post := range posts {
			u := userMap[post.UserID]
			items = append(items, model.PostListItem{
				ID: post.ID, UserID: post.UserID, Username: u.Username, Nickname: u.Nickname,
				Title: post.Title, Summary: post.Summary, CoverImage: post.CoverImage,
				CategoryID: post.CategoryID, Status: post.Status,
				ViewCount: post.ViewCount, LikeCount: post.LikeCount, CommentCount: post.CommentCount,
				IsTop: post.IsTop, IsEssence: post.IsEssence,
				CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}
	return &model.SearchResponse{Total: total, Type: "post", Posts: items}, nil
}

func (s *SearchService) SearchUsers(req *model.SearchRequest) (*model.SearchResponse, error) {
	users, total, err := s.userRepo.Search(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("搜索失败")
	}
	items := make([]model.SearchUser, 0, len(users))
	for _, u := range users {
		items = append(items, model.SearchUser{
			ID: u.ID, Username: u.Username, Nickname: u.Nickname, Avatar: u.Avatar, Bio: u.Bio,
		})
	}
	return &model.SearchResponse{Total: total, Type: "user", Users: items}, nil
}
