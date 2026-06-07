package service

import (
	"errors"

	"community-server/internal/db/mysql"
	"community-server/internal/model"
	"community-server/internal/repository"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PostService struct {
	notifSvc     *NotificationService
	postRepo     repository.PostRepository
	postLikeRepo repository.PostLikeRepository
	postFavRepo  repository.PostFavoriteRepository
	commentRepo  repository.CommentRepository
	userRepo     repository.UserRepository
	followRepo   repository.FollowRepository
	categoryRepo repository.CategoryRepository
}

func NewPostService(
	notifSvc *NotificationService,
	postRepo repository.PostRepository,
	postLikeRepo repository.PostLikeRepository,
	postFavRepo repository.PostFavoriteRepository,
	commentRepo repository.CommentRepository,
	userRepo repository.UserRepository,
	followRepo repository.FollowRepository,
	categoryRepo repository.CategoryRepository,
) *PostService {
	return &PostService{
		notifSvc: notifSvc, postRepo: postRepo,
		postLikeRepo: postLikeRepo, postFavRepo: postFavRepo,
		commentRepo: commentRepo, userRepo: userRepo, followRepo: followRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *PostService) CreatePost(userID uint, req *model.CreatePostRequest) (uint, error) {
	post := mysql.Post{
		UserID: userID, Title: req.Title, Content: req.Content,
		Summary: req.Summary, CoverImage: req.CoverImage, CategoryID: req.CategoryID, Status: 1,
	}
	if req.Status != 0 {
		post.Status = req.Status
	}
	if err := s.postRepo.Create(&post); err != nil {
		zap.S().Error("发布帖子失败", "userId", userID, "title", req.Title, "error", err)
		return 0, errors.New("发布帖子失败")
	}
	// 同步分类帖子数
	if post.CategoryID > 0 {
		s.categoryRepo.IncrementPostCount(post.CategoryID, +1)
	}
	zap.S().Info("发布帖子成功", "postId", post.ID, "userId", userID, "title", req.Title)
	return post.ID, nil
}

func (s *PostService) GetPostByID(postID uint) (*mysql.Post, error) {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}
	s.postRepo.UpdateColumn(postID, "view_count", gorm.Expr("view_count + 1"))
	return post, nil
}

func (s *PostService) GetPost(postID uint) (*model.PostDetailResponse, error) {
	post, err := s.GetPostByID(postID)
	if err != nil {
		return nil, err
	}
	user, _ := s.userRepo.FindByID(post.UserID)
	username, nickname, avatar := "", "", ""
	if user != nil {
		username, nickname, avatar = user.Username, user.Nickname, user.Avatar
	}
	return &model.PostDetailResponse{
		ID: post.ID, UserID: post.UserID, Username: username, Nickname: nickname, Avatar: avatar,
		Title: post.Title, Content: post.Content, Summary: post.Summary, CoverImage: post.CoverImage,
		CategoryID: post.CategoryID, Status: post.Status, ViewCount: post.ViewCount,
		LikeCount: post.LikeCount, CommentCount: post.CommentCount,
		IsTop: post.IsTop, IsEssence: post.IsEssence,
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *PostService) GetPostList(req *model.PostListRequest) (*model.PostListResponse, error) {
	q := repository.PostListQuery{
		CategoryID: req.CategoryID, Keyword: req.Keyword, Sort: req.Sort,
		Page: req.Page, PageSize: req.PageSize, StatusExclude: 2,
	}
	posts, total, err := s.postRepo.List(q)
	if err != nil {
		return nil, errors.New("获取帖子列表失败")
	}
	return s.buildPostList(posts, total), nil
}

func (s *PostService) UpdatePost(userID, postID uint, req *model.UpdatePostRequest) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
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
		return s.postRepo.Update(postID, updates)
	}
	return nil
}

func (s *PostService) DeletePost(userID, postID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.UserID != userID {
		return errors.New("无权删除此帖子")
	}
	// 事务在 Repository 层处理，这里直接顺序调用
	if err := s.postRepo.SoftDelete(postID); err != nil {
		zap.S().Error("软删除帖子失败", "postId", postID, "error", err)
		return errors.New("删除帖子失败")
	}
	s.commentRepo.SoftDeleteByPostID(postID)
	// 同步分类帖子数
	if post.CategoryID > 0 {
		s.categoryRepo.IncrementPostCount(post.CategoryID, -1)
	}
	zap.S().Info("删除帖子成功", "postId", postID, "userId", userID)
	return nil
}

func (s *PostService) LikePost(userID, postID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	exists, _ := s.postLikeRepo.Exists(userID, postID)
	if exists {
		return errors.New("已点赞")
	}
	if err := s.postLikeRepo.Create(&mysql.PostLike{UserID: userID, PostID: postID}); err != nil {
		return errors.New("点赞失败")
	}
	s.postRepo.UpdateColumn(postID, "like_count", gorm.Expr("like_count + 1"))
	if post.UserID != userID && s.notifSvc != nil {
		s.notifSvc.CreateAndPush(post.UserID, userID, model.NotifyLike, postID, "赞了你的帖子")
	}
	return nil
}

func (s *PostService) UnlikePost(userID, postID uint) error {
	exists, _ := s.postLikeRepo.Exists(userID, postID)
	if !exists {
		return errors.New("未点赞")
	}
	s.postLikeRepo.Delete(userID, postID)
	s.postRepo.UpdateColumn(postID, "like_count", gorm.Expr("GREATEST(like_count - 1, 0)"))
	return nil
}

func (s *PostService) FavoritePost(userID, postID uint) error {
	if _, err := s.postRepo.FindByID(postID); err != nil {
		return errors.New("帖子不存在")
	}
	exists, _ := s.postFavRepo.Exists(userID, postID)
	if exists {
		return errors.New("已收藏")
	}
	return s.postFavRepo.Create(&mysql.PostFavorite{UserID: userID, PostID: postID})
}

func (s *PostService) UnfavoritePost(userID, postID uint) error {
	exists, _ := s.postFavRepo.Exists(userID, postID)
	if !exists {
		return errors.New("未收藏")
	}
	return s.postFavRepo.Delete(userID, postID)
}

func (s *PostService) GetUserPosts(userID uint, req *model.PostListRequest) (*model.PostListResponse, error) {
	q := repository.PostListQuery{
		UserID: userID, Page: req.Page, PageSize: req.PageSize, StatusExclude: 2,
	}
	posts, total, err := s.postRepo.List(q)
	if err != nil {
		return nil, errors.New("获取帖子列表失败")
	}
	return s.buildPostList(posts, total), nil
}

func (s *PostService) GetUserFavorites(userID uint, page, pageSize int) (*model.PostListResponse, error) {
	favs, total, err := s.postFavRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, errors.New("获取收藏列表失败")
	}
	if len(favs) == 0 {
		return &model.PostListResponse{Total: total, Items: []model.PostListItem{}}, nil
	}
	postIDs := make([]uint, len(favs))
	for i, f := range favs {
		postIDs[i] = f.PostID
	}
	posts, _, err := s.postRepo.List(repository.PostListQuery{
		Page: 1, PageSize: len(postIDs), StatusExclude: 2,
	})
	// filter to only matched posts
	postMap := make(map[uint]mysql.Post)
	for _, p := range posts {
		postMap[p.ID] = p
	}
	return s.buildPostListFromIDs(postIDs, postMap, total), nil
}

func (s *PostService) GetUserLikedPosts(userID uint, page, pageSize int) (*model.PostListResponse, error) {
	likes, total, err := s.postLikeRepo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, errors.New("获取点赞列表失败")
	}
	if len(likes) == 0 {
		return &model.PostListResponse{Total: total, Items: []model.PostListItem{}}, nil
	}
	postIDs := make([]uint, len(likes))
	for i, l := range likes {
		postIDs[i] = l.PostID
	}
	posts, _, err := s.postRepo.List(repository.PostListQuery{
		Page: 1, PageSize: len(postIDs), StatusExclude: 2,
	})
	postMap := make(map[uint]mysql.Post)
	for _, p := range posts {
		postMap[p.ID] = p
	}
	return s.buildPostListFromIDs(postIDs, postMap, total), nil
}

func (s *PostService) GetFollowFeed(userID uint, page, pageSize int) (*model.PostListResponse, error) {
	follows, err := s.followRepo.FindAllFollowing(userID)
	if err != nil || len(follows) == 0 {
		return &model.PostListResponse{Total: 0, Items: []model.PostListItem{}}, nil
	}
	followIDs := make([]uint, len(follows))
	for i, f := range follows {
		followIDs[i] = f.FollowID
	}
	q := repository.PostListQuery{
		UserIDs: followIDs, Page: page, PageSize: pageSize, StatusExclude: 2, Sort: "newest",
	}
	posts, total, err := s.postRepo.List(q)
	if err != nil {
		return nil, errors.New("获取动态失败")
	}
	return s.buildPostList(posts, total), nil
}

// --- helpers ---

func (s *PostService) buildPostList(posts []mysql.Post, total int64) *model.PostListResponse {
	items := make([]model.PostListItem, 0, len(posts))
	if len(posts) == 0 {
		return &model.PostListResponse{Total: total, Items: items}
	}
	userIDs := make([]uint, len(posts))
	for i, p := range posts {
		userIDs[i] = p.UserID
	}
	users, _ := s.userRepo.FindByIDs(userIDs)
	userMap := make(map[uint]mysql.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}
	for _, p := range posts {
		u := userMap[p.UserID]
		items = append(items, model.PostListItem{
			ID: p.ID, UserID: p.UserID, Username: u.Username, Nickname: u.Nickname,
			Title: p.Title, Summary: p.Summary, CoverImage: p.CoverImage,
			CategoryID: p.CategoryID, Status: p.Status,
			ViewCount: p.ViewCount, LikeCount: p.LikeCount, CommentCount: p.CommentCount,
			IsTop: p.IsTop, IsEssence: p.IsEssence,
			CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &model.PostListResponse{Total: total, Items: items}
}

func (s *PostService) buildPostListFromIDs(postIDs []uint, postMap map[uint]mysql.Post, total int64) *model.PostListResponse {
	items := make([]model.PostListItem, 0, len(postIDs))
	if len(postIDs) == 0 {
		return &model.PostListResponse{Total: total, Items: items}
	}
	// collect all user IDs from matched posts
	userIDs := make([]uint, 0, len(postMap))
	for _, p := range postMap {
		userIDs = append(userIDs, p.UserID)
	}
	users, _ := s.userRepo.FindByIDs(userIDs)
	userMap := make(map[uint]mysql.User, len(users))
	for _, u := range users {
		userMap[u.ID] = u
	}
	for _, pid := range postIDs {
		p, ok := postMap[pid]
		if !ok {
			continue
		}
		u := userMap[p.UserID]
		items = append(items, model.PostListItem{
			ID: p.ID, UserID: p.UserID, Username: u.Username, Nickname: u.Nickname,
			Title: p.Title, Summary: p.Summary, CoverImage: p.CoverImage,
			CategoryID: p.CategoryID, Status: p.Status,
			ViewCount: p.ViewCount, LikeCount: p.LikeCount, CommentCount: p.CommentCount,
			IsTop: p.IsTop, IsEssence: p.IsEssence,
			CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return &model.PostListResponse{Total: total, Items: items}
}
