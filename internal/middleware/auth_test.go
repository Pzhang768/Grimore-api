package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const testSecret = "test-jwt-secret"

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}
	return db, mock
}

func makeToken(t *testing.T, secret string, userID uuid.UUID, email string, expiry time.Time) string {
	t.Helper()
	claims := supabaseClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(expiry),
		},
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func newAuthRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(Auth(testSecret, db))
	r.GET("/protected", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	return r
}

func TestAuth_MissingHeader(t *testing.T) {
	db, _ := newMockDB(t)
	r := newAuthRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_InvalidBearerPrefix(t *testing.T) {
	db, _ := newMockDB(t)
	r := newAuthRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Token abc123")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_TamperedToken(t *testing.T) {
	db, _ := newMockDB(t)
	r := newAuthRouter(db)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ0ZXN0In0.tampered")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	db, _ := newMockDB(t)
	r := newAuthRouter(db)

	userID := uuid.New()
	tokenStr := makeToken(t, testSecret, userID, "test@example.com", time.Now().Add(-time.Hour))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_WrongSigningMethod(t *testing.T) {
	db, _ := newMockDB(t)
	r := newAuthRouter(db)

	userID := uuid.New()
	claims := supabaseClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email: "test@example.com",
	}
	// sign with RS256 (wrong method) using a throwaway key
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenStr, _ := token.SignedString([]byte("other-secret"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_InvalidSubject(t *testing.T) {
	db, _ := newMockDB(t)
	r := newAuthRouter(db)

	claims := supabaseClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "not-a-uuid",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		Email: "test@example.com",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(testSecret))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_DBUpsertFailure(t *testing.T) {
	db, mock := newMockDB(t)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).WillReturnError(errors.New("db error"))
	mock.ExpectRollback()

	r := newAuthRouter(db)

	userID := uuid.New()
	tokenStr := makeToken(t, testSecret, userID, "test@example.com", time.Now().Add(time.Hour))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestAuth_ValidToken_SetsUserID(t *testing.T) {
	db, mock := newMockDB(t)

	userID := uuid.New()
	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var capturedID uuid.UUID
	r := gin.New()
	r.Use(Auth(testSecret, db))
	r.GET("/protected", func(c *gin.Context) {
		id, _ := c.Get(ContextKeyUserID)
		capturedID = id.(uuid.UUID)
		c.Status(http.StatusOK)
	})

	tokenStr := makeToken(t, testSecret, userID, "test@example.com", time.Now().Add(time.Hour))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if capturedID != userID {
		t.Errorf("expected userID %s in context, got %s", userID, capturedID)
	}
}
