package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"task-manager/domain"
	"task-manager/errs"
	"task-manager/usecases"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// --- MongoDB Implementation ---

type mongoUserRepository struct {
	collection *mongo.Collection
}

// mongoUser is a private struct to handle BSON mapping.
type mongoUser struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username"`
	PasswordHash string             `bson:"password_hash"`
	Role         string             `bson:"role"`
}

func NewMongoUserRepository(collection *mongo.Collection) usecases.UserRepository {
	return &mongoUserRepository{collection: collection}
}

// func (r *mongoUserRepository) Generate() string {
// 	return primitive.NewObjectID().Hex()
// }

func (r *mongoUserRepository) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mUser := mongoUser{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}
	res, err := r.collection.InsertOne(ctx, mUser)
	if err != nil {
		return err
	}
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return nil
}

func (r *mongoUserRepository) GetByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var mUser mongoUser
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&mUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Printf("WARN: Login failed for username '%s': user not found", username)
			return nil, errs.ErrUserNotFound
		}
		log.Printf("ERROR: Database error during login for username '%s': %v", username, err)
		return nil, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	user := domain.User{
		ID:           mUser.ID.Hex(),
		Username:     mUser.Username,
		Password:     "",
		PasswordHash: mUser.PasswordHash,
		Role:         mUser.Role,
	}
	return &user, nil
}

func (r *mongoUserRepository) GetByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user domain.User
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errs.ErrInvalidUserId
		}
		return nil, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	return &user, nil
}

func (r *mongoUserRepository) UpdateUserStatus(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Attempting to update user status in repository ")
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}

	update := bson.M{
		"$set": bson.M{
			"role": domain.RoleAdmin,
		},
	}
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objId}, update)
	if err != nil {
		return fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	if result.MatchedCount == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *mongoUserRepository) Count() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	return count, nil
}

func (r *mongoUserRepository) CheckUsername(username string) (exist bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	count, err := r.collection.CountDocuments(ctx, bson.M{"username": username})
	if err != nil {
		return false, fmt.Errorf("%w: %v", errs.ErrUnexpected, err)
	}
	return count > 0, nil
}
