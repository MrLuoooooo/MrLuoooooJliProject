package service

import (
	"errors"

	"community-server/DB/mysql"
	"community-server/internal/model"
)

type SearchService struct{}

func NewSearchService() *SearchService {
	return &SearchService{}
}

func (s *SearchService) SearchPosts(req *model.SearchRequest) (*model.SearchResponse, error) {
	var posts []mysql.Post
	var total int64

	query := mysql.DB.Model(&mysql.Post{}).
		Where("status != 2 AND (title LIKE ? OR content LIKE ? OR summary LIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := query.Order("is_top DESC, created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&posts)

	if result.Error != nil {
		return nil, errors.New("搜索失败")
	}

	items := make([]model.PostListItem, 0, len(posts))
	for _, post := range posts {
		var user mysql.User
		mysql.DB.Where("id = ?", post.UserID).First(&user)

		items = append(items, model.PostListItem{
			ID:           post.ID,
			UserID:       post.UserID,
			Username:     user.Username,
			Nickname:     user.Nickname,
			Title:        post.Title,
			Summary:      post.Summary,
			CoverImage:   post.CoverImage,
			CategoryID:   post.CategoryID,
			Status:       post.Status,
			ViewCount:    post.ViewCount,
			LikeCount:    post.LikeCount,
			CommentCount: post.CommentCount,
			IsTop:        post.IsTop,
			IsEssence:    post.IsEssence,
			CreatedAt:    post.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &model.SearchResponse{
		Total: total,
		Type:  "post",
		Posts: items,
	}, nil
}

func (s *SearchService) SearchUsers(req *model.SearchRequest) (*model.SearchResponse, error) {
	var users []mysql.User
	var total int64

	query := mysql.DB.Model(&mysql.User{}).
		Where("status = 1 AND (username LIKE ? OR nickname LIKE ? OR bio LIKE ?)",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&users)

	if result.Error != nil {
		return nil, errors.New("搜索失败")
	}

	items := make([]model.SearchUser, 0, len(users))
	for _, user := range users {
		items = append(items, model.SearchUser{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
			Bio:      user.Bio,
		})
	}

	return &model.SearchResponse{
		Total: total,
		Type:  "user",
		Users: items,
	}, nil
}
