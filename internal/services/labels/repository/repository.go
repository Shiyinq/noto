package repository

import (
	"context"
	"errors"
	"noto/internal/services/labels/model"
	"noto/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LabelRepository interface {
	CheckAndInsertLabel(label *model.LabelCreate) (*model.LabelCreate, error)
	CreateLabel(label *model.LabelCreate) (*model.LabelCreate, error)
	GetLabels(userId primitive.ObjectID) ([]model.LabelResponse, error)
	DeleteLabel(userId primitive.ObjectID, labelId primitive.ObjectID) error
	AddBookLabel(book *model.BookLabel) (*model.AddBookLabelResponse, error)
	DeleteBookLabel(book *model.BookLabel) error
	GetBookByLabel(userId primitive.ObjectID, labelName string, page int, limit int) (*model.PaginatedBookResponse, error)
}

type LabelRepositoryImpl struct {
	labels      *mongo.Collection
	book_labels *mongo.Collection
}

func NewLabelRepository(db *mongo.Database) LabelRepository {
	return &LabelRepositoryImpl{labels: db.Collection("labels"), book_labels: db.Collection("book_labels")}
}

func (r *LabelRepositoryImpl) CheckAndInsertLabel(label *model.LabelCreate) (*model.LabelCreate, error) {
	filter := bson.M{"userId": label.UserId, "name": label.Name}
	existingLabel := r.labels.FindOne(context.Background(), filter)

	if existingLabel.Err() != nil {
		if existingLabel.Err() == mongo.ErrNoDocuments {
			newLabel, err := r.labels.InsertOne(context.Background(), label)
			if err != nil {
				return nil, err
			}
			label.ID = newLabel.InsertedID.(primitive.ObjectID)
		} else {
			return nil, existingLabel.Err()
		}
	} else {
		err := existingLabel.Decode(&label)
		if err != nil {
			return nil, err
		}
	}

	return label, nil
}

func (r *LabelRepositoryImpl) CreateLabel(label *model.LabelCreate) (*model.LabelCreate, error) {
	label.CreatedAt = time.Now()
	label.UpdatedAt = time.Now()

	result, err := r.CheckAndInsertLabel(label)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *LabelRepositoryImpl) AddBookLabel(book *model.BookLabel) (*model.AddBookLabelResponse, error) {
	label := &model.LabelCreate{
		UserId:    book.UserId,
		Name:      book.LabelName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.CheckAndInsertLabel(label)
	if err != nil {
		return nil, err
	}

	bookLabel := bson.M{
		"userId":  book.UserId,
		"bookId":  book.BookId,
		"labelId": result.ID,
	}
	bookLabelResult, err := r.book_labels.InsertOne(context.Background(), bookLabel)
	if err != nil {
		return nil, err
	}

	var newBookLabel model.AddBookLabelResponse
	err = r.book_labels.FindOne(context.Background(), bson.M{"_id": bookLabelResult.InsertedID}).Decode(&newBookLabel)
	if err != nil {
		return nil, err
	}

	return &newBookLabel, nil
}

func (r *LabelRepositoryImpl) DeleteBookLabel(book *model.BookLabel) error {
	filter := bson.M{"userId": book.UserId, "name": book.LabelName}
	labelExist := r.labels.FindOne(context.Background(), filter)

	if labelExist.Err() != nil {
		return labelExist.Err()
	}

	var existingLabel model.LabelResponse
	labelExist.Decode(&existingLabel)

	filterDelete := bson.M{"bookId": book.BookId, "labelId": existingLabel.ID}

	deleted, err := r.book_labels.DeleteMany(context.Background(), filterDelete)
	if err != nil {
		return err
	}

	if deleted.DeletedCount == 0 {
		return errors.New("book label not found or not deleted")
	}

	return nil
}

func (r *LabelRepositoryImpl) GetLabels(userId primitive.ObjectID) ([]model.LabelResponse, error) {
	var labels []model.LabelResponse

	cursor, err := r.labels.Find(context.Background(), bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &labels); err != nil {
		return nil, err
	}

	if len(labels) == 0 {
		return []model.LabelResponse{}, nil
	}

	return labels, nil
}

func (r *LabelRepositoryImpl) DeleteLabel(userId primitive.ObjectID, labelId primitive.ObjectID) error {
	filter := bson.M{"_id": labelId, "userId": userId}
	deleted, err := r.labels.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if deleted.DeletedCount == 0 {
		return errors.New("label not found or not deleted")
	}

	return nil
}

func (r *LabelRepositoryImpl) GetBookByLabel(userId primitive.ObjectID, labelName string, page int, limit int) (*model.PaginatedBookResponse, error) {
	pipeline := mongo.Pipeline{
		{{
			Key: "$match", Value: bson.M{
				"userId": userId,
				"name":   labelName,
			},
		}},
		{{
			Key: "$lookup", Value: bson.M{
				"from":         "book_labels",
				"localField":   "_id",
				"foreignField": "labelId",
				"as":           "book_labels",
			},
		}},
		{{
			Key: "$unwind", Value: bson.M{
				"path":                       "$book_labels",
				"preserveNullAndEmptyArrays": true,
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"_id":     0,
				"labelId": "$book_labels.labelId",
				"bookId":  "$book_labels.bookId",
			},
		}},
		{{
			Key: "$lookup", Value: bson.M{
				"from":         "books",
				"localField":   "bookId",
				"foreignField": "_id",
				"as":           "books",
			},
		}},
		{{
			Key: "$unwind", Value: bson.M{
				"path":                       "$books",
				"preserveNullAndEmptyArrays": false,
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"_id":        "$books._id",
				"title":      "$books.title",
				"createdAt":  "$books.createdAt",
				"updatedAt":  "$books.updatedAt",
				"isArchived": "$books.isArchived",
			},
		}},
		{{
			Key: "$group", Value: bson.M{
				"_id":        "$_id",
				"title":      bson.M{"$first": "$title"},
				"createdAt":  bson.M{"$first": "$createdAt"},
				"updatedAt":  bson.M{"$first": "$updatedAt"},
				"isArchived": bson.M{"$first": "$isArchived"},
			},
		}},
		{{
			Key: "$sort", Value: bson.M{"updatedAt": -1},
		}},
		{{
			Key: "$lookup", Value: bson.M{
				"from":         "book_labels",
				"localField":   "_id",
				"foreignField": "bookId",
				"as":           "book_labels",
			},
		}},
		{{
			Key: "$lookup", Value: bson.M{
				"from":         "labels",
				"localField":   "book_labels.labelId",
				"foreignField": "_id",
				"as":           "labels",
			},
		}},
		{{
			Key: "$project", Value: bson.M{
				"book_labels":      0,
				"labels.userId":    0,
				"labels.createdAt": 0,
				"labels.updatedAt": 0,
			},
		}},
		{{
			Key: "$sort", Value: bson.M{"updatedAt": -1},
		}},
		{{Key: "$facet", Value: utils.PaginationAggregate(page, limit)}},
		{{Key: "$unwind", Value: "$metadata"}},
	}

	cursor, err := r.labels.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var books []model.PaginatedBookResponse
	if err = cursor.All(context.Background(), &books); err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return &model.PaginatedBookResponse{
			Data:     []model.BookResponse{},
			Metadata: model.PaginationMetadata{},
		}, nil
	}

	return &books[0], err
}
