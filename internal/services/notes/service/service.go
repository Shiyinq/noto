package service

import (
	model "noto/internal/services/notes/model"
	"noto/internal/services/notes/repository"
)

type NoteService interface {
	GetAllNotes(bookId string) ([]model.NoteResponse, error)
	CreateNote(note *model.Note) (*model.Note, error)
	UpdateNote(bookId string, noteId string, note *model.NoteUpdate) (*model.NoteResponse, error)
	DeleteNote(bookId string, noteId string) error
}

type NoteServiceImpl struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &NoteServiceImpl{noteRepo: noteRepo}
}

func (r *NoteServiceImpl) GetAllNotes(bookId string) ([]model.NoteResponse, error) {
	var notes, err = r.noteRepo.GetAllNotes(bookId)
	return notes, err
}

func (r *NoteServiceImpl) CreateNote(note *model.Note) (*model.Note, error) {
	var newNote, err = r.noteRepo.CreateNote(note)
	return newNote, err
}

func (r *NoteServiceImpl) UpdateNote(bookId string, noteId string, note *model.NoteUpdate) (*model.NoteResponse, error) {
	var updated, err = r.noteRepo.UpdateNote(bookId, noteId, note)
	return updated, err
}

func (r *NoteServiceImpl) DeleteNote(bookId string, noteId string) error {
	var err = r.noteRepo.DeleteNote(bookId, noteId)
	return err
}
