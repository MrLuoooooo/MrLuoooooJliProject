package service

import (
	"errors"

	"community-server/DB/mysql"
	"community-server/internal/model"

	"go.uber.org/zap"
)

type PostService struct{}

func NewPostService() *PostService {
	return &PostService{}
}

func (s *PostService) CreatePost(userID uint, req *model.CreatePostRequest) (uint, error) {
	post := mysql.Post{
		UserID:     userID,
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		Status:     1,
	}

	if req.Status != 0 {
		post.Status = req.Status
	}

	result := mysql.DB.Create(&post)
	if result.Error != nil {
		zap.S().Error("发布帖子失败", "userId", userID, "title", req.Title, "error", result.Error)
		return 0, errors.New("发布帖子失败")
	}

	zap.S().Info("发布帖子成功", "postId", post.ID, "userId", userID, "title", req.Title)
	return post.ID, nil
}

func (s *PostService) GetPostByID(postID uint) (*mysql.Post, error) {
	var post mysql.Post
	result := mysql.DB.Where("id = ? AND status != 2", postID).First(&post)
	if result.Error != nil {
		return nil, errors.New("帖子不存在")
	}

	mysql.DB.Model(&post).UpdateColumn("view_count", post.ViewCount+1)

	return &post, nil
}

func (s *PostService) GetPostList(req *model.PostListRequest) (*model.PostListResponse, error) {
	var posts []mysql.Post
	var total int64

	query := mysql.DB.Model(&mysql.Post{}).Where("status != 2")

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
		return nil, errors.New("获取帖子列表失败")
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

	return &model.PostListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *PostService) UpdatePost(userID, postID uint, req *model.UpdatePostRequest) error {
	var post mysql.Post
	result := mysql.DB.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return errors.New("帖子不存在")
	}

	if post.UserID != userID {
		return errors.New("无权编辑此帖子")
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Summary != "" {
		updates["summary"] = req.Summary
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Status > 0 {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		mysql.DB.Model(&post).Updates(updates)
	}

	return nil
}

func (s *PostService) DeletePost(userID, postID uint) error {
	var post mysql.Post
	result := mysql.DB.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return errors.New("帖子不存在")
	}

	if post.UserID != userID {
		return errors.New("无权删除此帖子")
	}

	mysql.DB.Model(&post).Update("status", 2)

	zap.S().Info("删除帖子成功", "postId", postID, "userId", userID)
	return nil
}

func (s *PostService) LikePost(userID, postID uint) error {
	var post mysql.Post
	result := mysql.DB.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return errors.New("帖子不存在")
	}

	var existingLike mysql.PostLike
	result = mysql.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingLike)
	if result.Error == nil {
		return errors.New("已点赞")
	}

	postLike := mysql.PostLike{
		UserID: userID,
		PostID: postID,
	}

	result = mysql.DB.Create(&postLike)
	if result.Error != nil {
		return errors.New("点赞失败")
	}

	mysql.DB.Model(&post).UpdateColumn("like_count", post.LikeCount+1)

	return nil
}

func (s *PostService) UnlikePost(userID, postID uint) error {
	var postLike mysql.PostLike
	result := mysql.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&postLike)
	if result.Error != nil {
		return errors.New("未点赞")
	}

	mysql.DB.Delete(&postLike)

	var post mysql.Post
	mysql.DB.Where("id = ?", postID).First(&post)
	if post.LikeCount > 0 {
		mysql.DB.Model(&post).UpdateColumn("like_count", post.LikeCount-1)
	}

	return nil
}

func (s *PostService) FavoritePost(userID, postID uint) error {
	var post mysql.Post
	result := mysql.DB.Where("id = ?", postID).First(&post)
	if result.Error != nil {
		return errors.New("帖子不存在")
	}

	var existingFavorite mysql.PostFavorite
	result = mysql.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&existingFavorite)
	if result.Error == nil {
		return errors.New("已收藏")
	}

	postFavorite := mysql.PostFavorite{
		UserID: userID,
		PostID: postID,
	}

	result = mysql.DB.Create(&postFavorite)
	if result.Error != nil {
		return errors.New("收藏失败")
	}

	return nil
}

func (s *PostService) UnfavoritePost(userID, postID uint) error {
	var postFavorite mysql.PostFavorite
	result := mysql.DB.Where("user_id = ? AND post_id = ?", userID, postID).First(&postFavorite)
	if result.Error != nil {
		return errors.New("未收藏")
	}

	mysql.DB.Delete(&postFavorite)

	return nil
}

func (s *PostService) GetUserPosts(userID uint, req *model.PostListRequest) (*model.PostListResponse, error) {
	var posts []mysql.Post
	var total int64

	query := mysql.DB.Model(&mysql.Post{}).Where("user_id = ? AND status != 2", userID)

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&posts)

	if result.Error != nil {
		return nil, errors.New("获取帖子列表失败")
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

	return &model.PostListResponse{
		Total: total,
		Items: items,
	}, nil
}
