package handlers

import (
	"context"
	"net/http"
	"time"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CardHandler struct {
	db *pgxpool.Pool
}

func NewCardHandler(db *pgxpool.Pool) *CardHandler {
	return &CardHandler{db: db}
}

func (h *CardHandler) GetCards(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.db.Query(ctx, `
		SELECT id, user_id, title, created_at
		FROM cards
		WHERE user_id=$1
		ORDER BY created_at DESC
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get cards"})
		return
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var card models.Card
		_ = rows.Scan(&card.ID, &card.UserID, &card.Title, &card.CreatedAt)
		cards = append(cards, card)
	}

	c.JSON(http.StatusOK, cards)
}

func (h *CardHandler) CreateCard(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}

	cardID := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.db.Exec(ctx, `
		INSERT INTO cards (id, user_id, title)
		VALUES ($1, $2, $3)
	`, cardID, userID, req.Title)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create card"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      cardID,
		"user_id": userID,
		"title":   req.Title,
	})
}

func (h *CardHandler) DeleteCard(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := h.db.Exec(ctx, `
		DELETE FROM cards WHERE id=$1 AND user_id=$2
	`, id, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed delete card"})
		return
	}

	if cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Card not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card deleted"})
}
