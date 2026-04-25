package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		fields := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"duration", time.Since(start),
		}

		if len(c.Errors) > 0 {
			slog.Error("request", append(fields, "errors", c.Errors.String())...)
		} else {
			slog.Info("request", fields...)
		}
	}
}
