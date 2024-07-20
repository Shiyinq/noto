package handler

import (
	model "noto/internal/services/notes/model"
	service "noto/internal/services/notes/service"

	"github.com/gofiber/fiber/v2"
)

type NoteHandler interface {
	GetNotes(c *fiber.Ctx) error
	GetNote(c *fiber.Ctx) error
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

func (s *NoteHandlerImpl) GetNotes(c *fiber.Ctx) error {
	notes, err := s.noteService.GetAllNotes()
	if err != nil {
		return err
	}
	return c.JSON(notes)
}

func (s *NoteHandlerImpl) GetNote(c *fiber.Ctx) error {
	id := c.Params("id")
	note, err := s.noteService.GetNoteByID(id)
	if err != nil {
		return c.Status(404).SendString("Note not found")
	}
	return c.JSON(note)
}

func (s *NoteHandlerImpl) CreateNote(c *fiber.Ctx) error {
	note := new(model.Note)
	if err := c.BodyParser(note); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	newNote, _ := s.noteService.CreateNote(note)
	return c.JSON(newNote)
}

func (s *NoteHandlerImpl) UpdateNote(c *fiber.Ctx) error {
	id := c.Params("id")
	newNote := new(model.Note)
	if err := c.BodyParser(newNote); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	updatedNote, err := s.noteService.UpdateNoteByID(id, newNote)
	if err != nil {
		return c.Status(404).SendString("Note not found")
	}
	return c.JSON(updatedNote)
}

func (s *NoteHandlerImpl) DeleteNote(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := s.noteService.DeleteNoteByID(id); err != nil {
		return c.Status(404).SendString("Note not found")
	}
	return c.SendString("Note deleted")
}
