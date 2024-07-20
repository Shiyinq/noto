package routes

import (
	books_router "noto/internal/services/books"
	labels_router "noto/internal/services/labels"
	notes_router "noto/internal/services/notes"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	books_router.BooksRouter(app)
	notes_router.NotesRouter(app)
	labels_router.LabelsRouter(app)
}
