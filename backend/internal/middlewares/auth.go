package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/saimonpoliveira/busato_tasks/backend/internal/utils"
)

const (
	ContextUserIDKey = "userID"
	ContextEmailKey  = "email"
	ContextRoleKey   = "role"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			utils.Unauthorized(c, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextEmailKey, claims.Email)
		c.Set(ContextRoleKey, claims.Role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(ContextRoleKey)
		if !exists || role != "admin" {
			utils.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
