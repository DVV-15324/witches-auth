package middleware

import (
	"time"
	"github.com/DVV-15324/witches/pkg/core/response/logger"
	w_utils "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	wcmd_utils "github.com/DVV-15324/witches/cmd/utils"
)

func TimingMiddleware(log *logger.ModelLogger, cfg *wcmd_utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// Lấy tid từ RequestContext
		tid := w_utils.GetTid(c.Request.Context(), cfg)
		sub := w_utils.GetSub(c.Request.Context(), cfg)

		ms := time.Since(start).Milliseconds()

		if ms > cfg.SlowThreshold {
			log.WarnWithFields("Slow request",
				zap.String("trace_id(tid)", tid),
				zap.String("subject(sub)", sub),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.Int64("duration_ms", ms),
			)
		} else {
			log.InfoWithFields("Request completed",
				zap.String("trace_id(tid)", tid),
				zap.String("subject(sub)", sub),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status", c.Writer.Status()),
				zap.Int64("duration_ms", ms),
			)
		}

		c.Header("X-Response-Time", time.Since(start).String())
	}
}
