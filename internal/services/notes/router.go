package notes_router

import (
	"noto/internal/services/notes/handler"
	"noto/internal/services/notes/repository"
	"noto/internal/services/notes/service"

	"github.com/gofiber/fiber/v2"
)

var repo = repository.NewNoteRepository()
var serv = service.NewNoteService(repo)
var hand = handler.NewNoteHandler(serv)

func NotesRouter(app *fiber.App) {
	app.Get("/notes", hand.GetNotes)
	app.Get("/notes/:id", hand.GetNote)
	app.Post("/notes", hand.CreateNote)
	app.Put("/notes/:id", hand.UpdateNote)
	app.Delete("/notes/:id", hand.DeleteNote)
}
