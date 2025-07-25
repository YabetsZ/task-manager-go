package infrastructure

import (
	"errors"
	"net/http"
	"task-manager/domain"
	"task-manager/usecases"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a gin.HandlerFunc for JWT authentication and authorization.
func AuthMiddleware(userUsecase usecases.UserUsecase, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}

		tokenString := authHeader[7:]

		claims := CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
			return []byte(JWTSecret), nil
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

		// look out this error
		user, err := userUsecase.GetUserByID(claims.UserID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		// Admins can access user routes, but not the other way around
		if requiredRole == domain.RoleUser && user.Role != domain.RoleUser && user.Role != domain.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		// Only admins can access admin routes
		if requiredRole == domain.RoleAdmin && user.Role != domain.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions, admin access required"})
			return
		}

		// Set user in context for downstream handlers
		c.Set("user", user)
		c.Next()
	}
}
