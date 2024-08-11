package service

import (
	"noto/internal/services/labels/model"
	"noto/internal/services/labels/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LabelService interface {
	CreateLabel(label *model.LabelCreate) (*model.LabelCreate, error)
	GetLabels(userId primitive.ObjectID) ([]model.LabelResponse, error)
	DeleteLabel(userId primitive.ObjectID, labelId primitive.ObjectID) error
	AddBookLabel(book *model.BookLabel) (*model.AddBookLabelResponse, error)
	DeleteBookLabel(book *model.BookLabel) error
	GetBookByLabel(userId primitive.ObjectID, labelName string, page int, limit int) (*model.PaginatedBookResponse, error)
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

func (r *LabelServiceImpl) GetLabels(userId primitive.ObjectID) ([]model.LabelResponse, error) {
	return r.labelRepo.GetLabels(userId)
}

func (r *LabelServiceImpl) DeleteLabel(userId primitive.ObjectID, labelId primitive.ObjectID) error {
	return r.labelRepo.DeleteLabel(userId, labelId)
}

func (r *LabelServiceImpl) AddBookLabel(book *model.BookLabel) (*model.AddBookLabelResponse, error) {
	return r.labelRepo.AddBookLabel(book)
}

func (r *LabelServiceImpl) DeleteBookLabel(book *model.BookLabel) error {
	return r.labelRepo.DeleteBookLabel(book)
}

func (r *LabelServiceImpl) GetBookByLabel(userId primitive.ObjectID, labelName string, page int, limit int) (*model.PaginatedBookResponse, error) {
	return r.labelRepo.GetBookByLabel(userId, labelName, page, limit)
}
