package repository

import (
	"errors"
	model "noto/internal/services/notes/model"
)

type NoteRepository interface {
	GetAllNotes() ([]model.Note, error)
	GetNoteByID(id string) (*model.Note, error)
	CreateNote(Note *model.Note) (*model.Note, error)
	UpdateNoteByID(id string, NewNote *model.Note) (*model.Note, error)
	DeleteNoteByID(id string) error
}

type NoteRepositoryImpl struct {
	notes []model.Note
}

func NewNoteRepository() NoteRepository {
	return &NoteRepositoryImpl{notes: []model.Note{}}
}

func (r *NoteRepositoryImpl) GetAllNotes() ([]model.Note, error) {
	return r.notes, nil
}

func (r *NoteRepositoryImpl) GetNoteByID(id string) (*model.Note, error) {
	for _, note := range r.notes {
		if note.ID == id {
			return &note, nil
		}
	}
	return nil, errors.New("note not found")
}

func (r *NoteRepositoryImpl) CreateNote(note *model.Note) (*model.Note, error) {
	var _ = append(r.notes, *note)
	return note, nil
}

func (r *NoteRepositoryImpl) UpdateNoteByID(id string, newNote *model.Note) (*model.Note, error) {
	for i, note := range r.notes {
		if note.ID == id {
			r.notes[i] = *newNote
			return newNote, nil
		}
	}
	return nil, errors.New("note not found")
}

func (r *NoteRepositoryImpl) DeleteNoteByID(id string) error {
	for i, note := range r.notes {
		if note.ID == id {
			var _ = append(r.notes[:i], r.notes[i+1:]...)
			return nil
		}
	}
	return errors.New("note not found")
}
