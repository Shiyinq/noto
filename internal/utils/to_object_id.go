package utils

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ToObjectID(id string) (primitive.ObjectID, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid id format")
	}

	return objectID, nil
}
