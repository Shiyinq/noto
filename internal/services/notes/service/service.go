package service

import (
	model "noto/internal/services/notes/model"
	"noto/internal/services/notes/repository"
)

type NoteService interface {
	GetAllNotes() ([]model.Note, error)
	GetNoteByID(id string) (*model.Note, error)
	CreateNote(note *model.Note) (*model.Note, error)
	UpdateNoteByID(id string, NewNote *model.Note) (*model.Note, error)
	DeleteNoteByID(id string) error
}

type NoteServiceImpl struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &NoteServiceImpl{noteRepo: noteRepo}
}

func (r *NoteServiceImpl) GetAllNotes() ([]model.Note, error) {
	var notes, err = r.noteRepo.GetAllNotes()
	return notes, err
}

func (r *NoteServiceImpl) GetNoteByID(id string) (*model.Note, error) {
	var note, err = r.noteRepo.GetNoteByID(id)
	return note, err
}

func (r *NoteServiceImpl) CreateNote(note *model.Note) (*model.Note, error) {
	var newNote, err = r.noteRepo.CreateNote(note)
	return newNote, err
}

func (r *NoteServiceImpl) UpdateNoteByID(id string, newNote *model.Note) (*model.Note, error) {
	var updated, err = r.noteRepo.UpdateNoteByID(id, newNote)
	return updated, err
}

func (r *NoteServiceImpl) DeleteNoteByID(id string) error {
	var err = r.noteRepo.DeleteNoteByID(id)
	return err
}
