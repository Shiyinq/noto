package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LabelCreateSwagger struct {
	Name string `json:"name" bson:"name"`
}

type LabelCreate struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
}

type LabelUpdate struct {
	Name string `json:"name" bson:"name"`
}

type LabelResponse struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time `json:"updated_at" bson:"updatedAt"`
}
