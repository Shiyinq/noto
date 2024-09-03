package service

import (
	"context"
	"noto/internal/config"
	"noto/internal/services/books/model"
	"noto/internal/services/books/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	now    time.Time
	userId primitive.ObjectID
	bookId primitive.ObjectID
	serv   BookService
	repo   repository.BookRepository
)

func TestMain(m *testing.M) {
	config.LoadConfig()
	repo = repository.NewBookRepository(config.DB)
	serv = NewBookService(repo)

	now = time.Now()
	userId = primitive.NewObjectID()
	bookId = primitive.NewObjectID()

	m.Run()

	cleanupDatabase()
}

func cleanupDatabase() {
	collection := config.DB.Collection("books")
	filter := bson.M{
		"$or": []bson.M{
			{"_id": bookId},
			{"userId": userId},
		},
	}
	_, err := collection.DeleteMany(context.Background(), filter)
	if err != nil {
		panic(err)
	}
}

func TestCreateBook(t *testing.T) {
	newBook := model.BookCreate{
		ID:         bookId,
		UserId:     userId,
		Title:      "Test Book",
		CreatedAt:  now,
		UpdatedAt:  now,
		IsArchived: false,
	}

	res, err := serv.CreateBook(&newBook)

	require.NoError(t, err, "Failed to create book")
	assert.Equal(t, &newBook, res, "Data result should match expected")
}

func TestGetBooks(t *testing.T) {
	res, err := serv.GetBooks(userId, false, 1, 10)

	require.NoError(t, err, "Failed to get books")
	assert.NotNil(t, res, "Data should not nil")
}

func TestGetBook(t *testing.T) {
	res, err := serv.GetBook(userId, bookId)

	if err != nil {
		if err.Error() != "book not found" {
			require.NoError(t, err, "Failed to get book")
		}
		assert.Equal(t, "book not found", err.Error())
	}

	if res != nil {
		assert.Equal(t, bookId.Hex(), res.ID, "BookID should match expected")
	}
}

func TestUpdateBook(t *testing.T) {
	title := "Update Book"
	updated := model.BookUpdate{
		ID:     bookId,
		UserId: userId,
		Title:  title,
	}

	res, err := serv.UpdateBook(&updated)
	require.NoError(t, err, "Failed to update book")
	assert.Equal(t, title, res.Title, "Title result should match expected")
}

func TestArchiveBook(t *testing.T) {
	archive := true
	archived := model.ArchiveBook{
		ID:         bookId,
		UserId:     userId,
		IsArchived: archive,
	}

	res, err := serv.ArchiveBook(&archived)
	require.NoError(t, err, "Failed to archive book")
	assert.Equal(t, archive, res.IsArchived, "IsArchived result should match expected")
}
