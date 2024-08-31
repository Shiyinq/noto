package labels_router

import (
	"noto/internal/config"
	"noto/internal/services/labels/handler"
	"noto/internal/services/labels/repository"
	"noto/internal/services/labels/service"

	"github.com/gofiber/fiber/v2"
)

func LabelsRouter(router fiber.Router) {
	var repo = repository.NewLabelRepository(config.DB)
	var serv = service.NewLabelService(repo)
	var hand = handler.NewLabelHandler(serv)

	router.Post("/labels", hand.CreateLabel)
	router.Get("/labels", hand.GetLabels)
	router.Delete("/labels/:labelId", hand.DeleteLabel)
	router.Post("/books/:bookId/labels", hand.AddBookLabel)
	router.Delete("/books/:bookId/labels", hand.DeleteBookLabel)
	router.Get("/labels/:labelName/books", hand.GetBookByLabel)
}
