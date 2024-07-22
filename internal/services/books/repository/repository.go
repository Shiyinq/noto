package repository

import (
	"context"
	"errors"
	"time"

	"noto/internal/config"
	"noto/internal/services/books/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookRepository interface {
	CreateBook(book *model.Book) (*model.Book, error)
	GetAllBooks() ([]model.BookResponse, error)
	GetBookByID(id string) (*model.BookResponse, error)
}

type BookRepositoryImpl struct {
	books *mongo.Collection
}

func NewBookRepository() BookRepository {
	return &BookRepositoryImpl{books: config.DB.Collection("books")}
}

func (r *BookRepositoryImpl) CreateBook(book *model.Book) (*model.Book, error) {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	_, err := r.books.InsertOne(context.Background(), book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepositoryImpl) GetAllBooks() ([]model.BookResponse, error) {
	var books []model.BookResponse
	cursor, err := r.books.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &books); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepositoryImpl) GetBookByID(id string) (*model.BookResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	var book model.BookResponse
	err = r.books.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return &book, nil
}
