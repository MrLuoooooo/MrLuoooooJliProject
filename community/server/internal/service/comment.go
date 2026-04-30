package service

import (
	"errors"

	"community-server/DB/mysql"
	"community-server/internal/model"
)

type CommentService struct{}

func NewCommentService() *CommentService {
	return &CommentService{}
}

func (s *CommentService) CreateComment(userID uint, req *model.CreateCommentRequest) (uint, error) {
	var post mysql.Post
	result := mysql.DB.Where("id = ? AND status != 2", req.PostID).First(&post)
	if result.Error != nil {
		return 0, errors.New("帖子不存在")
	}

	if req.ParentID > 0 {
		var parentComment mysql.Comment
		result = mysql.DB.Where("id = ? AND post_id = ?", req.ParentID, req.PostID).First(&parentComment)
		if result.Error != nil {
			return 0, errors.New("父评论不存在")
		}
	}

	comment := mysql.Comment{
		UserID:   userID,
		PostID:   req.PostID,
		ParentID: req.ParentID,
		Content:  req.Content,
		Status:   1,
	}

	result = mysql.DB.Create(&comment)
	if result.Error != nil {
		return 0, errors.New("发表评论失败")
	}

	mysql.DB.Model(&post).UpdateColumn("comment_count", post.CommentCount+1)

	return comment.ID, nil
}

func (s *CommentService) GetCommentList(req *model.CommentListRequest) (*model.CommentListResponse, error) {
	var comments []mysql.Comment
	var total int64

	query := mysql.DB.Model(&mysql.Comment{}).Where("post_id = ? AND parent_id = 0 AND status != 2", req.PostID)

	query.Count(&total)

	offset := (req.Page - 1) * req.PageSize
	if offset < 0 {
		offset = 0
	}

	result := query.Order("created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&comments)

	if result.Error != nil {
		return nil, errors.New("获取评论列表失败")
	}

	items := make([]model.CommentListItem, 0, len(comments))
	for _, comment := range comments {
		var user mysql.User
		mysql.DB.Where("id = ?", comment.UserID).First(&user)

		item := model.CommentListItem{
			ID:        comment.ID,
			UserID:    comment.UserID,
			Username:  user.Username,
			Nickname:  user.Nickname,
			Avatar:    user.Avatar,
			PostID:    comment.PostID,
			ParentID:  comment.ParentID,
			Content:   comment.Content,
			Status:    comment.Status,
			LikeCount: comment.LikeCount,
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		var replies []mysql.Comment
		mysql.DB.Where("post_id = ? AND parent_id = ? AND status != 2", req.PostID, comment.ID).
			Order("created_at ASC").
			Find(&replies)

		if len(replies) > 0 {
			item.Replies = make([]model.CommentListItem, 0, len(replies))
			for _, reply := range replies {
				var replyUser mysql.User
				mysql.DB.Where("id = ?", reply.UserID).First(&replyUser)

				item.Replies = append(item.Replies, model.CommentListItem{
					ID:        reply.ID,
					UserID:    reply.UserID,
					Username:  replyUser.Username,
					Nickname:  replyUser.Nickname,
					Avatar:    replyUser.Avatar,
					PostID:    reply.PostID,
					ParentID:  reply.ParentID,
					Content:   reply.Content,
					Status:    reply.Status,
					LikeCount: reply.LikeCount,
					CreatedAt: reply.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		}

		items = append(items, item)
	}

	return &model.CommentListResponse{
		Total: total,
		Items: items,
	}, nil
}

func (s *CommentService) UpdateComment(userID, commentID uint, req *model.UpdateCommentRequest) error {
	var comment mysql.Comment
	result := mysql.DB.Where("id = ?", commentID).First(&comment)
	if result.Error != nil {
		return errors.New("评论不存在")
	}

	if comment.UserID != userID {
		return errors.New("无权编辑此评论")
	}

	mysql.DB.Model(&comment).Update("content", req.Content)

	return nil
}

func (s *CommentService) DeleteComment(userID, commentID uint) error {
	var comment mysql.Comment
	result := mysql.DB.Where("id = ?", commentID).First(&comment)
	if result.Error != nil {
		return errors.New("评论不存在")
	}

	if comment.UserID != userID {
		return errors.New("无权删除此评论")
	}

	mysql.DB.Model(&comment).Update("status", 2)

	var post mysql.Post
	mysql.DB.Where("id = ?", comment.PostID).First(&post)
	if post.CommentCount > 0 {
		mysql.DB.Model(&post).UpdateColumn("comment_count", post.CommentCount-1)
	}

	var replies []mysql.Comment
	mysql.DB.Where("parent_id = ?", commentID).Find(&replies)
	for _, reply := range replies {
		mysql.DB.Model(&reply).Update("status", 2)
		if post.CommentCount > 0 {
			mysql.DB.Model(&post).UpdateColumn("comment_count", post.CommentCount-1)
		}
	}

	return nil
}

func (s *CommentService) LikeComment(userID, commentID uint) error {
	var comment mysql.Comment
	result := mysql.DB.Where("id = ?", commentID).First(&comment)
	if result.Error != nil {
		return errors.New("评论不存在")
	}

	var existingLike mysql.CommentLike
	result = mysql.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingLike)
	if result.Error == nil {
		return errors.New("已点赞")
	}

	commentLike := mysql.CommentLike{
		UserID:    userID,
		CommentID: commentID,
	}

	result = mysql.DB.Create(&commentLike)
	if result.Error != nil {
		return errors.New("点赞失败")
	}

	mysql.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount+1)

	return nil
}

func (s *CommentService) UnlikeComment(userID, commentID uint) error {
	var commentLike mysql.CommentLike
	result := mysql.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&commentLike)
	if result.Error != nil {
		return errors.New("未点赞")
	}

	mysql.DB.Delete(&commentLike)

	var comment mysql.Comment
	mysql.DB.Where("id = ?", commentID).First(&comment)
	if comment.LikeCount > 0 {
		mysql.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount-1)
	}

	return nil
}
