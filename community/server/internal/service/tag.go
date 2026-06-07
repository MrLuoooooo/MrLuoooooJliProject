package service

import (
	"errors"

	"community-server/internal/db/mysql"
	"community-server/internal/model"
	"community-server/internal/repository"
)

type TagService struct {
	tagRepo     repository.TagRepository
	postTagRepo repository.PostTagRepository
	postRepo    repository.PostRepository
	userRepo    repository.UserRepository
}

func NewTagService(
	tagRepo repository.TagRepository,
	postTagRepo repository.PostTagRepository,
	postRepo repository.PostRepository,
	userRepo repository.UserRepository,
) *TagService {
	return &TagService{tagRepo: tagRepo, postTagRepo: postTagRepo, postRepo: postRepo, userRepo: userRepo}
}

func (s *TagService) CreateTag(req *model.CreateTagRequest) (uint, error) {
	if _, err := s.tagRepo.FindByName(req.Name); err == nil {
		return 0, errors.New("标签已存在")
	}
	tag := mysql.Tag{Name: req.Name, Description: req.Description, Status: 1}
	if err := s.tagRepo.Create(&tag); err != nil {
		return 0, errors.New("创建标签失败")
	}
	return tag.ID, nil
}

func (s *TagService) GetTagList(req *model.TagListRequest) (*model.TagListResponse, error) {
	tags, total, err := s.tagRepo.List(req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("获取标签列表失败")
	}
	items := make([]model.TagResponse, 0, len(tags))
	for _, t := range tags {
		items = append(items, model.TagResponse{
			ID: t.ID, Name: t.Name, Description: t.Description, PostCount: t.PostCount,
		})
	}
	return &model.TagListResponse{Total: total, Items: items}, nil
}

func (s *TagService) UpdateTag(tagID uint, req *model.UpdateTagRequest) error {
	if _, err := s.tagRepo.FindByID(tagID); err != nil {
		return errors.New("标签不存在")
	}
	return s.tagRepo.Update(tagID, map[string]interface{}{"name": req.Name, "description": req.Description})
}

func (s *TagService) DeleteTag(tagID uint) error {
	if _, err := s.tagRepo.FindByID(tagID); err != nil {
		return errors.New("标签不存在")
	}
	return s.tagRepo.Delete(tagID)
}

func (s *TagService) AddPostTags(postID uint, tagIDs []uint) error {
	for _, tid := range tagIDs {
		if err := s.postTagRepo.Create(&mysql.PostTag{PostID: postID, TagID: tid}); err != nil {
			return err
		}
		s.tagRepo.IncrementPostCount(tid, 1)
	}
	return nil
}

func (s *TagService) RemovePostTag(postID, tagID uint) error {
	if err := s.postTagRepo.Delete(postID, tagID); err != nil {
		return err
	}
	s.tagRepo.IncrementPostCount(tagID, -1)
	return nil
}

func (s *TagService) GetPostTags(postID uint) ([]model.TagResponse, error) {
	pts, err := s.postTagRepo.FindByPostID(postID)
	if err != nil {
		return nil, err
	}
	tagIDs := make([]uint, len(pts))
	for i, pt := range pts {
		tagIDs[i] = pt.TagID
	}
	// batch load tags via repo
	tags, _ := s.tagRepo.FindByIDs(tagIDs)
	items := make([]model.TagResponse, 0, len(tags))
	for _, t := range tags {
		items = append(items, model.TagResponse{
			ID: t.ID, Name: t.Name, Description: t.Description, PostCount: t.PostCount,
		})
	}
	return items, nil
}

func (s *TagService) GetPostsByTag(tagID uint, page, pageSize int) (*model.PostListResponse, error) {
	pts, _, err := s.postTagRepo.FindByTagID(tagID, page, pageSize)
	if err != nil {
		return nil, err
	}
	if len(pts) == 0 {
		return &model.PostListResponse{Total: 0, Items: []model.PostListItem{}}, nil
	}
	postIDs := make([]uint, len(pts))
	for i, pt := range pts {
		postIDs[i] = pt.PostID
	}
	posts, _, err := s.postRepo.List(repository.PostListQuery{
		Page: 1, PageSize: len(postIDs), StatusExclude: 2,
	})
	postMap := make(map[uint]mysql.Post)
	for _, p := range posts {
		postMap[p.ID] = p
	}
	// collect user IDs
	userIDs := make([]uint, 0, len(postMap))
	for _, p := range postMap {
		userIDs = append(userIDs, p.UserID)
	}
	users, _ := s.userRepo.FindByIDs(userIDs)
	userMap := make(map[uint]mysql.User)
	for _, u := range users {
		userMap[u.ID] = u
	}
	items := make([]model.PostListItem, 0, len(postIDs))
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
	return &model.PostListResponse{Total: int64(len(items)), Items: items}, nil
}
