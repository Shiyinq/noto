package repository

import (
	"context"
	"errors"
	"noto/internal/config"
	model "noto/internal/services/notes/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoteRepository interface {
	GetAllNotes(bookId string) ([]model.NoteResponse, error)
	CreateNote(Note *model.Note) (*model.Note, error)
	UpdateNote(bookId string, noteId string, note *model.NoteUpdate) (*model.NoteResponse, error)
	DeleteNote(bookId string, noteId string) error
}

type NoteRepositoryImpl struct {
	notes *mongo.Collection
}

func NewNoteRepository() NoteRepository {
	return &NoteRepositoryImpl{notes: config.DB.Collection("notes")}
}

func (r *NoteRepositoryImpl) GetAllNotes(bookId string) ([]model.NoteResponse, error) {
	var notes []model.NoteResponse
	objectID, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return nil, errors.New("invalid bookId format")
	}

	filter := bson.M{"bookId": objectID}
	cursor, err := r.notes.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NoteRepositoryImpl) CreateNote(note *model.Note) (*model.Note, error) {
	var idErr error
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()
	note.BookId, idErr = primitive.ObjectIDFromHex(note.BookId.Hex())
	if idErr != nil {
		return nil, errors.New("invalid ID format")
	}

	_, err := r.notes.InsertOne(context.Background(), note)
	if err != nil {
		return nil, err
	}

	return note, nil
}

func (r *NoteRepositoryImpl) UpdateNote(bookId string, noteId string, note *model.NoteUpdate) (*model.NoteResponse, error) {
	bookObjectID, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return nil, errors.New("invalid bookId format")
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteId)
	if err != nil {
		return nil, errors.New("invalid noteId format")
	}

	filter := bson.M{"_id": noteObjectID, "bookId": bookObjectID}
	update := bson.M{
		"$set": bson.M{
			"text":      note.Text,
			"updatedAt": time.Now(),
		},
	}

	result := r.notes.FindOneAndUpdate(context.Background(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))

	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return nil, errors.New("note not found")
		}
		return nil, result.Err()
	}

	var updatedNote model.NoteResponse
	err = result.Decode(&updatedNote)
	if err != nil {
		return nil, err
	}

	return &updatedNote, nil
}

func (r *NoteRepositoryImpl) DeleteNote(bookId string, noteId string) error {
	bookObjectID, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return errors.New("invalid bookId format")
	}

	noteObjectID, err := primitive.ObjectIDFromHex(noteId)
	if err != nil {
		return errors.New("invalid noteId format")
	}

	filter := bson.M{"_id": noteObjectID, "bookId": bookObjectID}

	deleted, err := r.notes.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	if deleted.DeletedCount == 0 {
		return errors.New("note not found or not deleted")
	}

	return nil
}
