package middleware

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"community-server/DB/redis"
	"community-server/pkg/response"

	"github.com/gin-gonic/gin"
	redisclient "github.com/redis/go-redis/v9"
)

type RateLimiter struct{}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

func (rl *RateLimiter) getClient() *redisclient.Client {
	return redis.GetClient()
}

func (rl *RateLimiter) Allow(key string, capacity int64, refillRate float64) (bool, error) {
	client := rl.getClient()
	if client == nil {
		return true, nil
	}

	ctx := context.Background()
	tokensKey := fmt.Sprintf("%s:tokens", key)
	lastRefillKey := fmt.Sprintf("%s:last_refill", key)

	tokensStr, err := client.Get(ctx, tokensKey).Result()
	if err != nil && err != redisclient.Nil {
		return false, err
	}

	lastRefillStr, err := client.Get(ctx, lastRefillKey).Result()
	if err != nil && err != redisclient.Nil {
		return false, err
	}

	var tokens float64
	var lastRefill int64
	now := time.Now().Unix()

	if err == redisclient.Nil || tokensStr == "" {
		tokens = float64(capacity)
		lastRefill = now
	} else {
		tokens, _ = strconv.ParseFloat(tokensStr, 64)
		if lastRefillStr == "" {
			lastRefill = now
		} else {
			lastRefill, _ = strconv.ParseInt(lastRefillStr, 10, 64)
		}
	}

	elapsed := float64(now - lastRefill)
	added := elapsed * refillRate
	tokens = math.Min(float64(capacity), tokens+added)

	if tokens >= 1 {
		tokens--
		pipe := client.Pipeline()
		pipe.Set(ctx, tokensKey, tokens, time.Hour)
		pipe.Set(ctx, lastRefillKey, now, time.Hour)
		_, err := pipe.Exec(ctx)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	pipe := client.Pipeline()
	pipe.Set(ctx, tokensKey, tokens, time.Hour)
	pipe.Set(ctx, lastRefillKey, now, time.Hour)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return false, err
	}
	return false, nil
}

func RateLimitMiddleware(limiter *RateLimiter, capacity int64, refillRate float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter.getClient() == nil {
			c.Next()
			return
		}

		key := fmt.Sprintf("rate_limit:%s:%s", c.ClientIP(), c.Request.URL.Path)

		allowed, err := limiter.Allow(key, capacity, refillRate)
		if err != nil {
			response.ErrorWithMsg(c, response.CodeServerBusy, "服务器繁忙")
			c.Abort()
			return
		}

		if !allowed {
			c.Header("X-RateLimit-Limit", strconv.FormatInt(capacity, 10))
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", "1")
			response.ErrorWithMsg(c, response.CodeTooManyRequests, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
