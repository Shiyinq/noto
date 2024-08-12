package handler

import (
	_ "noto/internal/common"
	"noto/internal/services/notes/model"
	service "noto/internal/services/notes/service"
	"noto/internal/utils"

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
// @Security 	BearerAuth
// @Produce		json
// @Param 		Authorization header string false "Bearer token"
// @Param		bookId path string true "Book ID"
// @Param		page		query		int		false	"Page number for pagination"	minimum(1)
// @Param		limit		query		int		false	"Number of items per page"	minimum(1)
// @Success		200		{object}	model.PaginatedNoteResponse
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books/{bookId}/notes [get]
func (s *NoteHandlerImpl) GetNotes(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	notes, err := s.noteService.GetNotes(userId, bookId, page, limit)
	if err != nil {
		return utils.ErrorInternalServer(c, err.Error())
	}

	return c.JSON(notes)
}

// CreateNoe
// @Summary		Create a new note
// @Description	Create a new note
// @Tags		Notes
// @Security 	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		Authorization header string false "Bearer token"
// @Param		bookId path string true "Book ID"
// @Param		book	body		model.NoteCreateSwagger	true	"Note to create"
// @Success		201		{object}	model.NoteCreate
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books/{bookId}/notes [post]
func (s *NoteHandlerImpl) CreateNote(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	note := new(model.NoteCreate)
	if err := c.BodyParser(note); err != nil {
		return utils.ErrorBadRequest(c, "failed to parse json: "+err.Error())
	}

	note.UserId = userId
	note.BookId = bookId
	newNote, err := s.noteService.CreateNote(note)
	if err != nil {
		return utils.ErrorInternalServer(c, "failed to create note: "+err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newNote)
}

// UpdateNote
// @Summary		Update note
// @Description	Update note
// @Tags		Notes
// @Security 	BearerAuth
// @Accept		json
// @Produce		json
// @Param 		Authorization header string false "Bearer token"
// @Param		bookId path string true "Book ID"
// @Param		noteId path string true "Note ID"
// @Param		book	body		model.NoteUpdateSwagger	true	"Note to update"
// @Success		200		{object}	model.NoteResponse
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books/{bookId}/notes/{noteId} [patch]
func (s *NoteHandlerImpl) UpdateNote(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	noteId, err := utils.ToObjectID(c.Params("noteId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	note := new(model.NoteUpdate)
	if err := c.BodyParser(note); err != nil {
		return utils.ErrorBadRequest(c, "failed to parse json: "+err.Error())
	}

	note.ID = noteId
	note.UserId = userId
	note.BookId = bookId
	updatedNote, err := s.noteService.UpdateNote(note)
	if err != nil {
		return utils.ErrorInternalServer(c, "failed to update note: "+err.Error())
	}

	return c.JSON(updatedNote)
}

// DeleteNote
// @Summary		Delete note
// @Description	Deelete note
// @Tags		Notes
// @Security 	BearerAuth
// @Produce		json
// @Param 		Authorization header string false "Bearer token"
// @Param		bookId path string true "Book ID"
// @Param		noteId path string true "Note ID"
// @Success		200		{object} 	interface{}
// @Failure     400    	{object}   	common.ErrorResponse
// @Failure     401     {object}    common.ErrorResponse
// @Failure     500     {object}    common.ErrorResponse
// @Router		/api/books/{bookId}/notes/{noteId} [delete]
func (s *NoteHandlerImpl) DeleteNote(c *fiber.Ctx) error {
	userId, err := utils.GetUserID(c)
	if err != nil {
		return utils.ErrorUnauthorized(c, err.Error())
	}

	bookId, err := utils.ToObjectID(c.Params("bookId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	noteId, err := utils.ToObjectID(c.Params("noteId"))
	if err != nil {
		return utils.ErrorBadRequest(c, err.Error())
	}

	if err := s.noteService.DeleteNote(userId, bookId, noteId); err != nil {
		return utils.ErrorInternalServer(c, "failed to delete note: "+err.Error())
	}

	return c.JSON(fiber.Map{
		"success": "note deleted",
	})
}
