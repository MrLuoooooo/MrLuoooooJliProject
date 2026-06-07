package middleware

import (
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	totalRequests  int64
	activeRequests int64
)

// MetricsMiddleware 记录请求计数和耗时
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		atomic.AddInt64(&totalRequests, 1)
		atomic.AddInt64(&activeRequests, 1)

		c.Next()

		atomic.AddInt64(&activeRequests, -1)
		duration := time.Since(start)

		// 记录到 gin 默认日志（可扩展为 Prometheus）
		c.Header("X-Request-Duration", duration.String())
	}
}

// GetTotalRequests 获取总请求数
func GetTotalRequests() int64 {
	return atomic.LoadInt64(&totalRequests)
}

// GetActiveRequests 获取当前活跃请求数
func GetActiveRequests() int64 {
	return atomic.LoadInt64(&activeRequests)
}
