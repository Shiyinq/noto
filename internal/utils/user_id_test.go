package utils

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(*fiber.Ctx)
		expectedResult primitive.ObjectID
		expectedError  string
	}{
		{
			name: "Valid UserID",
			setupContext: func(c *fiber.Ctx) {
				c.Locals("userID", "507f1f77bcf86cd799439011")
			},
			expectedResult: func() primitive.ObjectID {
				id, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
				return id
			}(),
			expectedError: "",
		},
		{
			name: "UserID Not Found",
			setupContext: func(c *fiber.Ctx) {
			},
			expectedResult: primitive.NilObjectID,
			expectedError:  "user id not found",
		},
		{
			name: "Invalid UserID Format",
			setupContext: func(c *fiber.Ctx) {
				c.Locals("userID", "invalid-user-id")
			},
			expectedResult: primitive.NilObjectID,
			expectedError:  "invalid user id format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
			defer app.ReleaseCtx(ctx)

			tt.setupContext(ctx)

			userID, err := GetUserID(ctx)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, userID)
		})
	}
}
