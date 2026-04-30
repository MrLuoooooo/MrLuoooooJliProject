package service

import (
	"errors"

	"community-server/DB/mysql"
	"community-server/internal/model"
)

type TagService struct{}

func NewTagService() *TagService {
	return &TagService{}
}

func (s *TagService) CreateTag(req *model.CreateTagRequest) (uint, error) {
	var existingTag mysql.Tag
	result := mysql.DB.Where("name = ?", req.Name).First(&existingTag)
	if result.Error == nil {
		return 0, errors.New("标签已存在")
	}

	tag := mysql.Tag{
		Name:        req.Name,
		Description: req.Description,
		PostCount:   0,
		Status:      1,
	}

	result = mysql.DB.Create(&tag)
	if result.Error != nil {
		return 0, errors.New("创建标签失败")
	}

	return tag.ID, nil
}

func (s *TagService) GetTagList(req *model.TagListRequest) (*model.TagListResponse, error) {
	var tags []mysql.Tag
	var total int64

	query := mysql.DB.Model(&mysql.Tag{}).Where("status != 2")

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := query.Order("post_count DESC, created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&tags)

	if result.Error != nil {
		return nil, errors.New("获取标签列表失败")
	}

	items := make([]model.TagResponse, 0, len(tags))
	for _, tag := range tags {
		items = append(items, model.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			PostCount:   tag.PostCount,
			Status:      tag.Status,
			CreatedAt:   tag.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &model.TagListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *TagService) UpdateTag(tagID uint, req *model.UpdateTagRequest) error {
	var tag mysql.Tag
	result := mysql.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		return errors.New("标签不存在")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) > 0 {
		mysql.DB.Model(&tag).Updates(updates)
	}

	return nil
}

func (s *TagService) DeleteTag(tagID uint) error {
	var tag mysql.Tag
	result := mysql.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		return errors.New("标签不存在")
	}

	mysql.DB.Model(&tag).Update("status", 2)

	mysql.DB.Where("tag_id = ?", tagID).Delete(&mysql.PostTag{})

	return nil
}

func (s *TagService) AddPostTags(postID uint, tagIDs []uint) error {
	var post mysql.Post
	result := mysql.DB.Where("id = ? AND status != 2", postID).First(&post)
	if result.Error != nil {
		return errors.New("帖子不存在")
	}

	for _, tagID := range tagIDs {
		var tag mysql.Tag
		result = mysql.DB.Where("id = ? AND status != 2", tagID).First(&tag)
		if result.Error != nil {
			continue
		}

		var existingPostTag mysql.PostTag
		result = mysql.DB.Where("post_id = ? AND tag_id = ?", postID, tagID).First(&existingPostTag)
		if result.Error == nil {
			continue
		}

		postTag := mysql.PostTag{
			PostID: postID,
			TagID:  tagID,
		}

		mysql.DB.Create(&postTag)
		mysql.DB.Model(&tag).UpdateColumn("post_count", tag.PostCount+1)
	}

	return nil
}

func (s *TagService) RemovePostTag(postID, tagID uint) error {
	var postTag mysql.PostTag
	result := mysql.DB.Where("post_id = ? AND tag_id = ?", postID, tagID).First(&postTag)
	if result.Error != nil {
		return errors.New("标签关联不存在")
	}

	mysql.DB.Delete(&postTag)

	var tag mysql.Tag
	mysql.DB.Where("id = ?", tagID).First(&tag)
	if tag.PostCount > 0 {
		mysql.DB.Model(&tag).UpdateColumn("post_count", tag.PostCount-1)
	}

	return nil
}

func (s *TagService) GetPostTags(postID uint) ([]model.TagResponse, error) {
	var postTags []mysql.PostTag
	result := mysql.DB.Where("post_id = ?", postID).Find(&postTags)
	if result.Error != nil {
		return nil, errors.New("获取标签失败")
	}

	tags := make([]model.TagResponse, 0, len(postTags))
	for _, postTag := range postTags {
		var tag mysql.Tag
		mysql.DB.Where("id = ?", postTag.TagID).First(&tag)

		tags = append(tags, model.TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			PostCount:   tag.PostCount,
			Status:      tag.Status,
			CreatedAt:   tag.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return tags, nil
}

func (s *TagService) GetPostsByTag(tagID uint, page, pageSize int) (*model.PostListResponse, error) {
	var tag mysql.Tag
	result := mysql.DB.Where("id = ? AND status != 2", tagID).First(&tag)
	if result.Error != nil {
		return nil, errors.New("标签不存在")
	}

	var postTags []mysql.PostTag
	var total int64

	query := mysql.DB.Model(&mysql.PostTag{}).Where("tag_id = ?", tagID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	result = query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&postTags)

	if result.Error != nil {
		return nil, errors.New("获取帖子列表失败")
	}

	items := make([]model.PostListItem, 0, len(postTags))
	for _, postTag := range postTags {
		var post mysql.Post
		mysql.DB.Where("id = ? AND status != 2", postTag.PostID).First(&post)

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
