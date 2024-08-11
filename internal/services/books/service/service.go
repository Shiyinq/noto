package service

import (
	"noto/internal/services/books/model"
	"noto/internal/services/books/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookService interface {
	CreateBook(book *model.BookCreate) (*model.BookCreate, error)
	GetBooks(userId primitive.ObjectID, isArchived bool, page int, limit int) (*model.PaginatedBookResponse, error)
	GetBook(userId primitive.ObjectID, bookId primitive.ObjectID) (*model.BookResponse, error)
	UpdateBook(book *model.BookUpdate) (*model.BookResponse, error)
	ArchiveBook(book *model.ArchiveBook) (*model.BookResponse, error)
}

type BookServiceImpl struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &BookServiceImpl{bookRepo: bookRepo}
}

func (r *BookServiceImpl) CreateBook(book *model.BookCreate) (*model.BookCreate, error) {
	return r.bookRepo.CreateBook(book)
}

func (r *BookServiceImpl) GetBooks(userId primitive.ObjectID, isArchived bool, page int, limit int) (*model.PaginatedBookResponse, error) {
	return r.bookRepo.GetBooks(userId, isArchived, page, limit)
}

func (r *BookServiceImpl) GetBook(userId primitive.ObjectID, bookId primitive.ObjectID) (*model.BookResponse, error) {
	return r.bookRepo.GetBook(userId, bookId)
}

func (s *BookServiceImpl) UpdateBook(book *model.BookUpdate) (*model.BookResponse, error) {
	return s.bookRepo.UpdateBook(book)
}

func (s *BookServiceImpl) ArchiveBook(book *model.ArchiveBook) (*model.BookResponse, error) {
	return s.bookRepo.ArchiveBook(book)
}
