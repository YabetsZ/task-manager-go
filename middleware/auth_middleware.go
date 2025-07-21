package middleware

import (
	"errors"
	"net/http"
	"task-manager/data"
	"task-manager/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a gin.HandlerFunc for JWT authentication and authorization.
func AuthMiddleware(userService *data.UserService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}

		tokenString := authHeader[7:]

		claims := data.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
			return []byte(data.JWTSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			} else {
				// log.Panic(err.Error())
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			return
		}

		if !token.Valid {
			// log.Panic("just token being invalid after passing parsing")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		user, appErr := userService.GetUserByID(claims.UserID)

		if appErr != nil {
			c.AbortWithStatusJSON(appErr.Code, gin.H{"error": appErr.Msg})
			return
		}
		// Guests are not allowed to pass
		if requiredRole == models.RoleUser && user.Role != models.RoleUser && user.Role != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		// Only admins can access admin routes
		if requiredRole == models.RoleAdmin && user.Role != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}

		// Set user in context for downstream handlers
		c.Set("user", user)
		c.Next()
	}
}
