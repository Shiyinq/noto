package handler

import (
	"noto/internal/services/notes/model"
	service "noto/internal/services/notes/service"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler interface {
	GetNotes(c *fiber.Ctx) error
	CreateNote(c *fiber.Ctx) error
	UpdateNote(c *fiber.Ctx) error
	DeleteNote(c *fiber.Ctx) error
}

type NoteHandlerImpl struct {
	noteService service.NoteService
}

func NewNoteHandler(noteService service.NoteService) NoteHandler {
	return &NoteHandlerImpl{noteService: noteService}
}

// GetNotes
// @Summary		Get notes by book id
// @Description	Get notes by book id
// @Tags		Notes
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Success		200		{object}	[]model.NoteResponse
// @Router		/books/{bookId}/notes [get]
func (s *NoteHandlerImpl) GetNotes(c *fiber.Ctx) error {
	bookId := c.Params("bookId")
	notes, err := s.noteService.GetAllNotes(bookId)
	if err != nil {
		return err
	}
	return c.JSON(notes)
}

// CreateNoe
// @Summary		Create a new note
// @Description	Create a new note
// @Tags		Notes
// @Accept		json
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Param		book	body		model.NoteCreateSwagger	true	"Note to create"
// @Success		201		{object}	model.NoteCreate
// @Router		/books/{bookId}/notes [post]
func (s *NoteHandlerImpl) CreateNote(c *fiber.Ctx) error {
	bookId := c.Params("bookId")
	note := new(model.NoteCreate)
	if err := c.BodyParser(note); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	newNote, _ := s.noteService.CreateNote(bookId, note)
	return c.Status(fiber.StatusCreated).JSON(newNote)
}

// UpdateNote
// @Summary		Update note
// @Description	Update note
// @Tags		Notes
// @Accept		json
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Param		noteId path string true "Note ID"
// @Param		book	body		model.NoteUpdate	true	"Note to update"
// @Success		201		{object}	model.NoteResponse
// @Router		/books/{bookId}/notes/{noteId} [patch]
func (s *NoteHandlerImpl) UpdateNote(c *fiber.Ctx) error {
	bookId := c.Params("bookId")
	noteId := c.Params("noteId")
	note := new(model.NoteUpdate)
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	updatedNote, err := s.noteService.UpdateNote(bookId, noteId, note)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(updatedNote)
}

// DeleteNote
// @Summary		Delete note
// @Description	Deelete note
// @Tags		Notes
// @Accept		json
// @Produce		json
// @Param		bookId path string true "Book ID"
// @Param		noteId path string true "Note ID"
// @Success		201		{object} interface{}
// @Router		/books/{bookId}/notes/{noteId} [delete]
func (s *NoteHandlerImpl) DeleteNote(c *fiber.Ctx) error {
	bookId := c.Params("bookId")
	noteId := c.Params("noteId")

	if err := s.noteService.DeleteNote(bookId, noteId); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": "note deleted",
	})
}
