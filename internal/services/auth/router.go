package auth

import (
	"noto/internal/config"
	"noto/internal/services/auth/handler"
	"noto/internal/services/auth/repository"
	"noto/internal/services/auth/service"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func AuthRouter(app *fiber.App) {
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     config.GoogleClientID,
		ClientSecret: config.GoogleClientSecret,

		Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	authRepo := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepo, googleOauthConfig, config.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	auth := app.Group("/auth")
	auth.Get("/google", authHandler.HandleGoogleLogin)
	auth.Get("/google/callback", authHandler.HandleGoogleCallback)
}
