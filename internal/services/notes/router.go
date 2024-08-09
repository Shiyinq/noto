package notes_router

import (
	"noto/internal/services/notes/handler"
	"noto/internal/services/notes/repository"
	"noto/internal/services/notes/service"

	"github.com/gofiber/fiber/v2"
)

func NotesRouter(router fiber.Router) {
	var repo = repository.NewNoteRepository()
	var serv = service.NewNoteService(repo)
	var hand = handler.NewNoteHandler(serv)

	router.Get("/books/:bookId/notes", hand.GetNotes)
	router.Post("/books/:bookId/notes", hand.CreateNote)
	router.Patch("/books/:bookId/notes/:noteId", hand.UpdateNote)
	router.Delete("/books/:bookId/notes/:noteId", hand.DeleteNote)
}
