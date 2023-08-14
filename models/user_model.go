package models

import "time"

type User struct {
	Name      string    `bson:"name" json:"name,omitempty"`
	Email     string    `bson:"email" json:"email"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	Password  string    `bson:"password" json:"password"`
}
