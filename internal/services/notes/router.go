package notes_router

import (
	"noto/internal/services/notes/handler"
	"noto/internal/services/notes/repository"
	"noto/internal/services/notes/service"

	"github.com/gofiber/fiber/v2"
)

func NotesRouter(app *fiber.App) {
	var repo = repository.NewNoteRepository()
	var serv = service.NewNoteService(repo)
	var hand = handler.NewNoteHandler(serv)

	app.Get("/books/:bookId/notes", hand.GetNotes)
	app.Post("/books/:bookId/notes", hand.CreateNote)
	app.Patch("/books/:bookId/notes/:noteId", hand.UpdateNote)
	app.Delete("/books/:bookId/notes/:noteId", hand.DeleteNote)
}
