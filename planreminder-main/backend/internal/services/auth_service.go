package services

import (
	"context"
	"errors"
	"time"

	"backend/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/api/idtoken"
)

type AuthService struct {
	cfg config.Config
}

func NewAuthService(cfg config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

type GoogleUserPayload struct {
	GoogleID string
	Email    string
	Name     string
}

func (s *AuthService) VerifyGoogleToken(idToken string) (*GoogleUserPayload, error) {
	if idToken == "" {
		return nil, errors.New("missing id_token")
	}

	// audience optional: kalau ingin strict, isi dengan cfg.GoogleClientID
	payload, err := idtoken.Validate(context.Background(), idToken, "")
	if err != nil {
		return nil, errors.New("invalid google token")
	}

	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	if googleID == "" || email == "" {
		return nil, errors.New("missing required google payload fields")
	}

	return &GoogleUserPayload{
		GoogleID: googleID,
		Email:    email,
		Name:     name,
	}, nil
}

func (s *AuthService) GenerateJWT(userID string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWTSecret))
}
