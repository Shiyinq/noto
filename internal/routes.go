package routes

import (
	auth_router "noto/internal/services/auth"
	books_router "noto/internal/services/books"
	labels_router "noto/internal/services/labels"
	notes_router "noto/internal/services/notes"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	auth_router.AuthRouter(app)
	books_router.BooksRouter(app)
	notes_router.NotesRouter(app)
	labels_router.LabelsRouter(app)
}
