package utils

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserID(c *fiber.Ctx) (primitive.ObjectID, error) {
	userId, ok := c.Locals("userID").(string)
	if !ok {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusUnauthorized, "User ID not found")
	}

	objUserId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return primitive.NilObjectID, fiber.NewError(fiber.StatusBadRequest, "Invalid User ID format")
	}

	return objUserId, nil
}
