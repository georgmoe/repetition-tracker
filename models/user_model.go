package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name,omitempty"`
	Email     string             `bson:"email" json:"email"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Password  string             `bson:"password" json:"password"`
}
