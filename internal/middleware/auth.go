package middleware

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/Pzhang768/Grimore-api/internal/models"
)

const ContextKeyUserID = "userID"

type supabaseClaims struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

// decodeBase64URL decodes a base64url-encoded big integer (used for EC key coordinates).
func decodeBase64URL(s string) (*big.Int, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return new(big.Int).SetBytes(b), nil
}

func Auth(jwtSecret string, db *gorm.DB) gin.HandlerFunc {
	hmacSecret := []byte(jwtSecret)

	// P-256 public key from Supabase JWKS — used for ES256 tokens.
	// x and y are the base64url-encoded coordinates from /auth/v1/.well-known/jwks.json.
	const jwksX = "r3UIJ8pYajMUDPh1DgvvoRO77-654ECdy4Gzm0qVGuc"
	const jwksY = "EB6ZVt3gJx9WJhoMSqDPqyBU9a6yTU3JS_q4yo8qazg"

	x, _ := decodeBase64URL(jwksX)
	y, _ := decodeBase64URL(jwksY)
	ecKey := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")

		claims := &supabaseClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
			switch t.Method.(type) {
			case *jwt.SigningMethodHMAC:
				return hmacSecret, nil
			case *jwt.SigningMethodECDSA:
				return ecKey, nil
			default:
				return nil, errors.New("unexpected signing method")
			}
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token subject"})
			return
		}

		user := models.User{ID: userID, Email: claims.Email}
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"email"}),
		}).Create(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		c.Set(ContextKeyUserID, userID)
		c.Next()
	}
}
