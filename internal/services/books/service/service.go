package service

import (
	"noto/internal/services/books/model"
	"noto/internal/services/books/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService interface {
	CreateBook(book *model.BookCreate) (*model.BookCreate, error)
	GetBooks(userId primitive.ObjectID, isArchived bool) ([]model.BookResponse, error)
	GetBook(id string) (*model.BookResponse, error)
	UpdateBook(id string, title string) (*model.BookResponse, error)
	ArchiveBook(id string, book *model.ArchiveBook) (*model.BookResponse, error)
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

func (r *BookServiceImpl) GetBook(id string) (*model.BookResponse, error) {
	return r.bookRepo.GetBookByID(id)
}

func (s *BookServiceImpl) UpdateBook(id string, title string) (*model.BookResponse, error) {
	updatedBook, err := s.bookRepo.UpdateBook(id, title)
	if err != nil {
		return nil, err
	}
	return updatedBook, nil
}

func (s *BookServiceImpl) ArchiveBook(id string, book *model.ArchiveBook) (*model.BookResponse, error) {
	archived, err := s.bookRepo.ArchiveBook(id, book)
	if err != nil {
		return nil, err
	}
	return archived, nil
}
