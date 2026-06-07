package service

import (
	"errors"
	"fmt"

	"community-server/internal/db/mysql"
	"community-server/internal/im"
	"community-server/internal/model"
	"community-server/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CommentService struct {
	imClient         im.IMClient
	notifSvc         *NotificationService
	commentRepo      repository.CommentRepository
	commentLikeRepo  repository.CommentLikeRepository
	postRepo         repository.PostRepository
	userRepo         repository.UserRepository
}

func NewCommentService(
	imClient im.IMClient,
	notifSvc *NotificationService,
	commentRepo repository.CommentRepository,
	commentLikeRepo repository.CommentLikeRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
) *CommentService {
	return &CommentService{
		imClient: imClient, notifSvc: notifSvc,
		commentRepo: commentRepo, commentLikeRepo: commentLikeRepo, postRepo: postRepo, userRepo: userRepo,
	}
}

func (s *CommentService) CreateComment(userID uint, req *model.CreateCommentRequest) (uint, error) {
	post, err := s.postRepo.FindByID(req.PostID)
	if err != nil {
		return 0, errors.New("帖子不存在")
	}
	if req.ParentID > 0 {
		if _, err := s.commentRepo.FindByID(req.ParentID); err != nil {
			return 0, errors.New("父评论不存在")
		}
	}
	comment := mysql.Comment{
		UserID: userID, PostID: req.PostID, ParentID: req.ParentID, Content: req.Content, Status: 1,
	}
	if err := s.commentRepo.Create(&comment); err != nil {
		return 0, errors.New("发表评论失败")
	}
	s.postRepo.UpdateColumn(req.PostID, "comment_count", gorm.Expr("comment_count + 1"))

	if post.UserID != userID {
		content := fmt.Sprintf("评论了你的帖子《%s》", post.Title)
		s.notifSvc.CreateAndPush(post.UserID, userID, model.NotifyComment, post.ID, content)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					zap.S().Error("IM通知帖子作者 panic", "postId", post.ID, "panic", r)
				}
			}()
			body := fmt.Sprintf("你的帖子《%s》收到了一条新评论", post.Title)
			if err := s.imClient.SendSystemMsg(im.UserIDToStr(post.UserID), im.UserIDToStr(post.UserID), body); err != nil {
				zap.S().Warn("IM通知帖子作者失败", "postId", post.ID, "error", err)
			}
		}()
	}
	return comment.ID, nil
}

func (s *CommentService) GetCommentList(req *model.CommentListRequest) (*model.CommentListResponse, error) {
	comments, total, err := s.commentRepo.FindRootByPostID(req.PostID, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("获取评论列表失败")
	}
	items := make([]model.CommentListItem, 0, len(comments))
	if len(comments) == 0 {
		return &model.CommentListResponse{Total: total, Items: items}, nil
	}

	// batch load comment users
	commentUserIDs := make([]uint, len(comments))
	for i, c := range comments {
		commentUserIDs[i] = c.UserID
	}
	allUsers, _ := s.userRepo.FindByIDs(commentUserIDs)
	userMap := make(map[uint]mysql.User, len(allUsers))
	for _, u := range allUsers {
		userMap[u.ID] = u
	}

	// batch load replies
	commentIDs := make([]uint, len(comments))
	for i, c := range comments {
		commentIDs[i] = c.ID
	}
	allReplies, _ := s.commentRepo.FindRepliesByPostAndParents(req.PostID, commentIDs)

	// batch load reply users
	replyUserIDs := make([]uint, 0, len(allReplies))
	for _, r := range allReplies {
		replyUserIDs = append(replyUserIDs, r.UserID)
	}
	var replyUsers []mysql.User
	if len(replyUserIDs) > 0 {
		replyUsers, _ = s.userRepo.FindByIDs(replyUserIDs)
	}
	replyUserMap := make(map[uint]mysql.User, len(replyUsers))
	for _, u := range replyUsers {
		replyUserMap[u.ID] = u
	}

	repliesByParent := make(map[uint][]mysql.Comment)
	for _, r := range allReplies {
		repliesByParent[r.ParentID] = append(repliesByParent[r.ParentID], r)
	}

	for _, comment := range comments {
		u := userMap[comment.UserID]
		item := model.CommentListItem{
			ID: comment.ID, UserID: comment.UserID,
			Username: u.Username, Nickname: u.Nickname, Avatar: u.Avatar,
			PostID: comment.PostID, ParentID: comment.ParentID,
			Content: comment.Content, Status: comment.Status, LikeCount: comment.LikeCount,
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if replies, ok := repliesByParent[comment.ID]; ok {
			item.Replies = make([]model.CommentListItem, 0, len(replies))
			for _, reply := range replies {
				ru := replyUserMap[reply.UserID]
				item.Replies = append(item.Replies, model.CommentListItem{
					ID: reply.ID, UserID: reply.UserID,
					Username: ru.Username, Nickname: ru.Nickname, Avatar: ru.Avatar,
					PostID: reply.PostID, ParentID: reply.ParentID,
					Content: reply.Content, Status: reply.Status, LikeCount: reply.LikeCount,
					CreatedAt: reply.CreatedAt.Format("2006-01-02 15:04:05"),
				})
			}
		}
		items = append(items, item)
	}
	return &model.CommentListResponse{Total: total, Items: items}, nil
}

func (s *CommentService) UpdateComment(userID, commentID uint, req *model.UpdateCommentRequest) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	if comment.UserID != userID {
		return errors.New("无权编辑此评论")
	}
	return s.commentRepo.Update(commentID, map[string]interface{}{"content": req.Content})
}

func (s *CommentService) DeleteComment(userID, commentID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	if comment.UserID != userID {
		return errors.New("无权删除此评论")
	}
	return s.commentRepo.SoftDelete(commentID)
}

func (s *CommentService) LikeComment(userID, commentID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	exists, _ := s.commentLikeRepo.Exists(userID, commentID)
	if exists {
		return errors.New("已点赞")
	}
	if err := s.commentLikeRepo.Create(&mysql.CommentLike{UserID: userID, CommentID: commentID}); err != nil {
		return errors.New("点赞失败")
	}
	s.commentRepo.UpdateColumn(commentID, "like_count", gorm.Expr("like_count + 1"))
	if comment.UserID != userID && s.notifSvc != nil {
		s.notifSvc.CreateAndPush(comment.UserID, userID, model.NotifyLike, commentID, "赞了你的评论")
	}
	return nil
}

func (s *CommentService) UnlikeComment(userID, commentID uint) error {
	exists, _ := s.commentLikeRepo.Exists(userID, commentID)
	if !exists {
		return errors.New("未点赞")
	}
	s.commentLikeRepo.Delete(userID, commentID)
	s.commentRepo.UpdateColumn(commentID, "like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	return nil
}
