package middleware

import (
	"net/http"
	"strings"

	"github.com/Saurav-hehe/Complaint-Management/utils"
	"github.com/gin-gonic/gin"
)

func StaffAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// Only allow staff roles
		role, ok := claims["role"].(string)
		if !ok || (role != "electrician" && role != "plumber" && role != "carpenter") {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
		c.Set("email", claims["email"])
		c.Set("role", role)
		c.Next()
	}
}
