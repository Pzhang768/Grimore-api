package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		for _, e := range c.Errors {
			if e.Type != gin.ErrorTypePublic {
				slog.Error("internal error", "error", e.Error())
			}
		}

		err := c.Errors.Last()
		switch err.Type {
		case gin.ErrorTypePublic:
			status := c.Writer.Status()
			if status == http.StatusOK {
				status = http.StatusBadRequest
			}
			c.JSON(status, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}
}
