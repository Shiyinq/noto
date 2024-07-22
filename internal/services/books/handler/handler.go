package handler

import (
	"noto/internal/services/books/model"
	"noto/internal/services/books/service"

	"github.com/gofiber/fiber/v2"
)

type BookHandler interface {
	CreateBook(c *fiber.Ctx) error
	GetBooks(c *fiber.Ctx) error
	GetBook(c *fiber.Ctx) error
}

type BookHandlerImpl struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) BookHandler {
	return &BookHandlerImpl{bookService: bookService}
}

func (s *BookHandlerImpl) CreateBook(c *fiber.Ctx) error {
	note := new(model.Book)
	if err := c.BodyParser(note); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	newNote, _ := s.bookService.CreateBook(note)
	return c.JSON(newNote)
}

func (s *BookHandlerImpl) GetBooks(c *fiber.Ctx) error {
	books, err := s.bookService.GetAllBooks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(books)
}

func (s *BookHandlerImpl) GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	book, err := s.bookService.GetBook(id)
	if err != nil {
		if err.Error() == "book not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(book)
}
