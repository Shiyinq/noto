package labels_router

import (
	"noto/internal/services/labels/handler"
	"noto/internal/services/labels/repository"
	"noto/internal/services/labels/service"

	"github.com/gofiber/fiber/v2"
)

var repo = repository.NewLabelRepository()
var serv = service.NewLabelService(repo)
var hand = handler.NewLabelHandler(serv)

func LabelsRouter(app *fiber.App) {
	app.Get("/labels", hand.GetLabels)
	app.Get("/labels/:id", hand.GetLabel)
}
