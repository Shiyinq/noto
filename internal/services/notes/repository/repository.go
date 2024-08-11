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
	GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID) ([]model.NoteResponse, error)
	CreateNote(bookId string, note *model.NoteCreate) (*model.NoteCreate, error)
	UpdateNote(bookId string, noteId string, note *model.NoteUpdate) (*model.NoteResponse, error)
	DeleteNote(userId primitive.ObjectID, bookId primitive.ObjectID, noteId primitive.ObjectID) error
}

type NoteRepositoryImpl struct {
	notes *mongo.Collection
}

func NewNoteRepository() NoteRepository {
	return &NoteRepositoryImpl{notes: config.DB.Collection("notes")}
}

func (r *NoteRepositoryImpl) GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID) ([]model.NoteResponse, error) {
	var notes []model.NoteResponse

	filter := bson.M{"userId": userId, "bookId": bookId}
	cursor, err := r.notes.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &notes); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NoteRepositoryImpl) CreateNote(bookId string, note *model.NoteCreate) (*model.NoteCreate, error) {
	objectID, err := primitive.ObjectIDFromHex(bookId)
	if err != nil {
		return nil, errors.New("invalid bookId format")
	}

	note.BookId = objectID
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	newNote, err := r.notes.InsertOne(context.Background(), note)
	if err != nil {
		return nil, err
	}

	note.ID = newNote.InsertedID.(primitive.ObjectID)

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

func (r *NoteRepositoryImpl) DeleteNote(userId primitive.ObjectID, bookId primitive.ObjectID, noteId primitive.ObjectID) error {
	filter := bson.M{"_id": noteId, "userId": userId, "bookId": bookId}

	deleted, err := r.notes.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	if deleted.DeletedCount == 0 {
		return errors.New("note not found or not deleted")
	}

	return nil
}
