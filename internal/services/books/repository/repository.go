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
	CreateBook(book *model.BookCreate) (*model.BookCreate, error)
	GetBooks(userId primitive.ObjectID, isArchived bool) ([]model.BookResponse, error)
	GetBook(userId primitive.ObjectID, bookId primitive.ObjectID) (*model.BookResponse, error)
	UpdateBook(book *model.BookUpdate) (*model.BookResponse, error)
	ArchiveBook(book *model.ArchiveBook) (*model.BookResponse, error)
}

type BookRepositoryImpl struct {
	books *mongo.Collection
}

func NewBookRepository() BookRepository {
	return &BookRepositoryImpl{books: config.DB.Collection("books")}
}

func bookAgregate(matchCondition bson.D) mongo.Pipeline {
	return mongo.Pipeline{
		{
			{Key: "$match", Value: matchCondition},
		},
		{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "book_labels"},
					{Key: "localField", Value: "_id"},
					{Key: "foreignField", Value: "bookId"},
					{Key: "as", Value: "book_labels"},
				},
			},
		},
		{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "labels"},
					{Key: "localField", Value: "book_labels.labelId"},
					{Key: "foreignField", Value: "_id"},
					{Key: "as", Value: "labels"},
				},
			},
		},
		{
			{Key: "$project",
				Value: bson.D{
					{Key: "book_labels", Value: 0},
					{Key: "labels.createdAt", Value: 0},
					{Key: "labels.updatedAt", Value: 0},
				},
			},
		},
		{{Key: "$sort", Value: bson.D{{Key: "updatedAt", Value: -1}}}},
	}
}

func (r *BookRepositoryImpl) CreateBook(book *model.BookCreate) (*model.BookCreate, error) {
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	book.IsArchived = book.IsArchived || false

	newBook, err := r.books.InsertOne(context.Background(), book)
	if err != nil {
		return nil, err
	}

	book.ID = newBook.InsertedID.(primitive.ObjectID)

	return book, nil
}

func (r *BookRepositoryImpl) GetBooks(userId primitive.ObjectID, isArchived bool) ([]model.BookResponse, error) {
	var books []model.BookResponse

	filter := bson.D{
		{Key: "userId", Value: userId},
		{Key: "isArchived", Value: isArchived},
	}
	pipeline := bookAgregate(filter)

	cursor, err := r.books.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &books); err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return []model.BookResponse{}, nil
	}

	return books, nil
}

func (r *BookRepositoryImpl) GetBook(userId primitive.ObjectID, bookId primitive.ObjectID) (*model.BookResponse, error) {
	var book []model.BookResponse

	filter := bson.D{
		{Key: "userId", Value: userId},
		{Key: "_id", Value: bookId},
	}
	pipeline := bookAgregate(filter)

	cursor, err := r.books.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &book); err != nil {
		return nil, err
	}

	if len(book) == 0 {
		return nil, errors.New("book not found")
	}

	return &book[0], nil
}

func (r *BookRepositoryImpl) UpdateBook(book *model.BookUpdate) (*model.BookResponse, error) {
	filter := bson.M{"userId": book.UserId, "_id": book.ID}
	update := bson.M{
		"$set": bson.M{
			"title":     book.Title,
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
	err := result.Decode(&updatedBook)
	if err != nil {
		return nil, err
	}
	return &updatedBook, nil
}

func (r *BookRepositoryImpl) ArchiveBook(book *model.ArchiveBook) (*model.BookResponse, error) {
	filter := bson.M{"userId": book.UserId, "_id": book.ID}
	update := bson.M{
		"$set": bson.M{
			"isArchived": book.IsArchived,
			"updatedAt":  time.Now(),
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
	err := result.Decode(&updatedBook)
	if err != nil {
		return nil, err
	}

	return &updatedBook, nil
}
