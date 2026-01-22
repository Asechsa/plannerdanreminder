package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PushHandler struct {
	db *pgxpool.Pool
}

func NewPushHandler(db *pgxpool.Pool) *PushHandler {
	return &PushHandler{db: db}
}

func (h *PushHandler) Subscribe(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req struct {
		Endpoint string `json:"endpoint" binding:"required"`
		P256dh   string `json:"p256dh" binding:"required"`
		Auth     string `json:"auth" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint, p256dh, auth required"})
		return
	}

	subID := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.db.Exec(ctx, `
		INSERT INTO push_subscriptions (id, user_id, endpoint, p256dh, auth)
		VALUES ($1, $2, $3, $4, $5)
	`, subID, userID, req.Endpoint, req.P256dh, req.Auth)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed save subscription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Subscribed"})
}

func (h *PushHandler) Unsubscribe(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		return
	}

	var req struct {
		Endpoint string `json:"endpoint" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "endpoint required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := h.db.Exec(ctx, `
		DELETE FROM push_subscriptions WHERE user_id=$1 AND endpoint=$2
	`, userID, req.Endpoint)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed unsubscribe"})
		return
	}

	if cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed"})
}
