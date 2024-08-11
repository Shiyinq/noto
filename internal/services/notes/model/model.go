package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteCreateSwagger struct {
	Text string `json:"text" bson:"text"`
}

type NoteCreate struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId    primitive.ObjectID `json:"user_id" bson:"userId"`
	BookId    primitive.ObjectID `json:"book_id" bson:"bookId"`
	Text      string             `json:"text" bson:"text"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
}

type NoteUpdateSwagger struct {
	Text string `json:"text" bson:"text"`
}

type NoteUpdate struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId primitive.ObjectID `json:"user_id" bson:"userId"`
	BookId primitive.ObjectID `json:"book_id" bson:"bookId"`
	Text   string             `json:"text" bson:"text"`
}

type NoteResponse struct {
	ID        string             `json:"id" bson:"_id"`
	BookId    primitive.ObjectID `json:"book_id" bson:"bookId"`
	Text      string             `json:"text" bson:"text"`
	CreatedAt time.Time          `json:"created_at" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updatedAt"`
}

type PaginationMetadata struct {
	TotalData    int  `json:"totalData" bson:"totalData"`
	TotalPage    int  `json:"totalPage" bson:"totalPage"`
	PreviousPage *int `json:"previousPage" bson:"previousPage"`
	CurrentPage  int  `json:"currentPage" bson:"currentPage"`
	NextPage     *int `json:"nextPage" bson:"nextPage"`
}

type PaginatedNoteResponse struct {
	Metadata PaginationMetadata `json:"metadata" bson:"metadata"`
	Data     []NoteResponse     `json:"data" bson:"data"`
}
