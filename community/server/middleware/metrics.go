package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		_ = start

		c.Next()

		// TODO: 实现指标记录逻辑
		// duration := time.Since(start)
		// hasError := len(c.Errors) > 0
		// metrics.RecordRequest(duration, hasError)
	}
}
