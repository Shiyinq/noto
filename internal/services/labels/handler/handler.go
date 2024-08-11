package handler

import (
	"noto/internal/services/labels/model"
	"noto/internal/services/labels/service"
	"noto/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type LabelHandler interface {
	CreateLabel(c *fiber.Ctx) error
	GetLabels(c *fiber.Ctx) error
	DeleteLabel(c *fiber.Ctx) error
	AddBookLabel(c *fiber.Ctx) error
	DeleteBookLabel(c *fiber.Ctx) error
	GetBookByLabel(c *fiber.Ctx) error
}

type LabelHandlerImpl struct {
	labelService service.LabelService
}

func NewLabelHandler(labelService service.LabelService) LabelHandler {
	return &LabelHandlerImpl{labelService: labelService}
}

// CreateLabel
// @Summary		Create a new label
// @Description	Create a new label
// @Tags		Labels
// @Security 	BearerAuth
// @Accept		json
// @Produce		json
// @Param		book	body		model.LabelCreateSwagger	true	"Label to create"
// @Success		201		{object}	model.LabelCreate
// @Router		/api/labels [post]
func (s *LabelHandlerImpl) CreateLabel(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	label := new(model.LabelCreate)
	if err := c.BodyParser(label); err != nil {
		return utils.ErrorBadRequest(c, "failed to parse json: "+err.Error())
	}

	label.UserId = userId
	newLabel, err := s.labelService.CreateLabel(label)
	if err != nil {
		return utils.ErrorInternalServer(c, "failed to create label: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newLabel)
}

// GetLabels
// @Summary		Get all labels
// @Description	Get all labels
// @Tags		Labels
// @Security 	BearerAuth
// @Produce		json
// @Success		200		{object}	[]model.LabelResponse
// @Router		/api/labels [get]
func (s *LabelHandlerImpl) GetLabels(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	labels, err := s.labelService.GetLabels(userId)
	if err != nil {
		return utils.ErrorInternalServer(c, err.Error())
	}

	return c.JSON(labels)
}

// DeleteLabel
// @Summary		Delete label
// @Description	Delete label
// @Tags		Labels
// @Security 	BearerAuth
// @Accept		json
// @Produce		json
// @Param		labelId path string true "Label ID"
// @Success		200	{object} interface{}
// @Router		/api/labels/{labelId} [delete]
func (s *LabelHandlerImpl) DeleteLabel(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	labelId, err := utils.ToObjectID(c.Params("labelId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	if err := s.labelService.DeleteLabel(userId, labelId); err != nil {
		return utils.ErrorInternalServer(c, "failed to delete label: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": "label deleted",
	})
}

// AddBookLabel
// @Summary		Add label to book
// @Description	Add label to book
// @Tags		Labels
// @Security 	BearerAuth
// @Accept		json
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Param		book	body		model.BookLabelSwagger	true	"Label to add"
// @Success		201	{object}	model.AddBookLabelResponse
// @Router		/api/books/{bookId}/labels [post]
func (s *LabelHandlerImpl) AddBookLabel(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	label := new(model.BookLabel)
	if err := c.BodyParser(label); err != nil {
		return utils.ErrorBadRequest(c, "failed to parse json: "+err.Error())
	}

	label.UserId = userId
	label.BookId = bookId
	added, err := s.labelService.AddBookLabel(label)
	if err != nil {
		return utils.ErrorInternalServer(c, "failed to add book label: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(added)
}

// DeleteBookLabel
// @Summary		Delete label from book
// @Description	Delete label from book
// @Tags		Labels
// @Security 	BearerAuth
// @Accept		json
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Param		book	body		model.BookLabelSwagger	true	"Label to delete"
// @Success		200	{object} interface{}
// @Router		/api/books/{bookId}/labels [delete]
func (s *LabelHandlerImpl) DeleteBookLabel(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	label := new(model.BookLabel)
	if err := c.BodyParser(label); err != nil {
		return utils.ErrorBadRequest(c, "failed to parse json: "+err.Error())
	}

	label.UserId = userId
	label.BookId = bookId
	deleted := s.labelService.DeleteBookLabel(label)
	if deleted != nil {
		return utils.ErrorInternalServer(c, "failed to delete book label: "+deleted.Error())
	}

	return c.JSON(fiber.Map{
		"success": "deleted",
	})
}

// GetBookByLabel
// @Summary		Get book by label name
// @Description	Get book by label name
// @Tags		Labels
// @Security 	BearerAuth
// @Produce		json
// @Param		labelName path string true "Label Name"
// @Param		page		query		int		false	"Page number for pagination"	minimum(1)
// @Param		limit		query		int		false	"Number of items per page"	minimum(1)
// @Success		200	{object} model.PaginatedBookResponse
// @Router		/api/labels/{labelName}/books [get]
func (s *LabelHandlerImpl) GetBookByLabel(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	labelName := c.Params("labelName")
	if labelName == "" {
		return utils.ErrorBadRequest(c, "label name required!")
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	books, err := s.labelService.GetBookByLabel(userId, labelName, page, limit)

	if err != nil {
		return utils.ErrorInternalServer(c, err.Error())
	}

	return c.JSON(books)
}
