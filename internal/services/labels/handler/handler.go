package handler

import (
	"noto/internal/services/labels/service"

	"github.com/gofiber/fiber/v2"
)

type LabelHandler interface {
	GetLabels(c *fiber.Ctx) error
	GetLabel(c *fiber.Ctx) error
}

type LabelHandlerImpl struct {
	labelService service.LabelService
}

func NewLabelHandler(labelService service.LabelService) LabelHandler {
	return &LabelHandlerImpl{labelService: labelService}
}

func (s *LabelHandlerImpl) GetLabels(c *fiber.Ctx) error {
	labels := s.labelService.GetAllLabels()
	return c.JSON(labels)
}

func (s *LabelHandlerImpl) GetLabel(c *fiber.Ctx) error {
	id := c.Params("id")
	label, err := s.labelService.GetLabelByID(id)
	if err != nil {
		return c.Status(404).SendString("Label not found")
	}
	return c.JSON(label)
}
