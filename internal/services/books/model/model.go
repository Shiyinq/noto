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
	UserId     primitive.ObjectID `json:"user_id" bson:"userId"`
	Title      string             `json:"title" bson:"title"`
	CreatedAt  time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updatedAt"`
	IsArchived bool               `json:"is_archived" bson:"isArchived"`
}

type BookUpdateSwagger struct {
	Title string `json:"title" bson:"title"`
}

type BookUpdate struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"user_id" bson:"userId"`
	Title  string             `json:"title" bson:"title"`
}

type Label struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

type BookResponse struct {
	ID         string    `json:"id" bson:"_id"`
	Title      string    `json:"title" bson:"title"`
	CreatedAt  time.Time `json:"created_at" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updatedAt"`
	IsArchived bool      `json:"is_archived" bson:"isArchived"`
	Labels     []Label   `json:"labels" bson:"labels"`
}

type PaginationMetadata struct {
	TotalData    int  `json:"totalData" bson:"totalData"`
	TotalPage    int  `json:"totalPage" bson:"totalPage"`
	PreviousPage *int `json:"previousPage" bson:"previousPage"`
	CurrentPage  int  `json:"currentPage" bson:"currentPage"`
	NextPage     *int `json:"nextPage" bson:"nextPage"`
}

type PaginatedBookResponse struct {
	Metadata PaginationMetadata `json:"metadata" bson:"metadata"`
	Data     []BookResponse     `json:"data" bson:"data"`
}

type ArchiveBookSwagger struct {
	IsArchived bool `json:"is_archived" bson:"isArchived"`
}

type ArchiveBook struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId     primitive.ObjectID `json:"user_id" bson:"userId"`
	IsArchived bool               `json:"is_archived" bson:"isArchived"`
}
