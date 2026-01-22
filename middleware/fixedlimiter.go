package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type FixedWindowLimiter struct {
	mu          sync.Mutex
	ipCounters  map[string]int       // 每个IP的请求计数
	windowStart map[string]time.Time // 每个IP的窗口开始时间
	limit       int                  // 窗口内最大请求数
	window      time.Duration        // 窗口时长（如1分钟）
}

func NewFixedWindowLimiter() *FixedWindowLimiter {
	return &FixedWindowLimiter{
		ipCounters:  make(map[string]int),
		windowStart: make(map[string]time.Time),
		limit:       100,
		window:      time.Minute,
	}
}

// Gin中间件：IP限流
func (lim *FixedWindowLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP() // 获取客户端IP
		now := time.Now()

		lim.mu.Lock()
		defer lim.mu.Unlock()

		// 检查窗口是否过期，过期则重置计数
		if start, exists := lim.windowStart[ip]; !exists || now.Sub(start) > lim.window {
			lim.windowStart[ip] = now
			lim.ipCounters[ip] = 1
		} else {
			// 窗口内请求数超过阈值，拒绝
			if lim.ipCounters[ip] >= lim.limit {
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "请求过于频繁，请稍后再试"})
				c.Abort()
				return
			}
			lim.ipCounters[ip]++
		}

		c.Next()
	}
}
