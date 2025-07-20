package data

import (
	"context"
	"errors"
	"net/http"
	"task-manager/errs"
	"task-manager/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

const jwtSecret = "task_manager_secret"

type CustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type UserService struct {
	collection *mongo.Collection
}

// Constructor for UserService
func NewUserService(collection *mongo.Collection) *UserService {
	return &UserService{collection}
}

func (us *UserService) RegisterUser(user *models.User) *errs.AppError {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return errs.New(http.StatusInternalServerError, "Unexpected error", err)
	}
	user.PasswordHash = string(hashedPassword)

	// For the requirement: "If the database is empty, the first created user will be an admin."
	count, err := us.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return errs.New(http.StatusInternalServerError, "Unexpected error", err)
	}
	if count == 0 {
		user.Role = models.RoleAdmin
	} else {
		user.Role = models.RoleUser
	}

	// Generate a new ID
	user.ID = primitive.NewObjectID()

	us.collection.InsertOne(ctx, user)
	return nil
}

func (us *UserService) Login(username, password string) (string, *errs.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := us.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errs.New(http.StatusNotFound, "invalid credentials", err)
		}
		return "", errs.New(http.StatusInternalServerError, "Unexpected Error", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errs.New(http.StatusInternalServerError, "invalid credentials", err)
	}

	return us.generateJWT(&user)
}

func (us *UserService) generateJWT(user *models.User) (string, *errs.AppError) {
	claims := CustomClaims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errs.New(http.StatusInternalServerError, "unexpected error", err)
	}
	return signedToken, nil
}
