package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserID(c *fiber.Ctx) (primitive.ObjectID, error) {
	userId, ok := c.Locals("userID").(string)
	if !ok {
		return primitive.NilObjectID, errors.New("user id not found")
	}

	objUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid user id format")
	}

	return objUserId, nil
}
