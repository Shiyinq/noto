package main

import (
	routes "noto/internal"
	"os"

	"noto/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectMongoDB()

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: false,
	})

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		Output:     os.Stdout,
	}))

	routes.SetupRoutes(app)

	app.Listen(":8080")
}
