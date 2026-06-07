package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"community-server/internal/model"

	"github.com/gin-gonic/gin"
)

// ============================================
// mock services
// ============================================

type mockUserSvc struct {
	registerFn   func(req *model.RegisterRequest) (uint, error)
	loginFn      func(req *model.LoginRequest) (*model.LoginResponse, error)
	getUserFn    func(userID uint) (*model.UserProfileResponse, error)
	updateProfFn func(userID uint, req *model.UpdateProfileRequest) error
}

func (m *mockUserSvc) Register(req *model.RegisterRequest) (uint, error) {
	if m.registerFn != nil {
		return m.registerFn(req)
	}
	return 1, nil
}
func (m *mockUserSvc) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	if m.loginFn != nil {
		return m.loginFn(req)
	}
	return &model.LoginResponse{Token: "test-token", UserID: 1, Username: "test"}, nil
}
func (m *mockUserSvc) GetUserByID(userID uint) (*model.UserProfileResponse, error) {
	if m.getUserFn != nil {
		return m.getUserFn(userID)
	}
	return &model.UserProfileResponse{ID: 1, Username: "test", Nickname: "Test"}, nil
}
func (m *mockUserSvc) UpdateProfile(userID uint, req *model.UpdateProfileRequest) error {
	if m.updateProfFn != nil {
		return m.updateProfFn(userID, req)
	}
	return nil
}

func (m *mockUserSvc) ForgotPassword(email string) error { return nil }
func (m *mockUserSvc) ResetPassword(token, newPassword string) error { return nil }

type mockPostSvc struct {
	createPostFn   func(userID uint, req *model.CreatePostRequest) (uint, error)
	getPostListFn  func(req *model.PostListRequest) (*model.PostListResponse, error)
	getPostFn      func(postID uint) (*model.PostDetailResponse, error)
	updatePostFn   func(userID, postID uint, req *model.UpdatePostRequest) error
	deletePostFn   func(userID, postID uint) error
	likePostFn     func(userID, postID uint) error
	unlikePostFn   func(userID, postID uint) error
	favoritePostFn func(userID, postID uint) error
	unfavPostFn    func(userID, postID uint) error
	getUserPostsFn func(userID uint, req *model.PostListRequest) (*model.PostListResponse, error)
	favoritesFn    func(userID uint, page, pageSize int) (*model.PostListResponse, error)
	likedFn        func(userID uint, page, pageSize int) (*model.PostListResponse, error)
	feedFn         func(userID uint, page, pageSize int) (*model.PostListResponse, error)
}

func (m *mockPostSvc) CreatePost(userID uint, req *model.CreatePostRequest) (uint, error) {
	if m.createPostFn != nil {
		return m.createPostFn(userID, req)
	}
	return 1, nil
}
func (m *mockPostSvc) GetPostList(req *model.PostListRequest) (*model.PostListResponse, error) {
	if m.getPostListFn != nil {
		return m.getPostListFn(req)
	}
	return &model.PostListResponse{}, nil
}
func (m *mockPostSvc) GetPost(postID uint) (*model.PostDetailResponse, error) {
	if m.getPostFn != nil {
		return m.getPostFn(postID)
	}
	return &model.PostDetailResponse{ID: postID, Title: "test post"}, nil
}
func (m *mockPostSvc) UpdatePost(userID, postID uint, req *model.UpdatePostRequest) error {
	if m.updatePostFn != nil {
		return m.updatePostFn(userID, postID, req)
	}
	return nil
}
func (m *mockPostSvc) DeletePost(userID, postID uint) error {
	if m.deletePostFn != nil {
		return m.deletePostFn(userID, postID)
	}
	return nil
}
func (m *mockPostSvc) LikePost(userID, postID uint) error {
	if m.likePostFn != nil {
		return m.likePostFn(userID, postID)
	}
	return nil
}
func (m *mockPostSvc) UnlikePost(userID, postID uint) error {
	if m.unlikePostFn != nil {
		return m.unlikePostFn(userID, postID)
	}
	return nil
}
func (m *mockPostSvc) FavoritePost(userID, postID uint) error {
	if m.favoritePostFn != nil {
		return m.favoritePostFn(userID, postID)
	}
	return nil
}
func (m *mockPostSvc) UnfavoritePost(userID, postID uint) error {
	if m.unfavPostFn != nil {
		return m.unfavPostFn(userID, postID)
	}
	return nil
}
func (m *mockPostSvc) GetUserPosts(userID uint, req *model.PostListRequest) (*model.PostListResponse, error) {
	if m.getUserPostsFn != nil {
		return m.getUserPostsFn(userID, req)
	}
	return &model.PostListResponse{}, nil
}
func (m *mockPostSvc) GetUserFavorites(userID uint, page, pageSize int) (*model.PostListResponse, error) {
	if m.favoritesFn != nil {
		return m.favoritesFn(userID, page, pageSize)
	}
	return &model.PostListResponse{}, nil
}
func (m *mockPostSvc) GetUserLikedPosts(userID uint, page, pageSize int) (*model.PostListResponse, error) {
	if m.likedFn != nil {
		return m.likedFn(userID, page, pageSize)
	}
	return &model.PostListResponse{}, nil
}
func (m *mockPostSvc) GetFollowFeed(userID uint, page, pageSize int) (*model.PostListResponse, error) {
	if m.feedFn != nil {
		return m.feedFn(userID, page, pageSize)
	}
	return &model.PostListResponse{}, nil
}

// ============================================
// 用户接口测试
// ============================================

func setupUserHandler(svc UserService) (*gin.Engine, *UserHandler) {
	gin.SetMode(gin.TestMode)
	h := NewUserHandler(svc)
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/profile", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		h.GetProfile(c)
	})
	return r, h
}

func TestUserRegister_Success(t *testing.T) {
	svc := &mockUserSvc{
		registerFn: func(req *model.RegisterRequest) (uint, error) {
			if req.Username == "" {
				return 0, errors.New("用户名不能为空")
			}
			return 42, nil
		},
	}
	r, _ := setupUserHandler(svc)

	body := `{"username":"newuser","password":"pass123","email":"a@b.com"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望 200，实际 %d", w.Code)
	}

	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("期望 code 0，实际 %d", resp.Code)
	}
	if resp.Data["user_id"] != float64(42) {
		t.Errorf("期望 user_id 42，实际 %v", resp.Data["user_id"])
	}
}

func TestUserRegister_InvalidParam(t *testing.T) {
	svc := &mockUserSvc{}
	r, _ := setupUserHandler(svc)

	// 缺少必填字段 username
	body := `{"password":"123"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code == 0 {
		t.Error("缺少必填字段应该返回错误")
	}
}

func TestUserLogin_Success(t *testing.T) {
	svc := &mockUserSvc{
		loginFn: func(req *model.LoginRequest) (*model.LoginResponse, error) {
			return &model.LoginResponse{
				Token: "jwt-token", UserID: 1, Username: "test", Nickname: "Tester",
			}, nil
		},
	}
	r, _ := setupUserHandler(svc)

	body := `{"username":"test","password":"pass"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code int                     `json:"code"`
		Data *model.LoginResponse    `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("期望 code 0，实际 %d", resp.Code)
	}
	if resp.Data.Token != "jwt-token" {
		t.Errorf("期望 token jwt-token，实际 %s", resp.Data.Token)
	}
}

func TestUserLogin_WrongPassword(t *testing.T) {
	svc := &mockUserSvc{
		loginFn: func(req *model.LoginRequest) (*model.LoginResponse, error) {
			return nil, errors.New("用户名或密码错误")
		},
	}
	r, _ := setupUserHandler(svc)

	body := `{"username":"test","password":"wrong"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 401 {
		t.Errorf("期望 401，实际 %d", resp.Code)
	}
}

func TestUserProfile_Unauthenticated(t *testing.T) {
	svc := &mockUserSvc{}
	gin.SetMode(gin.TestMode)
	h := NewUserHandler(svc)
	r := gin.New()
	r.GET("/profile", h.GetProfile) // 不设 user_id

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/profile", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 401 {
		t.Errorf("未登录应该返回 401，实际 %d", resp.Code)
	}
}

func TestUserProfile_Success(t *testing.T) {
	svc := &mockUserSvc{
		getUserFn: func(uid uint) (*model.UserProfileResponse, error) {
			return &model.UserProfileResponse{
				ID: 1, Username: "test", Nickname: "Tester",
				Avatar: "https://example.com/avatar.png", Bio: "hello",
				Email: "a@b.com", Status: 1,
			}, nil
		},
	}
	r, _ := setupUserHandler(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/profile", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int                          `json:"code"`
		Data *model.UserProfileResponse   `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("期望 code 0，实际 %d", resp.Code)
	}
	if resp.Data.Username != "test" {
		t.Errorf("期望 username test，实际 %s", resp.Data.Username)
	}
}

// ============================================
// 帖子接口测试
// ============================================

func setupPostHandler(svc PostService) (*gin.Engine, *PostHandler) {
	gin.SetMode(gin.TestMode)
	h := NewPostHandler(svc)
	r := gin.New()
	r.GET("/posts/:id", h.GetPost)
	r.POST("/posts", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		h.CreatePost(c)
	})
	return r, h
}

func TestPostGet_Success(t *testing.T) {
	svc := &mockPostSvc{
		getPostFn: func(id uint) (*model.PostDetailResponse, error) {
			return &model.PostDetailResponse{
				ID: id, Title: "Test Title", Content: "Test content",
				Username: "author", Nickname: "Author", Status: 1,
			}, nil
		},
	}
	r, _ := setupPostHandler(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/posts/1", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int                     `json:"code"`
		Data *model.PostDetailResponse `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("期望 code 0，实际 %d", resp.Code)
	}
	if resp.Data.Title != "Test Title" {
		t.Errorf("期望 title 'Test Title'，实际 %s", resp.Data.Title)
	}
}

func TestPostGet_NotFound(t *testing.T) {
	svc := &mockPostSvc{
		getPostFn: func(id uint) (*model.PostDetailResponse, error) {
			return nil, errors.New("帖子不存在")
		},
	}
	r, _ := setupPostHandler(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/posts/999", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 404 {
		t.Errorf("帖子不存在应该返回 404，实际 %d", resp.Code)
	}
}

func TestPostGet_InvalidID(t *testing.T) {
	svc := &mockPostSvc{}
	r, _ := setupPostHandler(svc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/posts/abc", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 400 {
		t.Errorf("无效ID应该返回 400，实际 %d", resp.Code)
	}
}

func TestPostCreate_Unauthenticated(t *testing.T) {
	svc := &mockPostSvc{}
	gin.SetMode(gin.TestMode)
	h := NewPostHandler(svc)
	r := gin.New()
	r.POST("/posts", h.CreatePost) // 不设 user_id

	body := `{"title":"t","content":"c"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 401 {
		t.Errorf("未登录应该返回 401，实际 %d", resp.Code)
	}
}

func TestPostCreate_Success(t *testing.T) {
	svc := &mockPostSvc{
		createPostFn: func(uid uint, req *model.CreatePostRequest) (uint, error) {
			return 100, nil
		},
	}
	r, _ := setupPostHandler(svc)

	body := `{"title":"New Post","content":"Hello world"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/posts", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("期望 code 0，实际 %d", resp.Code)
	}
	if resp.Data["post_id"] != float64(100) {
		t.Errorf("期望 post_id 100，实际 %v", resp.Data["post_id"])
	}
}
