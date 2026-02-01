package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ya-breeze/kin-core/auth"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := []byte("test-secret")

	t.Run("Valid Token", func(t *testing.T) {
		token, _ := auth.GenerateToken(1, 10, secret, time.Hour)

		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {
			userID := c.MustGet("user_id").(uint)
			familyID := c.MustGet("family_id").(uint)
			if userID != 1 || familyID != 10 {
				t.Errorf("Context values mismatch: user=%d, family=%d", userID, familyID)
			}
			c.Status(http.StatusOK)
		})

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, c.Request)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Missing Header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {})

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, c.Request)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", w.Code)
		}
	})

	t.Run("Invalid Token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.Use(AuthMiddleware(secret))
		r.GET("/test", func(c *gin.Context) {})

		c.Request, _ = http.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Authorization", "Bearer invalid-token")
		r.ServeHTTP(w, c.Request)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected 401, got %d", w.Code)
		}
	})
}
