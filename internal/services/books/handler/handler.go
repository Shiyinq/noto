package handler

import (
	"noto/internal/services/books/service"

	"github.com/gofiber/fiber/v2"
)

type BookHandler interface {
	GetBooks(c *fiber.Ctx) error
	GetBook(c *fiber.Ctx) error
}

type BookHandlerImpl struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) BookHandler {
	return &BookHandlerImpl{bookService: bookService}
}

func (s *BookHandlerImpl) GetBooks(c *fiber.Ctx) error {
	books, _ := s.bookService.GetAllBooks()
	return c.JSON(books)
}

func (s *BookHandlerImpl) GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	label, err := s.bookService.GetBook(id)
	if err != nil {
		return c.Status(404).SendString("Book not found")
	}
	return c.JSON(label)
}
