package middleware

import (
	"sync"
	"time"

	"property-hub-backend/config"
	"property-hub-backend/utils"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu       sync.Mutex
	attempts map[string][]time.Time
}

var loginLimiter = &rateLimiter{attempts: make(map[string][]time.Time)}
var globalLimiter = &rateLimiter{attempts: make(map[string][]time.Time)}

func (rl *rateLimiter) allow(key string, maxAttempts int, window time.Duration) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	// Clean old entries
	var recent []time.Time
	for _, t := range rl.attempts[key] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= maxAttempts {
		rl.attempts[key] = recent
		return false
	}

	recent = append(recent, now)
	rl.attempts[key] = recent
	return true
}

// RateLimitLogin limits login attempts per IP
func RateLimitLogin(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !loginLimiter.allow(ip, cfg.RateLimitLoginPerMinute, time.Minute) {
			utils.TooManyRequests(c, "RATE_LOGIN_LIMIT", "Terlalu banyak percobaan login. Silakan coba lagi dalam 1 menit.")
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitGlobal limits all requests per IP
func RateLimitGlobal(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !globalLimiter.allow(ip, cfg.RateLimitGlobalPerMinute, time.Minute) {
			utils.TooManyRequests(c, "RATE_GLOBAL_LIMIT", "Terlalu banyak permintaan. Silakan coba lagi nanti.")
			c.Abort()
			return
		}
		c.Next()
	}
}
