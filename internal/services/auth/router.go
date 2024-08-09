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

func AuthRouter(router fiber.Router) {
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

	router.Get("/auth/google", authHandler.HandleGoogleLogin)
	router.Get("/auth/google/callback", authHandler.HandleGoogleCallback)
}
