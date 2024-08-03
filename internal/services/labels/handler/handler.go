package handler

import (
	"noto/internal/services/labels/model"
	"noto/internal/services/labels/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LabelHandler interface {
	CreateLabel(c *fiber.Ctx) error
	GetLabels(c *fiber.Ctx) error
	DeleteLabel(c *fiber.Ctx) error
	AddBookLabel(c *fiber.Ctx) error
	DeleteBookLabel(c *fiber.Ctx) error
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
// @Accept		json
// @Produce		json
// @Param		book	body		model.LabelCreateSwagger	true	"Label to create"
// @Success		201		{object}	model.LabelCreate
// @Router		/labels [post]
func (s *LabelHandlerImpl) CreateLabel(c *fiber.Ctx) error {
	label := new(model.LabelCreate)
	if err := c.BodyParser(label); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newLabel, err := s.labelService.CreateLabel(label)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(newLabel)
}

// GetLabels
// @Summary		Get all labels
// @Description	Get all labels
// @Tags		Labels
// @Produce		json
// @Success		200		{object}	[]model.LabelResponse
// @Router		/labels [get]
func (s *LabelHandlerImpl) GetLabels(c *fiber.Ctx) error {
	labels, err := s.labelService.GetLabels()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(labels)
}

// DeleteLabel
// @Summary		Delete label
// @Description	Delete label
// @Tags		Labels
// @Accept		json
// @Produce		json
// @Param		labelId path string true "Label ID"
// @Success		200	{object} interface{}
// @Router		/labels/{labelId} [delete]
func (s *LabelHandlerImpl) DeleteLabel(c *fiber.Ctx) error {
	labelId := c.Params("labelId")

	if err := s.labelService.DeleteLabel(labelId); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": "label deleted",
	})
}

// AddBookLabel
// @Summary		Add label to book
// @Description	Add label to book
// @Tags		Labels
// @Accept		json
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Param		book	body		model.BookLabelSwagger	true	"Label to add"
// @Success		201	{object}	model.AddBookLabelResponse
// @Router		/books/{bookId}/labels [post]
func (s *LabelHandlerImpl) AddBookLabel(c *fiber.Ctx) error {
	bookId := c.Params("bookId")
	label := new(model.BookLabel)
	objectBookId, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	label.BookId = objectBookId

	if err := c.BodyParser(label); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	added, err := s.labelService.AddBookLabel(label)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(added)
}

// DeleteBookLabel
// @Summary		Delete label from book
// @Description	Delete label from book
// @Tags		Labels
// @Accept		json
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Param		book	body		model.BookLabelSwagger	true	"Label to delete"
// @Success		200	{object} interface{}
// @Router		/books/{bookId}/labels [delete]
func (s *LabelHandlerImpl) DeleteBookLabel(c *fiber.Ctx) error {
	bookId := c.Params("bookId")
	label := new(model.BookLabel)
	objectBookId, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	label.BookId = objectBookId

	if err := c.BodyParser(label); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	deleted := s.labelService.DeleteBookLabel(label)
	if deleted != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": deleted.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": "deleted",
	})
}
