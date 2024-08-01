package service

import (
	"noto/internal/services/labels/model"
	"noto/internal/services/labels/repository"
)

type LabelService interface {
	CreateLabel(label *model.Label) (*model.Label, error)
	GetLabels() ([]model.LabelResponse, error)
	DeleteLabel(labelId string) error
}

type LabelServiceImpl struct {
	labelRepo repository.LabelRepository
}

func NewLabelService(labelRepo repository.LabelRepository) LabelService {
	return &LabelServiceImpl{labelRepo: labelRepo}
}

func (r *LabelServiceImpl) CreateLabel(label *model.Label) (*model.Label, error) {
	return r.labelRepo.CreateLabel(label)
}

func (r *LabelServiceImpl) GetLabels() ([]model.LabelResponse, error) {
	return r.labelRepo.GetLabels()
}

func (r *LabelServiceImpl) DeleteLabel(labelId string) error {
	return r.labelRepo.DeleteLabel(labelId)
}
