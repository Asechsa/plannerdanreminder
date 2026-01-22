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

type SubCardHandler struct {
	db *pgxpool.Pool
}

func NewSubCardHandler(db *pgxpool.Pool) *SubCardHandler {
	return &SubCardHandler{db: db}
}

func (h *SubCardHandler) GetSubCards(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	cardID := c.Param("card_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.db.Query(ctx, `
		SELECT id, card_id, title, created_at
		FROM sub_cards
		WHERE card_id=$1
		ORDER BY created_at DESC
	`, cardID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get subcards"})
		return
	}
	defer rows.Close()

	var subcards []models.SubCard
	for rows.Next() {
		var s models.SubCard
		_ = rows.Scan(&s.ID, &s.CardID, &s.Title, &s.CreatedAt)
		subcards = append(subcards, s)
	}

	c.JSON(http.StatusOK, subcards)
}

func (h *SubCardHandler) CreateSubCard(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	cardID := c.Param("card_id")

	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title required"})
		return
	}

	subID := uuid.New().String()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.db.Exec(ctx, `
		INSERT INTO sub_cards (id, card_id, title)
		VALUES ($1, $2, $3)
	`, subID, cardID, req.Title)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create subcard"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      subID,
		"card_id": cardID,
		"title":   req.Title,
	})
}

func (h *SubCardHandler) DeleteSubCard(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := h.db.Exec(ctx, `
		DELETE FROM sub_cards WHERE id=$1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed delete subcard"})
		return
	}

	if cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subcard not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SubCard deleted"})
}
