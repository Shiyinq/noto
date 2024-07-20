package repository

import (
	"errors"
	"noto/internal/services/books/model"
)

type BookRepository interface {
	GetAllBooks() ([]model.Book, error)
	GetBookByID(id string) (*model.Book, error)
	// CreateBook(book *model.Book) (*model.Book, error)
	// UpdateBook(book *model.Book) (*model.Book, error)
	// DeleteBook(book *model.Book) error
}

type BookRepositoryImpl struct {
	books []model.Book
}

func NewBookRepository() BookRepository {
	return &BookRepositoryImpl{books: []model.Book{}}
}

func (r *BookRepositoryImpl) GetAllBooks() ([]model.Book, error) {
	return r.books, nil
}

func (r *BookRepositoryImpl) GetBookByID(id string) (*model.Book, error) {
	for _, label := range r.books {
		if label.ID == id {
			return &label, nil
		}
	}
	return nil, errors.New("label not found")
}
