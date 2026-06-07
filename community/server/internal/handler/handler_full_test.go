package handler

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"community-server/internal/model"

	"github.com/gin-gonic/gin"
)

// ============================================
// 评论接口测试
// ============================================

type mockCommentSvc struct {
	createFn  func(userID uint, req *model.CreateCommentRequest) (uint, error)
	listFn    func(req *model.CommentListRequest) (*model.CommentListResponse, error)
	updateFn  func(userID, commentID uint, req *model.UpdateCommentRequest) error
	deleteFn  func(userID, commentID uint) error
	likeFn    func(userID, commentID uint) error
	unlikeFn  func(userID, commentID uint) error
}

func (m *mockCommentSvc) CreateComment(uid uint, r *model.CreateCommentRequest) (uint, error) {
	if m.createFn != nil { return m.createFn(uid, r) }; return 1, nil
}
func (m *mockCommentSvc) GetCommentList(r *model.CommentListRequest) (*model.CommentListResponse, error) {
	if m.listFn != nil { return m.listFn(r) }; return &model.CommentListResponse{}, nil
}
func (m *mockCommentSvc) UpdateComment(uid, cid uint, r *model.UpdateCommentRequest) error {
	if m.updateFn != nil { return m.updateFn(uid, cid, r) }; return nil
}
func (m *mockCommentSvc) DeleteComment(uid, cid uint) error {
	if m.deleteFn != nil { return m.deleteFn(uid, cid) }; return nil
}
func (m *mockCommentSvc) LikeComment(uid, cid uint) error {
	if m.likeFn != nil { return m.likeFn(uid, cid) }; return nil
}
func (m *mockCommentSvc) UnlikeComment(uid, cid uint) error {
	if m.unlikeFn != nil { return m.unlikeFn(uid, cid) }; return nil
}

func setupCommentHandler(svc CommentService) (*gin.Engine, *CommentHandler) {
	gin.SetMode(gin.TestMode)
	h := &CommentHandler{commentService: svc}
	r := gin.New()
	r.POST("/comments", func(c *gin.Context) { c.Set("user_id", uint(1)); h.CreateComment(c) })
	r.GET("/comments", h.GetCommentList)
	r.PUT("/comments/:id", func(c *gin.Context) { c.Set("user_id", uint(1)); h.UpdateComment(c) })
	r.DELETE("/comments/:id", func(c *gin.Context) { c.Set("user_id", uint(1)); h.DeleteComment(c) })
	r.POST("/comments/:id/like", func(c *gin.Context) { c.Set("user_id", uint(1)); h.LikeComment(c) })
	r.DELETE("/comments/:id/like", func(c *gin.Context) { c.Set("user_id", uint(1)); h.UnlikeComment(c) })
	return r, h
}

func TestCommentCreate_Success(t *testing.T) {
	svc := &mockCommentSvc{createFn: func(uid uint, r *model.CreateCommentRequest) (uint, error) { return 42, nil }}
	r, _ := setupCommentHandler(svc)
	body := `{"post_id":1,"content":"nice post"}`
	w := doPost(r, "/comments", body)
	assertCode(t, w, 0)
}

func TestCommentCreate_Unauthenticated(t *testing.T) {
	svc := &mockCommentSvc{}
	gin.SetMode(gin.TestMode)
	h := &CommentHandler{commentService: svc}
	rr := gin.New()
	rr.POST("/comments", h.CreateComment)
	w := doPost(rr, "/comments", `{"post_id":1}`)
	assertCode(t, w, 401)
}

func TestCommentList_Success(t *testing.T) {
	svc := &mockCommentSvc{listFn: func(r *model.CommentListRequest) (*model.CommentListResponse, error) {
		return &model.CommentListResponse{Total: 5, Items: []model.CommentListItem{}}, nil
	}}
	r, _ := setupCommentHandler(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/comments?post_id=1&page=1", nil)
	r.ServeHTTP(w, req)
	assertCode(t, w, 0)
}

func TestCommentDelete_NotFound(t *testing.T) {
	svc := &mockCommentSvc{deleteFn: func(uid, cid uint) error { return errors.New("评论不存在") }}
	r, _ := setupCommentHandler(svc)
	w := httptest.NewRequest("DELETE", "/comments/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, w)
	assertCode(t, rr, 500)
}

func TestCommentLike_Success(t *testing.T) {
	svc := &mockCommentSvc{}
	r, _ := setupCommentHandler(svc)
	w := httptest.NewRequest("POST", "/comments/1/like", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, w)
	assertCode(t, rr, 0)
}

func TestCommentUnlike_Success(t *testing.T) {
	svc := &mockCommentSvc{}
	r, _ := setupCommentHandler(svc)
	w := httptest.NewRequest("DELETE", "/comments/1/like", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, w)
	assertCode(t, rr, 0)
}

// ============================================
// 关注接口测试
// ============================================

type mockFollowSvc struct {
	followFn   func(userID, followID uint) error
	unfollowFn func(userID, followID uint) error
}

func (m *mockFollowSvc) FollowUser(userID uint, req *model.FollowRequest) error {
	if m.followFn != nil { return m.followFn(userID, req.FollowID) }; return nil
}
func (m *mockFollowSvc) UnfollowUser(userID, followID uint) error {
	if m.unfollowFn != nil { return m.unfollowFn(userID, followID) }; return nil
}
func (m *mockFollowSvc) GetFollowers(userID uint, page, pageSize int) (*model.FollowListResponse, error) {
	return &model.FollowListResponse{}, nil
}
func (m *mockFollowSvc) GetFollowing(userID uint, page, pageSize int) (*model.FollowListResponse, error) {
	return &model.FollowListResponse{}, nil
}
func (m *mockFollowSvc) IsFollowing(userID, targetID uint) (bool, error) { return false, nil }
func (m *mockFollowSvc) GetFollowCounts(userID uint) (int64, int64) { return 0, 0 }

func setupFollowHandler(svc FollowService) (*gin.Engine, *FollowHandler) {
	gin.SetMode(gin.TestMode)
	h := &FollowHandler{followService: svc}
	r := gin.New()
	r.POST("/follows", func(c *gin.Context) { c.Set("user_id", uint(1)); h.FollowUser(c) })
	r.DELETE("/follows/:id", func(c *gin.Context) { c.Set("user_id", uint(1)); h.UnfollowUser(c) })
	return r, h
}

func TestFollowUser_Success(t *testing.T) {
	svc := &mockFollowSvc{}
	r, _ := setupFollowHandler(svc)
	w := doPost(r, "/follows", `{"follow_id":2}`)
	assertCode(t, w, 0)
}

func TestFollowUser_SelfFollow(t *testing.T) {
	svc := &mockFollowSvc{followFn: func(uid, fid uint) error { return errors.New("不能关注自己") }}
	r, _ := setupFollowHandler(svc)
	w := doPost(r, "/follows", `{"follow_id":2}`)
	assertCode(t, w, 500)
}

func TestUnfollowUser_Success(t *testing.T) {
	svc := &mockFollowSvc{}
	r, _ := setupFollowHandler(svc)
	w := httptest.NewRequest("DELETE", "/follows/2", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, w)
	assertCode(t, rr, 0)
}

// ============================================
// 通知接口测试
// ============================================

type mockNotifSvc struct {
	listFn       func(userID uint, page, pageSize int) (*model.NotificationListResponse, error)
	markFn       func(notifID, userID uint) error
	markAllFn    func(userID uint) error
	unreadFn     func(userID uint) (int64, error)
}

func (m *mockNotifSvc) GetList(userID uint, page, pageSize int) (*model.NotificationListResponse, error) {
	if m.listFn != nil { return m.listFn(userID, page, pageSize) }
	return &model.NotificationListResponse{Total: 0, Items: []model.NotificationResponse{}, Unread: 0}, nil
}
func (m *mockNotifSvc) MarkRead(notifID, userID uint) error {
	if m.markFn != nil { return m.markFn(notifID, userID) }; return nil
}
func (m *mockNotifSvc) MarkAllRead(userID uint) error {
	if m.markAllFn != nil { return m.markAllFn(userID) }; return nil
}
func (m *mockNotifSvc) GetUnreadCount(userID uint) (int64, error) {
	if m.unreadFn != nil { return m.unreadFn(userID) }; return 5, nil
}

func setupNotifHandler(svc *mockNotifSvc) (*gin.Engine, *NotificationHandler) {
	gin.SetMode(gin.TestMode)
	h := &NotificationHandler{notificationService: svc}
	r := gin.New()
	r.GET("/notifications", func(c *gin.Context) { c.Set("user_id", uint(1)); h.GetList(c) })
	r.PUT("/notifications/:id/read", func(c *gin.Context) { c.Set("user_id", uint(1)); h.MarkRead(c) })
	r.PUT("/notifications/read-all", func(c *gin.Context) { c.Set("user_id", uint(1)); h.MarkAllRead(c) })
	r.GET("/notifications/unread", func(c *gin.Context) { c.Set("user_id", uint(1)); h.GetUnreadCount(c) })
	return r, h
}

func TestNotificationList_Success(t *testing.T) {
	svc := &mockNotifSvc{}
	r, _ := setupNotifHandler(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/notifications?page=1", nil)
	r.ServeHTTP(w, req)
	assertCode(t, w, 0)
}

func TestNotificationUnauthenticated(t *testing.T) {
	svc := &mockNotifSvc{}
	gin.SetMode(gin.TestMode)
	h := &NotificationHandler{notificationService: svc}
	r := gin.New()
	r.GET("/notifications", h.GetList)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/notifications", nil)
	r.ServeHTTP(w, req)
	assertCode(t, w, 401)
}

func TestNotificationMarkRead_Success(t *testing.T) {
	svc := &mockNotifSvc{}
	r, _ := setupNotifHandler(svc)
	w := httptest.NewRequest("PUT", "/notifications/1/read", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, w)
	assertCode(t, rr, 0)
}

func TestNotificationMarkAllRead_Success(t *testing.T) {
	svc := &mockNotifSvc{}
	r, _ := setupNotifHandler(svc)
	w := httptest.NewRequest("PUT", "/notifications/read-all", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, w)
	assertCode(t, rr, 0)
}

func TestNotificationUnreadCount(t *testing.T) {
	svc := &mockNotifSvc{unreadFn: func(uid uint) (int64, error) { return 7, nil }}
	r, _ := setupNotifHandler(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/notifications/unread", nil)
	r.ServeHTTP(w, req)
	var resp struct { Code int; Data model.UnreadCountResponse `json:"data"` }
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 { t.Errorf("expected code 0, got %d", resp.Code) }
	if resp.Data.Count != 7 { t.Errorf("expected count 7, got %d", resp.Data.Count) }
}

// ============================================
// 管理接口测试
// ============================================

type mockAdminSvc struct {
	userListFn   func(req *model.AdminUserListRequest) (*model.AdminUserListResponse, error)
	deleteUserFn func(userID uint) error
}

func (m *mockAdminSvc) GetUserList(req *model.AdminUserListRequest) (*model.AdminUserListResponse, error) {
	if m.userListFn != nil { return m.userListFn(req) }
	return &model.AdminUserListResponse{Total: 0, Items: []model.AdminUserInfo{}}, nil
}
func (m *mockAdminSvc) DeleteUser(userID uint) error {
	if m.deleteUserFn != nil { return m.deleteUserFn(userID) }; return nil
}
func (m *mockAdminSvc) UpdateUserAdminType(userID uint, adminType int) error { return nil }
func (m *mockAdminSvc) UpdateUserStatus(userID uint, status int) error { return nil }
func (m *mockAdminSvc) GetPostList(req *model.AdminPostListRequest) (*model.AdminPostListResponse, error) {
	return &model.AdminPostListResponse{}, nil
}
func (m *mockAdminSvc) DeletePost(postID uint) error { return nil }
func (m *mockAdminSvc) SetPostTop(postID uint, isTop bool) error { return nil }
func (m *mockAdminSvc) SetPostEssence(postID uint, isEssence bool) error { return nil }

type mockIM struct{}
func (m *mockIM) RegisterUser(uid, nick string) error { return nil }
func (m *mockIM) SendPrivateMsg(s, t, c string) error { return nil }
func (m *mockIM) SendSystemMsg(s, t, c string) error { return nil }
func (m *mockIM) SendGroupMsg(s, g, c string) error { return nil }
func (m *mockIM) QueryUserInfo(uid string) (map[string]interface{}, error) { return nil, nil }
func (m *mockIM) CreateGroup(gid, name, owner string, members []string) error { return nil }
func (m *mockIM) QueryOnlineStatus(uids []string) (map[string]bool, error) { return nil, nil }
func (m *mockIM) SendBroadcastMsg(sender, content string) error { return nil }
func (m *mockIM) AddBot(botID, name, _ string) error { return nil }

func setupAdminHandler(svc *mockAdminSvc) (*gin.Engine, *AdminHandler) {
	gin.SetMode(gin.TestMode)
	h := &AdminHandler{adminService: svc, imClient: &mockIM{}}
	r := gin.New()
	r.GET("/admin/users", h.GetUserList)
	r.DELETE("/admin/users/:id", h.DeleteUser)
	return r, h
}

func TestAdminGetUserList_Success(t *testing.T) {
	svc := &mockAdminSvc{
		userListFn: func(req *model.AdminUserListRequest) (*model.AdminUserListResponse, error) {
			return &model.AdminUserListResponse{
				Total: 2,
				Items: []model.AdminUserInfo{
					{ID: 1, Username: "admin", AdminType: 1, Status: 1},
					{ID: 2, Username: "user1", AdminType: 0, Status: 1},
				},
			}, nil
		},
	}
	r, _ := setupAdminHandler(svc)
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/admin/users?page=1", nil)
	r.ServeHTTP(w, req)
	assertCode(t, w, 0)
}

func TestAdminDeleteUser_Success(t *testing.T) {
	svc := &mockAdminSvc{}
	r, _ := setupAdminHandler(svc)
	w := httptest.NewRequest("DELETE", "/admin/users/2", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, w)
	assertCode(t, rr, 0)
}

// ============================================
// helpers
// ============================================

func doPost(r *gin.Engine, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func assertCode(t *testing.T, w *httptest.ResponseRecorder, expected int) {
	t.Helper()
	var resp struct{ Code int }
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != expected {
		t.Errorf("expected code %d, got %d (body: %s)", expected, resp.Code, w.Body.String())
	}
}
