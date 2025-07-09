package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/SoliMark/gotasker-pro/internal/constant"
	"github.com/SoliMark/gotasker-pro/internal/util"
)

func TestJWTAuthMiddleware(t *testing.T) {
	// 1) 初始化 JWTMaker
	secretKey := "test_secret_key"
	jwtMaker := util.NewJWTMaker(secretKey)

	// 2) 產生一個有效的 token
	token, err := jwtMaker.GenerateToken(123, time.Minute)
	require.NoError(t, err)

	// 3) 建立一個帶 Middleware 的測試路由
	router := gin.New()
	router.GET("/protected",
		JWTAuthMiddleware(jwtMaker),
		func(c *gin.Context) {
			// 從 context 拿 user_id
			userID, exists := c.Get(constant.ContextUserIDKey)
			require.True(t, exists)

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
			})
		},
	)

	t.Run("valid token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set(constant.HeaderAuthorization, "Bearer "+token)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), `"user_id":123`)
	})

	t.Run("missing token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set(constant.HeaderAuthorization, token) // 缺少 Bearer 前綴

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set(constant.HeaderAuthorization, "Bearer invalid.token.here")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
