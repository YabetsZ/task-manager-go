package infrastructure

import (
	"log"
	"net/http"
	"task-manager/domain"
	"task-manager/errs"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTServiceV5 struct{}

func NewJWTServiceV5() *JWTServiceV5 {
	return &JWTServiceV5{}
}

const JWTSecret = "task_manager_secret"

type CustomClaims struct {
	UserID   string // `json:"user_id"`
	Username string // `json:"username"`
	Role     string // `json:"role"`
	jwt.RegisteredClaims
}

func (js *JWTServiceV5) GenerateJWT(user *domain.User) (string, error) {
	claims := CustomClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		log.Printf("ERROR: Failed to generate JWT for user '%s': %v", user.Username, err)
		return "", errs.New(http.StatusInternalServerError, "unexpected error", err)
	}
	return signedToken, nil
}
