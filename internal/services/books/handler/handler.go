package handler

import (
	"noto/internal/services/books/model"
	"noto/internal/services/books/service"
	"noto/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BookHandler interface {
	CreateBook(c *fiber.Ctx) error
	GetBooks(c *fiber.Ctx) error
	GetBook(c *fiber.Ctx) error
	UpdateBook(c *fiber.Ctx) error
	ArchiveBook(c *fiber.Ctx) error
}

type BookHandlerImpl struct {
	bookService service.BookService
}

func NewBookHandler(bookService service.BookService) BookHandler {
	return &BookHandlerImpl{bookService: bookService}
}

// CreateBook
// @Summary		Create a new book
// @Description	Create a new book
// @Tags		Books
// @Security 	BearerAuth
// @Accept		json
// @Produce		json
// @Param		book	body		model.BookCreateSwagger	true	"Book to create"
// @Success		201		{object}	model.BookCreate
// @Router		/api/books [post]
func (s *BookHandlerImpl) CreateBook(c *fiber.Ctx) error {
	objUserId, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	book := new(model.BookCreate)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse JSON: " + err.Error(),
		})
	}

	book.UserID = objUserId
	newBook, err := s.bookService.CreateBook(book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create book: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newBook)
}

// GetBooks
// @Summary		Get all book
// @Description	Get all book
// @Tags		Books
// @Security 	BearerAuth
// @Produce		json
// @Success		200		{object}	[]model.BookResponse
// @Router		/api/books [get]
func (s *BookHandlerImpl) GetBooks(c *fiber.Ctx) error {
	objUserId, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	isArchivedStr := c.Query("is_archived")

	var isArchived bool
	if isArchivedStr != "" {
		var errConv error
		isArchived, errConv = strconv.ParseBool(isArchivedStr)
		if errConv != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid value for is_archived"})
		}
	}

	books, err := s.bookService.GetBooks(objUserId, isArchived)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(books)
}

// GetBook
// @Summary		Get book by id
// @Description	Get book by id
// @Tags		Books
// @Security 	BearerAuth
// @Produce		json
// @Param 		id path string true "Book ID"
// @Success		200		{object}	model.BookResponse
// @Failure 	404 {object} fiber.Map
// @Router		/api/books/{id} [get]
func (s *BookHandlerImpl) GetBook(c *fiber.Ctx) error {
	objUserId, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	bookId, err := utils.ToObjectID(c.Params("id"))
	if err != nil {
		return err
	}

	book, err := s.bookService.GetBook(objUserId, bookId)
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

// UpdateBook
// @Summary		Update book by id
// @Description	Update book by id
// @Tags		Books
// @Security 	BearerAuth
// @Produce		json
// @Accept		json
// @Param 		id path string true "Book ID"
// @Param		book	body		model.BookUpdateSwagger	true	"Book to update"
// @Success		200		{object}	model.BookResponse
// @Router		/api/books/{id} [put]
func (s *BookHandlerImpl) UpdateBook(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	bookId, err := utils.ToObjectID(c.Params("id"))
	if err != nil {
		return err
	}

	book := new(model.BookUpdate)
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	title := book.Title
	if title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Title is required",
		})
	}

	book.ID = bookId
	book.UserID = userId
	updatedBook, err := s.bookService.UpdateBook(book)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	return c.JSON(updatedBook)
}

// ArchiveBook
// @Summary		Archive book by id
// @Description	Archive book by id
// @Tags		Books
// @Security 	BearerAuth
// @Produce		json
// @Accept		json
// @Param 		id path string true "Book ID"
// @Param		book	body		model.ArchiveBook true	"Book to archive"
// @Success		200		{object}	model.BookResponse
// @Router		/api/books/{id} [patch]
func (s *BookHandlerImpl) ArchiveBook(c *fiber.Ctx) error {
	objUserId, err := utils.GetUserID(c)
	if err != nil {
		return err
	}

	bookId := c.Params("id")
	archive := new(model.ArchiveBook)
	if err := c.BodyParser(&archive); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	archived, err := s.bookService.ArchiveBook(objUserId, bookId, archive)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	return c.JSON(archived)
}
