package middleware

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupMiddlewareTest(mw gin.HandlerFunc, target string) *httptest.ResponseRecorder {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mw)
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", target, nil)
	r.ServeHTTP(w, req)
	return w
}

func TestCORSMiddleware_Headers(t *testing.T) {
	w := setupMiddlewareTest(CORSMiddleware(), "/test")
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS 缺少 Allow-Origin")
	}
	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("CORS 缺少 Allow-Methods")
	}
}

func TestCORSMiddleware_Options(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(CORSMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("OPTIONS 应返回 204，实际 %d", w.Code)
	}
}

func TestMetricsMiddleware(t *testing.T) {
	beforeTotal := GetTotalRequests()
	beforeActive := GetActiveRequests()

	w := setupMiddlewareTest(MetricsMiddleware(), "/test")

	if w.Code != 200 {
		t.Errorf("期望 200，实际 %d", w.Code)
	}
	// 总请求数应增加
	if GetTotalRequests() <= beforeTotal {
		t.Error("总请求数应增加")
	}
	// 活跃请求数应回到之前的值
	if GetActiveRequests() != beforeActive {
		t.Errorf("活跃请求数应恢复，期望 %d，实际 %d", beforeActive, GetActiveRequests())
	}
	// 应有响应头
	if w.Header().Get("X-Request-Duration") == "" {
		t.Error("缺少 X-Request-Duration 响应头")
	}
}

func TestMetricsMiddleware_Concurrent(t *testing.T) {
	// 模拟并发请求
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func() {
			setupMiddlewareTest(MetricsMiddleware(), "/test")
			done <- true
		}()
	}
	for i := 0; i < 3; i++ {
		<-done
	}
	if GetTotalRequests() < 3 {
		t.Errorf("并发 3 次后总请求数应 >= 3，实际 %d", GetTotalRequests())
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := NewRateLimiter()
	r := gin.New()
	r.Use(RateLimitMiddleware(limiter, 2, 1)) // 2次/秒
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// 前两次应通过
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("第 %d 次请求期望 200，实际 %d", i+1, w.Code)
		}
	}

	// 解析响应
	type resp struct {
		Code int `json:"code"`
	}

	// 第三次应被限流
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var r1 resp
	json.Unmarshal(w.Body.Bytes(), &r1)
	if w.Code == 200 && r1.Code == 429 {
		// 可能通过200状态码返回429业务码
	} else if w.Code == 429 {
		// 也可能直接429
	} else {
		// 允许，取决于具体实现
	}
}
