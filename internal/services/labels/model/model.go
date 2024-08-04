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
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
}

type BookLabelSwagger struct {
	LabelName string `json:"label_name" bson:"labelName"`
}

type BookLabel struct {
	BookId    primitive.ObjectID `json:"book_id" bson:"book"`
	LabelName string             `json:"label_name" bson:"labelName"`
}

type AddBookLabelResponse struct {
	ID      primitive.ObjectID `json:"id" bson:"_id"`
	LabelId primitive.ObjectID `json:"label_id" bson:"labelId"`
	BookId  primitive.ObjectID `json:"book_id" bson:"bookId"`
}

type BookResponse struct {
	ID         string    `json:"id" bson:"_id"`
	Title      string    `json:"title" bson:"title"`
	CreatedAt  time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updatedAt"`
	IsArchived bool      `json:"is_archived" bson:"isArchived"`
}
