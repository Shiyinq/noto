package service

import (
	"noto/internal/services/books/model"
	"noto/internal/services/books/repository"
)

type BookService interface {
	GetAllBooks() ([]model.Book, error)
	GetBook(id string) (*model.Book, error)
}

type BookServiceImpl struct {
	bookRepo repository.BookRepository
}

func NewBookService(bookRepo repository.BookRepository) BookService {
	return &BookServiceImpl{bookRepo: bookRepo}
}

func (r *BookServiceImpl) GetAllBooks() ([]model.Book, error) {
	return r.bookRepo.GetAllBooks()
}

func (r *BookServiceImpl) GetBook(id string) (*model.Book, error) {
	return r.bookRepo.GetBookByID(id)
}
