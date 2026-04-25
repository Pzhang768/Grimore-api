package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newRouter() *gin.Engine {
	r := gin.New()
	r.Use(ErrorHandler())
	return r
}

func TestErrorHandler_NoErrors(t *testing.T) {
	r := newRouter()
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if w.Body.String() != "" {
		t.Errorf("expected empty body, got %s", w.Body.String())
	}
}

func TestErrorHandler_PublicError_WithExplicitStatus(t *testing.T) {
	r := newRouter()
	r.GET("/", func(c *gin.Context) {
		c.Status(http.StatusUnprocessableEntity)
		_ = c.Error(&gin.Error{Err: errors.New("invalid input"), Type: gin.ErrorTypePublic})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
	if w.Body.String() == "" {
		t.Error("expected error body, got empty")
	}
}

func TestErrorHandler_PublicError_DefaultsTo400(t *testing.T) {
	r := newRouter()
	r.GET("/", func(c *gin.Context) {
		_ = c.Error(&gin.Error{Err: errors.New("bad request"), Type: gin.ErrorTypePublic})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 default for public error with 200 status, got %d", w.Code)
	}
}

func TestErrorHandler_InternalError_Returns500(t *testing.T) {
	r := newRouter()
	r.GET("/", func(c *gin.Context) {
		_ = c.Error(&gin.Error{Err: errors.New("db exploded"), Type: gin.ErrorTypePrivate})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	if w.Body.String() != `{"error":"internal server error"}` {
		t.Errorf("expected generic error body, got %s", w.Body.String())
	}
}

func TestErrorHandler_InternalError_DoesNotLeakDetails(t *testing.T) {
	r := newRouter()
	r.GET("/", func(c *gin.Context) {
		_ = c.Error(&gin.Error{Err: errors.New("secret db password is hunter2"), Type: gin.ErrorTypePrivate})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Body.String() != `{"error":"internal server error"}` {
		t.Errorf("internal error detail leaked: %s", w.Body.String())
	}
}
