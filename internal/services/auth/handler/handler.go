package handler

import (
	_ "noto/internal/common"
	_ "noto/internal/services/auth/model"
	"noto/internal/services/auth/service"
	"noto/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// HandleGoogleLogin godoc
// @Summary Initiate Google OAuth login
// @Description Redirects the user to Google's OAuth consent screen
// @Tags Auth
// @Produce json
// @Success 	302 	{string} string "Redirect to Google's OAuth consent screen"
// @Router /auth/google [get]
func (h *AuthHandler) HandleGoogleLogin(c *fiber.Ctx) error {
	url := h.authService.HandleGoogleLogin()
	return c.Redirect(url)
}

// HandleGoogleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Processes the OAuth code returned by Google and returns a token
// @Tags Auth
// @Accept json
// @Produce json
// @Param code query string true "The authorization code returned by Google"
// @Success 	200 	{object} 	model.AuthToken
// @Failure     500     {object}    common.ErrorResponse
// @Router /auth/google/callback [get]
func (h *AuthHandler) HandleGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		return utils.ErrorInternalServer(c, err.Error())
	}

	return c.JSON(token)
}
