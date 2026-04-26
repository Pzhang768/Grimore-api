package middleware

import (
	"encoding/json"
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

func decodeBody(t *testing.T, body string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		t.Fatalf("response body is not valid JSON: %v — body: %s", err, body)
	}
	return m
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
	body := decodeBody(t, w.Body.String())
	if body["error"] != "invalid input" {
		t.Errorf("expected error message 'invalid input', got %v", body["error"])
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
		t.Errorf("expected 400 default for public error, got %d", w.Code)
	}
	body := decodeBody(t, w.Body.String())
	if body["error"] == nil {
		t.Error("expected non-empty error field in response body")
	}
}

// TestErrorHandler_InternalError verifies 500 status, generic message, and no detail leakage.
func TestErrorHandler_InternalError(t *testing.T) {
	r := newRouter()
	r.GET("/", func(c *gin.Context) {
		_ = c.Error(&gin.Error{Err: errors.New("secret db password is hunter2"), Type: gin.ErrorTypePrivate})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	body := decodeBody(t, w.Body.String())
	if body["error"] != "internal server error" {
		t.Errorf("expected generic message, got %v — internal detail may have leaked", body["error"])
	}
}

// TestErrorHandler_MultiError_LastErrorWins documents that the last attached error drives the response.
// Private first, public last → public wins (responds with public message).
// Public first, private last → private wins (responds with 500).
func TestErrorHandler_MultiError_LastErrorWins(t *testing.T) {
	cases := []struct {
		name           string
		attachErrors   func(c *gin.Context)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "private then public — last is public",
			attachErrors: func(c *gin.Context) {
				_ = c.Error(&gin.Error{Err: errors.New("db error"), Type: gin.ErrorTypePrivate})
				_ = c.Error(&gin.Error{Err: errors.New("user error"), Type: gin.ErrorTypePublic})
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "user error",
		},
		{
			name: "public then private — last is private",
			attachErrors: func(c *gin.Context) {
				_ = c.Error(&gin.Error{Err: errors.New("user error"), Type: gin.ErrorTypePublic})
				_ = c.Error(&gin.Error{Err: errors.New("db error"), Type: gin.ErrorTypePrivate})
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := newRouter()
			r.GET("/", func(c *gin.Context) { tc.attachErrors(c) })

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected %d, got %d", tc.expectedStatus, w.Code)
			}
			body := decodeBody(t, w.Body.String())
			if body["error"] != tc.expectedBody {
				t.Errorf("expected %q, got %v", tc.expectedBody, body["error"])
			}
		})
	}
}
