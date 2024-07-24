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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookRepository interface {
	CreateBook(book *model.Book) (*model.Book, error)
	GetAllBooks(isArchived bool) ([]model.BookResponse, error)
	GetBookByID(id string) (*model.BookResponse, error)
	UpdateBook(id string, title string) (*model.BookResponse, error)
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
	book.IsArchived = book.IsArchived || false

	_, err := r.books.InsertOne(context.Background(), book)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepositoryImpl) GetAllBooks(isArchived bool) ([]model.BookResponse, error) {
	var books []model.BookResponse
	filter := bson.M{"isArchived": isArchived}
	cursor, err := r.books.Find(context.Background(), filter)
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

func (r *BookRepositoryImpl) UpdateBook(id string, title string) (*model.BookResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid ID format")
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"title":     title,
			"updatedAt": time.Now(),
		},
	}

	result := r.books.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("book not found")
		}
		return nil, result.Err()
	}

	var updatedBook model.BookResponse
	err = result.Decode(&updatedBook)
	if err != nil {
		return nil, err
	}
	return &updatedBook, nil
}
