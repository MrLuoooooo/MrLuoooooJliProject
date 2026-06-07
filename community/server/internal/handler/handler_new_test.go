package handler

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"community-server/internal/im"

	"github.com/gin-gonic/gin"
)

// mockIMClient 实现 im.IMClient 接口用于测试
type mockIMClient struct {
	statusFn func(uids []string) (map[string]bool, error)
}

func (m *mockIMClient) QueryOnlineStatus(uids []string) (map[string]bool, error) {
	if m.statusFn != nil {
		return m.statusFn(uids)
	}
	return map[string]bool{"1": true}, nil
}
func (m *mockIMClient) RegisterUser(uid, nick string) error { return nil }
func (m *mockIMClient) SendPrivateMsg(s, t, c string) error { return nil }
func (m *mockIMClient) SendSystemMsg(s, t, c string) error { return nil }
func (m *mockIMClient) SendGroupMsg(s, g, c string) error { return nil }
func (m *mockIMClient) QueryUserInfo(uid string) (map[string]interface{}, error) { return nil, nil }
func (m *mockIMClient) CreateGroup(g, n, o string, ms []string) error { return nil }
func (m *mockIMClient) SendBroadcastMsg(s, c string) error { return nil }
func (m *mockIMClient) AddBot(b, n, w string) error { return nil }

var _ im.IMClient = (*mockIMClient)(nil)

func TestOnlineStatus_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewStatusHandler(&mockIMClient{})
	r := gin.New()
	r.GET("/users/:id/online", h.GetOnlineStatus)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/1/online", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 0 {
		t.Errorf("期望 code 0，实际 %d", resp.Code)
	}
	isOnline, _ := resp.Data["is_online"].(bool)
	if !isOnline {
		t.Error("期望 is_online 为 true")
	}
}

func TestOnlineStatus_InvalidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewStatusHandler(&mockIMClient{})
	r := gin.New()
	r.GET("/users/:id/online", h.GetOnlineStatus)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/abc/online", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 400 {
		t.Errorf("无效 ID 应返回 400，实际 %d", resp.Code)
	}
}

func TestOnlineStatus_IMFails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewStatusHandler(&mockIMClient{
		statusFn: func(uids []string) (map[string]bool, error) {
			return map[string]bool{}, nil // 用户不在线
		},
	})
	r := gin.New()
	r.GET("/users/:id/online", h.GetOnlineStatus)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/users/999/online", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code == 0 && resp.Data != nil {
		isOnline, _ := resp.Data["is_online"].(bool)
		if isOnline {
			t.Error("不存在的用户应返回 is_online=false")
		}
	}
}

func TestBatchOnlineStatus_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewStatusHandler(&mockIMClient{
		statusFn: func(uids []string) (map[string]bool, error) {
			result := make(map[string]bool)
			for _, id := range uids {
				result[id] = true
			}
			return result, nil
		},
	})
	r := gin.New()
	r.POST("/users/online/batch", h.BatchOnlineStatus)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/users/online/batch", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Code int `json:"code"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Code != 400 {
		t.Errorf("空请求应返回 400，实际 %d", resp.Code)
	}
}
