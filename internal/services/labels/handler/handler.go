package handler

import (
	"noto/internal/services/labels/model"
	"noto/internal/services/labels/service"

	"github.com/gofiber/fiber/v2"
)

type LabelHandler interface {
	CreateLabel(c *fiber.Ctx) error
	GetLabels(c *fiber.Ctx) error
	DeleteLabel(c *fiber.Ctx) error
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
// @Param		book	body		model.Label	true	"Label to create"
// @Success		201		{object}	model.Label
// @Router		/labels [post]
func (s *LabelHandlerImpl) CreateLabel(c *fiber.Ctx) error {
	label := new(model.Label)
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

	return c.JSON(newLabel)
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
