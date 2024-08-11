package repository

import (
	"context"
	"errors"
	"noto/internal/config"
	"noto/internal/services/labels/model"
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
	GetBookByLabel(labelName string) ([]model.BookResponse, error)
}

type LabelRepositoryImpl struct {
	labels      *mongo.Collection
	book_labels *mongo.Collection
}

func NewLabelRepository() LabelRepository {
	return &LabelRepositoryImpl{labels: config.DB.Collection("labels"), book_labels: config.DB.Collection("book_labels")}
}

func (r *LabelRepositoryImpl) CheckAndInsertLabel(label *model.LabelCreate) (*model.LabelCreate, error) {
	filter := bson.M{"name": label.Name}
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
		Name:      book.LabelName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := r.CheckAndInsertLabel(label)
	if err != nil {
		return nil, err
	}

	bookLabel := bson.M{
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
	filter := bson.M{"name": book.LabelName}
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

func (r *LabelRepositoryImpl) GetBookByLabel(labelName string) ([]model.BookResponse, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "name", Value: labelName}}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "book_labels"},
			{Key: "localField", Value: "_id"},
			{Key: "foreignField", Value: "labelId"},
			{Key: "as", Value: "book_labels"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$book_labels"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "labelId", Value: "$book_labels.labelId"},
			{Key: "bookId", Value: "$book_labels.bookId"},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "books"},
			{Key: "localField", Value: "bookId"},
			{Key: "foreignField", Value: "_id"},
			{Key: "as", Value: "books"},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$books"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: "$books._id"},
			{Key: "title", Value: "$books.title"},
			{Key: "createdAt", Value: "$books.createdAt"},
			{Key: "updatedAt", Value: "$books.updatedAt"},
			{Key: "isArchived", Value: "$books.isArchived"},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "title", Value: bson.D{{Key: "$first", Value: "$title"}}},
			{Key: "createdAt", Value: bson.D{{Key: "$first", Value: "$createdAt"}}},
			{Key: "updatedAt", Value: bson.D{{Key: "$first", Value: "$updatedAt"}}},
			{Key: "isArchived", Value: bson.D{{Key: "$first", Value: "$isArchived"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "updatedAt", Value: -1}}}},
	}

	cursor, err := r.labels.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	var results []model.BookResponse
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return []model.BookResponse{}, nil
	}

	return results, err
}
