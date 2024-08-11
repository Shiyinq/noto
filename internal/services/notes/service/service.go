package service

import (
	model "noto/internal/services/notes/model"
	"noto/internal/services/notes/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NoteService interface {
	GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID) ([]model.NoteResponse, error)
	CreateNote(bookId string, note *model.NoteCreate) (*model.NoteCreate, error)
	UpdateNote(bookId string, noteId string, note *model.NoteUpdate) (*model.NoteResponse, error)
	DeleteNote(userId primitive.ObjectID, bookId primitive.ObjectID, noteId primitive.ObjectID) error
}

type NoteServiceImpl struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &NoteServiceImpl{noteRepo: noteRepo}
}

func (r *NoteServiceImpl) GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID) ([]model.NoteResponse, error) {
	var notes, err = r.noteRepo.GetNotes(userId, bookId)
	return notes, err
}

func (r *NoteServiceImpl) CreateNote(bookId string, note *model.NoteCreate) (*model.NoteCreate, error) {
	var newNote, err = r.noteRepo.CreateNote(bookId, note)
	return newNote, err
}

func (r *NoteServiceImpl) UpdateNote(bookId string, noteId string, note *model.NoteUpdate) (*model.NoteResponse, error) {
	var updated, err = r.noteRepo.UpdateNote(bookId, noteId, note)
	return updated, err
}

func (r *NoteServiceImpl) DeleteNote(userId primitive.ObjectID, bookId primitive.ObjectID, noteId primitive.ObjectID) error {
	var err = r.noteRepo.DeleteNote(userId, bookId, noteId)
	return err
}
