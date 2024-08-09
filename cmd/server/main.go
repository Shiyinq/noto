package main

import (
	routes "noto/internal"

	"noto/internal/config"
	"noto/internal/middleware"

	_ "noto/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// @title		NOTO API
// @version		1.0
// @description	Noto API

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

// @host		localhost:8080
// @BasePath	/
func main() {
	config.LoadConfig()

	app := fiber.New(fiber.Config{
		EnablePrintRoutes: false,
	})

	app.Use(middleware.NewLogger())

	app.Get("/swagger/*", swagger.HandlerDefault)
	routes.SetupRoutes(app)

	app.Use(middleware.NotFoundHandler)

	app.Listen(":8080")
}
