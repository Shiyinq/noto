package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToObjectID(t *testing.T) {
	test := []struct {
		name           string
		request        string
		expectedResult primitive.ObjectID
		expectedError  string
	}{
		{
			name:    "Valid ID",
			request: "507f1f77bcf86cd799439011",
			expectedResult: func() primitive.ObjectID {
				id, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				return id
			}(),
			expectedError: "",
		},
		{
			name:           "Invalid ID",
			request:        "123",
			expectedResult: primitive.NilObjectID,
			expectedError:  "invalid id format",
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			objectID, err := ToObjectID(tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			}

			assert.Equal(t, objectID, tt.expectedResult)
		})
	}
}
