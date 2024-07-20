package repository

import (
	"errors"
	"noto/internal/services/labels/model"
)

type LabelRepository interface {
	GetAllLabels() []model.Label
	GetLabelByID(id string) (*model.Label, error)
	// CreateLabel(label *model.Label) error
	// UpdateLabel(label *model.Label) error
	// DeleteLabel(id string) error
}

type LabelRepositoryImpl struct {
	labels []model.Label
}

func NewLabelRepository() LabelRepository {
	return &LabelRepositoryImpl{
		labels: []model.Label{},
	}
}

func (r *LabelRepositoryImpl) GetAllLabels() []model.Label {
	return r.labels
}

func (r *LabelRepositoryImpl) GetLabelByID(id string) (*model.Label, error) {
	for _, label := range r.labels {
		if label.ID == id {
			return &label, nil
		}
	}
	return nil, errors.New("label not found")
}
