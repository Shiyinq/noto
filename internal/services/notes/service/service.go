package service

import (
	model "noto/internal/services/notes/model"
	"noto/internal/services/notes/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteService interface {
	GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID, page int, limit int) (*model.PaginatedNoteResponse, error)
	CreateNote(note *model.NoteCreate) (*model.NoteCreate, error)
	UpdateNote(note *model.NoteUpdate) (*model.NoteResponse, error)
	DeleteNote(userId primitive.ObjectID, bookId primitive.ObjectID, noteId primitive.ObjectID) error
}

type NoteServiceImpl struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &NoteServiceImpl{noteRepo: noteRepo}
}

func (r *NoteServiceImpl) GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID, page int, limit int) (*model.PaginatedNoteResponse, error) {
	return r.noteRepo.GetNotes(userId, bookId, page, limit)
}

func (r *NoteServiceImpl) CreateNote(note *model.NoteCreate) (*model.NoteCreate, error) {
	return r.noteRepo.CreateNote(note)
}

func (r *NoteServiceImpl) UpdateNote(note *model.NoteUpdate) (*model.NoteResponse, error) {
	return r.noteRepo.UpdateNote(note)
}

func (r *NoteServiceImpl) DeleteNote(userId primitive.ObjectID, bookId primitive.ObjectID, noteId primitive.ObjectID) error {
	return r.noteRepo.DeleteNote(userId, bookId, noteId)
}
