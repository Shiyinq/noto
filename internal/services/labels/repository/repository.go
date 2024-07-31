package repository

import (
	"context"
	"noto/internal/config"
	"noto/internal/services/labels/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LabelRepository interface {
	CreateLabel(label *model.Label) (*model.Label, error)
	GetLabels() ([]model.LabelResponse, error)
}

type LabelRepositoryImpl struct {
	labels *mongo.Collection
}

func NewLabelRepository() LabelRepository {
	return &LabelRepositoryImpl{labels: config.DB.Collection("labels")}
}

func (r *LabelRepositoryImpl) CreateLabel(label *model.Label) (*model.Label, error) {
	label.CreatedAt = time.Now()
	label.UpdatedAt = time.Now()

	filter := bson.M{"name": label.Name}
	nameIsExist := r.labels.FindOne(context.Background(), filter)

	if nameIsExist.Err() != nil {
		if nameIsExist.Err() == mongo.ErrNoDocuments {
			_, err := r.labels.InsertOne(context.Background(), label)
			if err != nil {
				return nil, err
			}
			return label, nil
		}
	}

	var exisLabel model.Label
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
