package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/api/idtoken"
)

type AuthHandler struct {
	db  *pgxpool.Pool
	cfg config.Config
}

func NewAuthHandler(db *pgxpool.Pool, cfg config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type googleLoginReq struct {
	IDToken string `json:"id_token" binding:"required"`
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	var req googleLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_token required"})
		return
	}

	// Verify Google ID Token
	payload, err := idtoken.Validate(context.Background(), req.IDToken, "")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid google token"})
		return
	}

	googleID := payload.Subject
	email, _ := payload.Claims["email"].(string)
	name, _ := payload.Claims["name"].(string)

	if googleID == "" || email == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Google token missing required fields"})
		return
	}

	// Upsert user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var userID string
	err = h.db.QueryRow(ctx, `
		SELECT id FROM users WHERE google_id=$1
	`, googleID).Scan(&userID)

	if err != nil {
		// create new user
		userID = uuid.New().String()
		_, err2 := h.db.Exec(ctx, `
			INSERT INTO users (id, google_id, email, name)
			VALUES ($1, $2, $3, $4)
		`, userID, googleID, email, name)

		if err2 != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create user"})
			return
		}
	}

	// Create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 days
		"iat":     time.Now().Unix(),
	})

	signedToken, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
		"user": gin.H{
			"id":    userID,
			"email": email,
			"name":  name,
		},
	})
}
