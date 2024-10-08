package repository

import (
	"context"
	"errors"
	"time"

	"noto/internal/services/books/model"
	"noto/internal/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BookRepository interface {
	CreateBook(book *model.BookCreate) (*model.BookCreate, error)
	GetBooks(userId primitive.ObjectID, isArchived bool, page int, limit int) (*model.PaginatedBookResponse, error)
	GetBook(userId primitive.ObjectID, bookId primitive.ObjectID) (*model.BookResponse, error)
	UpdateBook(book *model.BookUpdate) (*model.BookResponse, error)
	ArchiveBook(book *model.ArchiveBook) (*model.BookResponse, error)
}

type BookRepositoryImpl struct {
	books *mongo.Collection
}

func NewBookRepository(db *mongo.Database) BookRepository {
	return &BookRepositoryImpl{books: db.Collection("books")}
}

func bookAgregate(matchCondition bson.D, page int, limit int, usePagination bool) mongo.Pipeline {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchCondition}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "book_labels",
			"localField":   "_id",
			"foreignField": "bookId",
			"as":           "book_labels",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from":         "labels",
			"localField":   "book_labels.labelId",
			"foreignField": "_id",
			"as":           "labels",
		}}},
		{{Key: "$project", Value: bson.M{
			"book_labels":      0,
			"labels.createdAt": 0,
			"labels.updatedAt": 0,
		}}},
		{{Key: "$sort", Value: bson.M{"updatedAt": -1}}},
	}

	if usePagination {
		pipeline = append(pipeline,
			bson.D{{Key: "$facet", Value: utils.PaginationAggregate(page, limit)}},
			bson.D{{Key: "$unwind", Value: "$metadata"}},
		)
	}

	return pipeline
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

func (r *BookRepositoryImpl) GetBooks(userId primitive.ObjectID, isArchived bool, page int, limit int) (*model.PaginatedBookResponse, error) {
	var books []model.PaginatedBookResponse

	filter := bson.D{
		{Key: "userId", Value: userId},
		{Key: "isArchived", Value: isArchived},
	}
	pipeline := bookAgregate(filter, page, limit, true)

	cursor, err := r.books.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.Background(), &books); err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return &model.PaginatedBookResponse{
			Data:     []model.BookResponse{},
			Metadata: model.PaginationMetadata{},
		}, nil
	}

	return &books[0], nil
}

func (r *BookRepositoryImpl) GetBook(userId primitive.ObjectID, bookId primitive.ObjectID) (*model.BookResponse, error) {
	var book []model.BookResponse

	filter := bson.D{
		{Key: "userId", Value: userId},
		{Key: "_id", Value: bookId},
	}
	pipeline := bookAgregate(filter, 0, 0, false)

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
