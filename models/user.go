package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Username     string             `json:"username" bson:"username"`
	Password     string             `json:"password,omitempty"`      // From frontend
	PasswordHash string             `bson:"password_hash,omitempty"` // Stored in DB
	Role         string             `json:"-" bson:"role"`
}
