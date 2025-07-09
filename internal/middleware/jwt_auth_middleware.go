package middleware

import (
	"net/http"
	"strings"

	"github.com/SoliMark/gotasker-pro/internal/constant"
	"github.com/SoliMark/gotasker-pro/internal/util"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtMaker *util.JWTMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(constant.HeaderAuthorization)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				constant.ErrorKey: "authorization header is missing",
			})
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				constant.ErrorKey: "invalid authorization header format",
			})
			return
		}

		tokenStr := parts[1]

		claims, err := jwtMaker.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				constant.ErrorKey: "invaild or expired token",
			})
			return
		}

		c.Set(constant.ContextUserIDKey, claims.UserID)
		c.Next()
	}
}
