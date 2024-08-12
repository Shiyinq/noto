package handler

import (
	_ "noto/internal/common"
	"noto/internal/services/books/model"
	"noto/internal/services/books/service"
	"noto/internal/utils"

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
// @Param 		Authorization header string false "Bearer token"
// @Param		book	body		model.BookCreateSwagger	true	"Book to create"
// @Success		201		{object}	model.BookCreate
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books [post]
func (s *BookHandlerImpl) CreateBook(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	book := new(model.BookCreate)
	if err := c.BodyParser(book); err != nil {
		return utils.ErrorBadRequest(c, "failed to parse json: "+err.Error())
	}

	book.UserId = userId
	newBook, err := s.bookService.CreateBook(book)
	if err != nil {
		return utils.ErrorInternalServer(c, "failed to create book: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newBook)
}

// GetBooks
// @Summary		Get all book
// @Description	Get all book
// @Tags		Books
// @Security 	BearerAuth
// @Produce		json
// @Param 		Authorization header string false "Bearer token"
// @Param		is_archived	query		bool	false	"Filter by archive status"
// @Param		page		query		int		false	"Page number for pagination"	minimum(1)
// @Param		limit		query		int		false	"Number of items per page"	minimum(1)
// @Success		200		{object}	model.PaginatedBookResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books [get]
func (s *BookHandlerImpl) GetBooks(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	isArchived := c.QueryBool("is_archived", false)

	books, err := s.bookService.GetBooks(userId, isArchived, page, limit)
	if err != nil {
		return utils.ErrorInternalServer(c, err.Error())
	}

	return c.JSON(books)
}

// GetBook
// @Summary		Get book by id
// @Description	Get book by id
// @Tags		Books
// @Security 	BearerAuth
// @Produce		json
// @Param 		Authorization header string false "Bearer token"
// @Param 		bookId path string true "Book ID"
// @Success		200		{object}	model.BookResponse
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books/{bookId} [get]
func (s *BookHandlerImpl) GetBook(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	book, err := s.bookService.GetBook(userId, bookId)
	if err != nil {
		if err.Error() == "book not found" {
			return utils.ErrorNotFound(c, err.Error())
		}
		return utils.ErrorInternalServer(c, err.Error())
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
// @Param 		Authorization header string false "Bearer token"
// @Param 		bookId path string true "Book ID"
// @Param		book	body		model.BookUpdateSwagger	true	"Book to update"
// @Success		200		{object}	model.BookResponse
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books/{bookId} [put]
func (s *BookHandlerImpl) UpdateBook(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	book := new(model.BookUpdate)
	if err := c.BodyParser(&book); err != nil {
		return utils.ErrorBadRequest(c, "cannot parse json: "+err.Error())
	}

	title := book.Title
	if title == "" {
		return utils.ErrorBadRequest(c, "title is required")
	}

	book.ID = bookId
	book.UserId = userId
	updatedBook, err := s.bookService.UpdateBook(book)
	if err != nil {
		return utils.ErrorInternalServer(c, "failed to update book: "+err.Error())
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
// @Param 		Authorization header string false "Bearer token"
// @Param 		bookId path string true "Book ID"
// @Param		book	body		model.ArchiveBookSwagger true	"Book to archive"
// @Success		200		{object}	model.BookResponse
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books/{bookId} [patch]
func (s *BookHandlerImpl) ArchiveBook(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	book := new(model.ArchiveBook)
	if err := c.BodyParser(&book); err != nil {
		return utils.ErrorBadRequest(c, "cannot parse json: "+err.Error())
	}

	book.ID = bookId
	book.UserId = userId
	archived, err := s.bookService.ArchiveBook(book)
	if err != nil {
		return utils.ErrorInternalServer(c, "failed to archive book: "+err.Error())
	}

	return c.JSON(archived)
}
