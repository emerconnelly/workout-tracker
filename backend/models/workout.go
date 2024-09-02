package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Workout struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name"          bson:"name"`
	Exercises []Exercise         `json:"exercises"     bson:"exercises"`
}
