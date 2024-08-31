package repository

import (
	"context"
	"errors"
	model "noto/internal/services/notes/model"
	"noto/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NoteRepository interface {
	GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID, page int, limit int) (*model.PaginatedNoteResponse, error)
	CreateNote(note *model.NoteCreate) (*model.NoteCreate, error)
	UpdateNote(note *model.NoteUpdate) (*model.NoteResponse, error)
	DeleteNote(userId primitive.ObjectID, bookId primitive.ObjectID, noteId primitive.ObjectID) error
}

type NoteRepositoryImpl struct {
	notes *mongo.Collection
}

func NewNoteRepository(db *mongo.Database) NoteRepository {
	return &NoteRepositoryImpl{notes: db.Collection("notes")}
}

func (r *NoteRepositoryImpl) GetNotes(userId primitive.ObjectID, bookId primitive.ObjectID, page int, limit int) (*model.PaginatedNoteResponse, error) {
	var notes []model.PaginatedNoteResponse

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"userId": userId, "bookId": bookId}}},
		{{Key: "$sort", Value: bson.M{"createdAt": -1}}},
		{{Key: "$facet", Value: utils.PaginationAggregate(page, limit)}},
		{{Key: "$unwind", Value: "$metadata"}},
	}

	cursor, err := r.notes.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.Background(), &notes); err != nil {
		return nil, err
	}

	if len(notes) == 0 {
		return &model.PaginatedNoteResponse{
			Data:     []model.NoteResponse{},
			Metadata: model.PaginationMetadata{},
		}, nil
	}

	return &notes[0], nil
}

func (r *NoteRepositoryImpl) CreateNote(note *model.NoteCreate) (*model.NoteCreate, error) {
	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	newNote, err := r.notes.InsertOne(context.Background(), note)
	if err != nil {
		return nil, err
	}

	note.ID = newNote.InsertedID.(primitive.ObjectID)

	return note, nil
}

func (r *NoteRepositoryImpl) UpdateNote(note *model.NoteUpdate) (*model.NoteResponse, error) {
	filter := bson.M{"_id": note.ID, "userId": note.UserId, "bookId": note.BookId}
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
	err := result.Decode(&updatedNote)
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
