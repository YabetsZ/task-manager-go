package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Role         string             `json:"role" bson:"role"`
}
