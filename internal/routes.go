package routes

import (
	"noto/internal/middleware"
	auth_router "noto/internal/services/auth"
	books_router "noto/internal/services/books"
	labels_router "noto/internal/services/labels"
	notes_router "noto/internal/services/notes"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	prefix := "/api"

	normal := app.Group("")
	auth_router.AuthRouter(normal)

	protected := app.Group(prefix, middleware.Protected)
	books_router.BooksRouter(protected)
	notes_router.NotesRouter(protected)
	labels_router.LabelsRouter(protected)
}
