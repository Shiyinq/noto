package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookCreateSwagger struct {
	Title string `json:"title" bson:"title"`
}

type BookCreate struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title      string             `json:"title" bson:"title"`
	CreatedAt  time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updatedAt"`
	IsArchived bool               `json:"is_archived" bson:"isArchived"`
}

type BookResponse struct {
	ID         string    `json:"id" bson:"_id"`
	Title      string    `json:"title" bson:"title"`
	CreatedAt  time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updatedAt"`
	IsArchived bool      `json:"is_archived" bson:"isArchived"`
}

type ArchiveBook struct {
	IsArchived bool `json:"is_archived" bson:"isArchived"`
}
