package service

import (
	"context"
	"encoding/json"
	"net/http"
	"noto/internal/services/auth/model"
	"noto/internal/services/auth/repository"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/oauth2"
)

type AuthService interface {
	HandleGoogleLogin() string
	HandleGoogleCallback(code string) (*model.AuthToken, error)
}

type authService struct {
	repo              repository.AuthRepository
	googleOauthConfig *oauth2.Config
	jwtSecret         []byte
}

func NewAuthService(repo repository.AuthRepository, googleOauthConfig *oauth2.Config, jwtSecret []byte) AuthService {
	return &authService{
		repo:              repo,
		googleOauthConfig: googleOauthConfig,
		jwtSecret:         jwtSecret,
	}
}

func (s *authService) HandleGoogleLogin() string {
	return s.googleOauthConfig.AuthCodeURL("state")
}

func (s *authService) HandleGoogleCallback(code string) (*model.AuthToken, error) {
	token, err := s.googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.getUserInfo(token.AccessToken)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:       userInfo["id"].(string),
		Email:    userInfo["email"].(string),
		Name:     userInfo["name"].(string),
		PhotoURL: userInfo["picture"].(string),
	}

	err = s.repo.FindOrCreateUser(context.Background(), user)
	if err != nil {
		return nil, err
	}

	jwtToken, err := s.createJWTToken(user)
	if err != nil {
		return nil, err
	}

	return &model.AuthToken{Token: jwtToken}, nil
}

func (s *authService) getUserInfo(accessToken string) (map[string]interface{}, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, err
}

func (s *authService) createJWTToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(s.jwtSecret)
}
