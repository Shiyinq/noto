package main

import (
	routes "noto/internal"

	"noto/internal/config"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectMongoDB()

	app := fiber.New()

	routes.SetupRoutes(app)

	app.Listen(":8080")
}
