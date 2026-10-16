package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

type LimitMiddleWare struct{}

// createRateLimitMiddleware tạo middleware rate limit với headers
func (limit *LimitMiddleWare) CreateRateLimitMiddleware(rateLimit *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy key từ IP (hoặc custom)
		key := c.ClientIP()

		// Lấy thông tin rate limit từ Redis
		ctx := c.Request.Context()
		res, err := rateLimit.Get(ctx, key)
		if err != nil {
			// Nếu lỗi Redis, vẫn cho request đi qua (fallback)
			c.Next()
			return
		}

		// THÊM HEADERS
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", res.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", res.Remaining))
		resetTime := time.Unix(res.Reset, 0)
		c.Header("X-RateLimit-Reset", resetTime.Format(time.RFC3339))

		// Thêm header thời gian còn lại (seconds)
		retryAfter := time.Until(resetTime).Seconds()
		if retryAfter > 0 {
			c.Header("Retry-After", fmt.Sprintf("%.0f", retryAfter))
		}
		// Nếu đã đạt giới hạn
		if res.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"limit":       res.Limit,
				"remaining":   res.Remaining,
				"reset":       res.Reset,
				"retry_after": retryAfter,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
