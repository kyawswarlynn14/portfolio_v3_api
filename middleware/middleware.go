package middleware

import (
	"net/http"
	"strings"

	token "portfolio/tokens"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header not provided"})
			c.Abort()
			return
		}

		// Check if the token is a Bearer token
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Extract token
		clientToken := fields[1]
		claims, err := token.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			c.Abort()
			return
		}

		// Set user information (e.g., email) in the request context
		c.Set("email", claims.Email)
		c.Set("userId", claims.User_ID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func Authorization(roles []int) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleFromMdw, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Role not found in request context"})
			c.Abort()
			return
		}

		userRole, ok := userRoleFromMdw.(int)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Invalid role format"})
			c.Abort()
			return
		}

		isAuthorized := false
		for _, role := range roles {
			if userRole == role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
