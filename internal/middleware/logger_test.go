package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogger_SuccessRequest(t *testing.T) {
	r := gin.New()
	r.Use(Logger())
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestLogger_ErrorRequest(t *testing.T) {
	r := gin.New()
	r.Use(Logger())
	r.GET("/", func(c *gin.Context) {
		_ = c.Error(&gin.Error{Err: errors.New("something failed"), Type: gin.ErrorTypePrivate})
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
