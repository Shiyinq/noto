package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestToObjectID(t *testing.T) {
	testCases := []struct {
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

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			objectID, err := ToObjectID(testCase.request)

			if testCase.expectedError != "" {
				assert.EqualError(t, err, testCase.expectedError)
			}

			assert.Equal(t, objectID, testCase.expectedResult)
		})
	}
}
