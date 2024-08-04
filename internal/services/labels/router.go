package labels_router

import (
	"noto/internal/services/labels/handler"
	"noto/internal/services/labels/repository"
	"noto/internal/services/labels/service"

	"github.com/gofiber/fiber/v2"
)

func LabelsRouter(app *fiber.App) {
	var repo = repository.NewLabelRepository()
	var serv = service.NewLabelService(repo)
	var hand = handler.NewLabelHandler(serv)

	app.Post("/labels", hand.CreateLabel)
	app.Get("/labels", hand.GetLabels)
	app.Delete("/labels/:labelId", hand.DeleteLabel)
	app.Post("/books/:bookId/labels", hand.AddBookLabel)
	app.Delete("/books/:bookId/labels", hand.DeleteBookLabel)
	app.Get("/labels/:labelName/books", hand.GetBookByLabel)
}
