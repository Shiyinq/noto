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
	GetLabels() ([]model.LabelResponse, error)
	DeleteLabel(labelId string) error
	AddBookLabel(book *model.AddBookLabel) (*model.AddBookLabelResponse, error)
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

func (r *LabelRepositoryImpl) AddBookLabel(book *model.AddBookLabel) (*model.AddBookLabelResponse, error) {
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

func (r *LabelRepositoryImpl) GetLabels() ([]model.LabelResponse, error) {
	var labels []model.LabelResponse

	cursor, err := r.labels.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &labels); err != nil {
		return nil, err
	}

	return labels, nil
}

func (r *LabelRepositoryImpl) DeleteLabel(labelId string) error {
	labelObjectID, err := primitive.ObjectIDFromHex(labelId)
	if err != nil {
		return errors.New("invalid bookId format")
	}

	filter := bson.M{"_id": labelObjectID}
	deleted, err := r.labels.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if deleted.DeletedCount == 0 {
		return errors.New("label note found or not deleted")
	}

	return nil
}
