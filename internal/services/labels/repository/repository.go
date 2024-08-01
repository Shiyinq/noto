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
	CreateLabel(label *model.LabelCreate) (*model.LabelCreate, error)
	GetLabels() ([]model.LabelResponse, error)
	DeleteLabel(labelId string) error
}

type LabelRepositoryImpl struct {
	labels *mongo.Collection
}

func NewLabelRepository() LabelRepository {
	return &LabelRepositoryImpl{labels: config.DB.Collection("labels")}
}

func (r *LabelRepositoryImpl) CreateLabel(label *model.LabelCreate) (*model.LabelCreate, error) {
	label.CreatedAt = time.Now()
	label.UpdatedAt = time.Now()

	filter := bson.M{"name": label.Name}
	nameIsExist := r.labels.FindOne(context.Background(), filter)

	if nameIsExist.Err() != nil {
		if nameIsExist.Err() == mongo.ErrNoDocuments {
			newLabel, err := r.labels.InsertOne(context.Background(), label)
			if err != nil {
				return nil, err
			}
			label.ID = newLabel.InsertedID.(primitive.ObjectID)
			return label, nil
		}
	}

	var exisLabel model.LabelCreate
	err := nameIsExist.Decode(&exisLabel)

	if err != nil {
		return nil, err
	}

	return &exisLabel, nil
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
