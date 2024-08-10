package service

import (
	"noto/internal/services/books/model"
	"noto/internal/services/books/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService interface {
	CreateBook(book *model.BookCreate) (*model.BookCreate, error)
	GetBooks(userId primitive.ObjectID, isArchived bool) ([]model.BookResponse, error)
	GetBook(userId primitive.ObjectID, bookId string) (*model.BookResponse, error)
	UpdateBook(book *model.BookUpdate) (*model.BookResponse, error)
	ArchiveBook(userId primitive.ObjectID, bookId string, book *model.ArchiveBook) (*model.BookResponse, error)
}

type BookServiceImpl struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &BookServiceImpl{bookRepo: bookRepo}
}

func (r *BookServiceImpl) CreateBook(book *model.BookCreate) (*model.BookCreate, error) {
	book, err := r.bookRepo.CreateBook(book)
	return book, err
}

func (r *BookServiceImpl) GetBooks(userId primitive.ObjectID, isArchived bool) ([]model.BookResponse, error) {
	return r.bookRepo.GetBooks(userId, isArchived)
}

func (r *BookServiceImpl) GetBook(userId primitive.ObjectID, bookId string) (*model.BookResponse, error) {
	return r.bookRepo.GetBook(userId, bookId)
}

func (s *BookServiceImpl) UpdateBook(book *model.BookUpdate) (*model.BookResponse, error) {
	updatedBook, err := s.bookRepo.UpdateBook(book)
	if err != nil {
		return nil, err
	}
	return updatedBook, nil
}

func (s *BookServiceImpl) ArchiveBook(userId primitive.ObjectID, bookId string, book *model.ArchiveBook) (*model.BookResponse, error) {
	archived, err := s.bookRepo.ArchiveBook(userId, bookId, book)
	if err != nil {
		return nil, err
	}
	return archived, nil
}
