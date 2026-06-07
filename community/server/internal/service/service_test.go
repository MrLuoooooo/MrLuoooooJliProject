package service

import (
	"errors"
	"testing"

	"community-server/internal/db/mysql"
	"community-server/internal/model"
	"community-server/internal/repository"
)

// ============================================
// 桩
// ============================================

// mockNotifPush 占位的通知推送函数
// 服务层测试不关心 IM 推送，直接传 nil 就行，测的是业务逻辑本身
func mockNotifPush(userID, fromID uint, ntype int, targetID uint, content string) {}

type testPostRepo struct {
	posts     map[uint]*mysql.Post
	createErr error
}

func newTestPostRepo() *testPostRepo { return &testPostRepo{posts: map[uint]*mysql.Post{}} }

func (r *testPostRepo) Create(post *mysql.Post) error {
	if r.createErr != nil { return r.createErr }
	post.ID = uint(len(r.posts)) + 1
	r.posts[post.ID] = post
	return nil
}
func (r *testPostRepo) FindByID(id uint) (*mysql.Post, error) {
	if p, ok := r.posts[id]; ok { return p, nil }
	return nil, errors.New("not found")
}
func (r *testPostRepo) Update(id uint, m map[string]interface{}) error       { return nil }
func (r *testPostRepo) UpdateColumn(id uint, col string, val interface{}) error { return nil }
func (r *testPostRepo) List(q repository.PostListQuery) ([]mysql.Post, int64, error) { return nil, 0, nil }
func (r *testPostRepo) Search(kw string, page, size int) ([]mysql.Post, int64, error) { return nil, 0, nil }
func (r *testPostRepo) Delete(id uint) error   { delete(r.posts, id); return nil }
func (r *testPostRepo) SoftDelete(id uint) error { return r.Delete(id) }

type testUserRepo struct {
	users map[uint]mysql.User
}

func newTestUserRepo() *testUserRepo { return &testUserRepo{users: map[uint]mysql.User{}} }
func (r *testUserRepo) add(id uint, username, nickname string) {
	r.users[id] = mysql.User{BaseModel: mysql.BaseModel{ID: id}, Username: username, Nickname: nickname}
}
func (r *testUserRepo) Create(u *mysql.User) error                  { return nil }
func (r *testUserRepo) FindByUsername(s string) (*mysql.User, error) { return nil, errors.New("not found") }
func (r *testUserRepo) FindByEmail(s string) (*mysql.User, error)    { return nil, errors.New("not found") }
func (r *testUserRepo) FindByID(id uint) (*mysql.User, error) {
	if u, ok := r.users[id]; ok { return &u, nil }
	return nil, errors.New("not found")
}
func (r *testUserRepo) Update(id uint, m map[string]interface{}) error { return nil }
func (r *testUserRepo) Delete(id uint) error                           { return nil }
func (r *testUserRepo) FindByIDs(ids []uint) ([]mysql.User, error) {
	result := make([]mysql.User, 0, len(ids))
	for _, id := range ids {
		if u, ok := r.users[id]; ok { result = append(result, u) }
	}
	return result, nil
}
func (r *testUserRepo) Search(kw string, page, size int) ([]mysql.User, int64, error) { return nil, 0, nil }
func (r *testUserRepo) List(page, size int) ([]mysql.User, int64, error)              { return nil, 0, nil }

type testCategoryRepo struct{}

func (r *testCategoryRepo) Create(c *mysql.Category) error                       { return nil }
func (r *testCategoryRepo) FindByID(id uint) (*mysql.Category, error)             { return nil, nil }
func (r *testCategoryRepo) Update(id uint, m map[string]interface{}) error        { return nil }
func (r *testCategoryRepo) Delete(id uint) error                                  { return nil }
func (r *testCategoryRepo) List() ([]mysql.Category, error)                       { return nil, nil }
func (r *testCategoryRepo) IncrementPostCount(catID uint, delta int) error        { return nil }

type testFollowRepo struct{}

func (r *testFollowRepo) Create(f *mysql.UserFollow) error                   { return nil }
func (r *testFollowRepo) Delete(userID, followID uint) error                 { return nil }
func (r *testFollowRepo) Exists(userID, followID uint) (bool, error)         { return false, nil }
func (r *testFollowRepo) CountFollowers(userID uint) int64                   { return 0 }
func (r *testFollowRepo) CountFollowing(userID uint) int64                   { return 0 }
func (r *testFollowRepo) FindFollowers(userID uint, page, size int) ([]mysql.UserFollow, int64, error) { return nil, 0, nil }
func (r *testFollowRepo) FindFollowing(userID uint, page, size int) ([]mysql.UserFollow, int64, error) { return nil, 0, nil }
func (r *testFollowRepo) FindAllFollowing(userID uint) ([]mysql.UserFollow, error) { return nil, nil }

type testCommentRepo struct {
	comments       map[uint]*mysql.Comment
	softDelPostIDs []uint
}

func newTestCommentRepo() *testCommentRepo { return &testCommentRepo{comments: map[uint]*mysql.Comment{}} }
func (r *testCommentRepo) Create(c *mysql.Comment) error { c.ID = 1; r.comments[c.ID] = c; return nil }
func (r *testCommentRepo) FindByID(id uint) (*mysql.Comment, error) {
	if c, ok := r.comments[id]; ok { return c, nil }
	return nil, errors.New("not found")
}
func (r *testCommentRepo) FindRootByPostID(postID uint, page, size int) ([]mysql.Comment, int64, error) { return nil, 0, nil }
func (r *testCommentRepo) FindRepliesByPostAndParents(postID uint, parentIDs []uint) ([]mysql.Comment, error) { return nil, nil }
func (r *testCommentRepo) Update(id uint, m map[string]interface{}) error       { return nil }
func (r *testCommentRepo) UpdateColumn(id uint, col string, val interface{}) error { return nil }
func (r *testCommentRepo) SoftDelete(id uint) error                             { delete(r.comments, id); return nil }
func (r *testCommentRepo) SoftDeleteByPostID(postID uint) error                 { r.softDelPostIDs = append(r.softDelPostIDs, postID); return nil }

type testResetRepo struct{ resets map[string]*mysql.PasswordReset }

func newTestResetRepo() *testResetRepo { return &testResetRepo{resets: map[string]*mysql.PasswordReset{}} }
func (r *testResetRepo) Create(reset *mysql.PasswordReset) error { r.resets[reset.Token] = reset; return nil }
func (r *testResetRepo) FindByToken(token string) (*mysql.PasswordReset, error) {
	if rs, ok := r.resets[token]; ok { return rs, nil }
	return nil, errors.New("not found")
}
func (r *testResetRepo) MarkUsed(id uint) error { return nil }

// ============================================
// 帖子服务测试
// ============================================

func makePostSvc(postRepo *testPostRepo, userRepo *testUserRepo) *PostService {
	svc := NewPostService(nil, postRepo, nil, nil, newTestCommentRepo(), userRepo, &testFollowRepo{}, &testCategoryRepo{})
	svc.notifSvc = nil
	return svc
}

func TestPostService_CreatePost_Success(t *testing.T) {
	repo := newTestPostRepo()
	svc := makePostSvc(repo, newTestUserRepo())

	id, err := svc.CreatePost(1, &model.CreatePostRequest{Title: "Test", Content: "Content"})
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}
	if id == 0 { t.Error("expected non-zero id") }
	if _, err := repo.FindByID(id); err != nil { t.Error("post not persisted") }
}

func TestPostService_CreatePost_RepoError(t *testing.T) {
	repo := newTestPostRepo()
	repo.createErr = errors.New("db error")
	svc := makePostSvc(repo, newTestUserRepo())

	_, err := svc.CreatePost(1, &model.CreatePostRequest{Title: "x", Content: "y"})
	if err == nil { t.Error("expected error on repo failure") }
}

func TestPostService_GetPost_Success(t *testing.T) {
	repo := newTestPostRepo()
	svc := makePostSvc(repo, newTestUserRepo())
	id, _ := svc.CreatePost(1, &model.CreatePostRequest{Title: "T", Content: "C"})

	post, err := svc.GetPostByID(id)
	if err != nil { t.Fatalf("GetPostByID: %v", err) }
	if post.Title != "T" { t.Errorf("title = %q", post.Title) }
}

func TestPostService_GetPost_NotFound(t *testing.T) {
	svc := makePostSvc(newTestPostRepo(), newTestUserRepo())
	_, err := svc.GetPostByID(999)
	if err == nil { t.Error("expected error for non-existent post") }
}

func TestPostService_UpdatePost_NotOwner(t *testing.T) {
	svc := makePostSvc(newTestPostRepo(), newTestUserRepo())
	id, _ := svc.CreatePost(1, &model.CreatePostRequest{Title: "T", Content: "C"})

	err := svc.UpdatePost(2, id, &model.UpdatePostRequest{Title: "Hacked"})
	if err == nil { t.Error("expected permission denied") }
}

func TestPostService_DeletePost_NotOwner(t *testing.T) {
	svc := makePostSvc(newTestPostRepo(), newTestUserRepo())
	id, _ := svc.CreatePost(1, &model.CreatePostRequest{Title: "T", Content: "C"})

	err := svc.DeletePost(2, id)
	if err == nil { t.Error("expected permission denied") }
}

func TestPostService_DeletePost_Success(t *testing.T) {
	repo := newTestPostRepo()
	svc := makePostSvc(repo, newTestUserRepo())
	id, _ := svc.CreatePost(1, &model.CreatePostRequest{Title: "T", Content: "C"})

	err := svc.DeletePost(1, id)
	if err != nil { t.Fatalf("DeletePost: %v", err) }
	_, err = repo.FindByID(id)
	if err == nil { t.Error("post should be deleted") }
}

func TestPostService_buildPostList_LoadsUsers(t *testing.T) {
	userRepo := newTestUserRepo()
	userRepo.add(1, "alice", "Alice")
	userRepo.add(2, "bob", "Bob")
	svc := makePostSvc(newTestPostRepo(), userRepo)

	posts := []mysql.Post{
		{BaseModel: mysql.BaseModel{ID: 1}, UserID: 1, Title: "A", Status: 1},
		{BaseModel: mysql.BaseModel{ID: 2}, UserID: 2, Title: "B", Status: 1},
	}
	resp := svc.buildPostList(posts, 2)
	if len(resp.Items) != 2 { t.Fatalf("items = %d", len(resp.Items)) }
	if resp.Items[0].Username != "alice" { t.Errorf("u0 = %q", resp.Items[0].Username) }
	if resp.Items[1].Nickname != "Bob" { t.Errorf("u1 = %q", resp.Items[1].Nickname) }
}

func TestPostService_buildPostList_Empty(t *testing.T) {
	svc := makePostSvc(newTestPostRepo(), newTestUserRepo())
	resp := svc.buildPostList(nil, 0)
	if len(resp.Items) != 0 { t.Error("expected empty items") }
}

func TestPostService_GetFollowFeed_NoFollowing(t *testing.T) {
	svc := makePostSvc(newTestPostRepo(), newTestUserRepo())
	resp, err := svc.GetFollowFeed(1, 1, 20)
	if err != nil { t.Fatal(err) }
	if resp.Total != 0 { t.Errorf("total = %d", resp.Total) }
}

func TestPostService_UpdatePost_NoOp(t *testing.T) {
	svc := makePostSvc(newTestPostRepo(), newTestUserRepo())
	id, _ := svc.CreatePost(1, &model.CreatePostRequest{Title: "T", Content: "C"})

	err := svc.UpdatePost(1, id, &model.UpdatePostRequest{})
	if err != nil { t.Errorf("empty update should not error: %v", err) }
}

// ============================================
// 用户服务测试
// ============================================

func TestUserService_ForgotPassword_EmailNotFound(t *testing.T) {
	svc := NewUserService(nil, newTestUserRepo(), newTestResetRepo())
	err := svc.ForgotPassword("nonexistent@test.com")
	if err == nil { t.Error("expected error for unknown email") }
}

func TestUserService_ResetPassword_InvalidToken(t *testing.T) {
	svc := NewUserService(nil, newTestUserRepo(), newTestResetRepo())
	err := svc.ResetPassword("invalid-token", "newpass123")
	if err == nil { t.Error("expected error for invalid token") }
}

// ============================================
// 评论服务测试
// ============================================

func makeCommentSvc() (*CommentService, *testCommentRepo) {
	commentRepo := newTestCommentRepo()
	postRepo := newTestPostRepo()
	postRepo.posts[1] = &mysql.Post{BaseModel: mysql.BaseModel{ID: 1}, UserID: 1, Title: "T", Status: 1}
	svc := NewCommentService(nil, nil, commentRepo, nil, postRepo, newTestUserRepo())
	return svc, commentRepo
}

func TestCommentService_CreateComment_PostNotFound(t *testing.T) {
	svc, _ := makeCommentSvc()
	_, err := svc.CreateComment(1, &model.CreateCommentRequest{PostID: 999, Content: "test"})
	if err == nil { t.Error("expected error for non-existent post") }
}

func TestCommentService_CreateComment_Success(t *testing.T) {
	svc, repo := makeCommentSvc()
	id, err := svc.CreateComment(1, &model.CreateCommentRequest{PostID: 1, Content: "nice post"})
	if err != nil { t.Fatalf("CreateComment: %v", err) }
	if id == 0 { t.Error("expected non-zero id") }
	if _, err := repo.FindByID(id); err != nil { t.Error("comment not persisted") }
}

func TestCommentService_UpdateComment_NotOwner(t *testing.T) {
	svc, _ := makeCommentSvc()
	id, _ := svc.CreateComment(1, &model.CreateCommentRequest{PostID: 1, Content: "hello"})

	err := svc.UpdateComment(2, id, &model.UpdateCommentRequest{Content: "hacked"})
	if err == nil { t.Error("expected permission denied") }
}

func TestCommentService_DeleteComment_NotOwner(t *testing.T) {
	svc, _ := makeCommentSvc()
	id, _ := svc.CreateComment(1, &model.CreateCommentRequest{PostID: 1, Content: "hello"})

	err := svc.DeleteComment(2, id)
	if err == nil { t.Error("expected permission denied") }
}

func TestCommentService_DeleteComment_Success(t *testing.T) {
	svc, repo := makeCommentSvc()
	id, _ := svc.CreateComment(1, &model.CreateCommentRequest{PostID: 1, Content: "hello"})

	err := svc.DeleteComment(1, id)
	if err != nil { t.Fatalf("DeleteComment: %v", err) }
	_, err = repo.FindByID(id)
	if err == nil { t.Error("comment should be deleted") }
}
