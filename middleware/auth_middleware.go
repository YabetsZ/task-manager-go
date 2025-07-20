package middleware

import (
	"errors"
	"net/http"
	"task-manager/data"
	"task-manager/errs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware creates a gin.HandlerFunc for JWT authentication and authorization.
func AuthMiddleware(userService *data.UserService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorizaiton")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}

		tokenString := authHeader[7:]
		claims := data.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			return data.JWTSecret, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token has expired"})
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			}
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		user, err := userService.GetUserByID(claims.UserID)

		if err != nil {
			appErr := &errs.AppError{}
			if errors.As(err, appErr) {
				c.AbortWithStatusJSON(appErr.Code, gin.H{"error": appErr.Msg})
				return
			}
		}

	}
}
