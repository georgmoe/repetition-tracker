package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Set struct {
	Variation string    `bson:"variation,omitempty" json:"variation"` // Use of helpers like bands, or push up bars...
	Weight    float32   `bson:"weight,omitempty" json:"weight"`       // If exercise is done with external weights
	Reps      int       `bson:"reps" json:"reps"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type Exercise struct {
	Name     string `bson:"name" json:"name"`
	Category string `bson:"category,omitempty" json:"category"` // PUSH, PULL, LEG, CORE, ...
	Sets     []Set  `bson:"sets" json:"sets"`
}

type Workout struct {
	Name      string             `bson:"name, omitempty" json:"name"`
	Date      time.Time          `bson:"date" json:"date"`
	Exercises []Exercise         `bson:"exercises" json:"exercises"`
	UserId    primitive.ObjectID `bson:"userId" json:"-"`
}
