package repository

import (
	"context"
	"noto/internal/services/auth/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthRepository interface {
	FindOrCreateUser(ctx context.Context, user *model.User) (*model.User, error)
}

type AuthRepositoryImpl struct {
	users *mongo.Collection
}

func NewAuthRepository(db *mongo.Database) AuthRepository {
	return &AuthRepositoryImpl{users: db.Collection("users")}
}

func (r *AuthRepositoryImpl) FindOrCreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	now := time.Now()
	filter := bson.M{"email": user.Email}
	update := bson.M{
		"$set": bson.M{
			"name":      user.Name,
			"photoUrl":  user.PhotoURL,
			"updatedAt": now,
		},
		"$setOnInsert": bson.M{
			"createdAt": now,
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := r.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, err
	}

	var currentUser model.User
	err = r.users.FindOne(ctx, filter).Decode(&currentUser)
	if err != nil {
		return nil, err
	}

	return &currentUser, err
}
