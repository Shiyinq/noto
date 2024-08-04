package service

import (
	"noto/internal/services/labels/model"
	"noto/internal/services/labels/repository"
)

type LabelService interface {
	CreateLabel(label *model.LabelCreate) (*model.LabelCreate, error)
	GetLabels() ([]model.LabelResponse, error)
	DeleteLabel(labelId string) error
	AddBookLabel(book *model.BookLabel) (*model.AddBookLabelResponse, error)
	DeleteBookLabel(book *model.BookLabel) error
	GetBookByLabel(labelName string) ([]model.BookResponse, error)
}

type LabelServiceImpl struct {
	labelRepo repository.LabelRepository
}

func NewLabelService(labelRepo repository.LabelRepository) LabelService {
	return &LabelServiceImpl{labelRepo: labelRepo}
}

func (r *LabelServiceImpl) CreateLabel(label *model.LabelCreate) (*model.LabelCreate, error) {
	return r.labelRepo.CreateLabel(label)
}

func (r *LabelServiceImpl) GetLabels() ([]model.LabelResponse, error) {
	return r.labelRepo.GetLabels()
}

func (r *LabelServiceImpl) DeleteLabel(labelId string) error {
	return r.labelRepo.DeleteLabel(labelId)
}

func (r *LabelServiceImpl) AddBookLabel(book *model.BookLabel) (*model.AddBookLabelResponse, error) {
	return r.labelRepo.AddBookLabel(book)
}

func (r *LabelServiceImpl) DeleteBookLabel(book *model.BookLabel) error {
	return r.labelRepo.DeleteBookLabel(book)
}

func (r *LabelServiceImpl) GetBookByLabel(labelName string) ([]model.BookResponse, error) {
	return r.labelRepo.GetBookByLabel(labelName)
}
