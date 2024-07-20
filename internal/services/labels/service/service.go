package service

import (
	model "noto/internal/services/labels/model"
	"noto/internal/services/labels/repository"
)

type LabelService interface {
	GetAllLabels() []model.Label
	GetLabelByID(id string) (*model.Label, error)
}

type LabelServiceImpl struct {
	labelRepo repository.LabelRepository
}

func NewLabelService(labelRepo repository.LabelRepository) LabelService {
	return &LabelServiceImpl{labelRepo: labelRepo}
}

func (r *LabelServiceImpl) GetAllLabels() []model.Label {
	return r.labelRepo.GetAllLabels()
}

func (r *LabelServiceImpl) GetLabelByID(id string) (*model.Label, error) {
	return r.labelRepo.GetLabelByID(id)
}
