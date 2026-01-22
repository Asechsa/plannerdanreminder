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

type TaskHandler struct {
	db *pgxpool.Pool
}

func NewTaskHandler(db *pgxpool.Pool) *TaskHandler {
	return &TaskHandler{db: db}
}

// =======================
// GET TASKS (include urgency)
// =======================
func (h *TaskHandler) GetTasks(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	subCardID := c.Param("subcard_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := h.db.Query(ctx, `
		SELECT id, sub_card_id, title, deadline_at, note, status, urgency, created_at
		FROM tasks
		WHERE sub_card_id=$1
		ORDER BY 
			CASE urgency
				WHEN 'urgent' THEN 1
				WHEN 'overdue' THEN 2
				ELSE 3
			END,
			deadline_at ASC
	`, subCardID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed get tasks"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		_ = rows.Scan(&t.ID, &t.SubCardID, &t.Title, &t.DeadlineAt, &t.Note, &t.Status, &t.Urgency, &t.CreatedAt)
		tasks = append(tasks, t)
	}

	c.JSON(http.StatusOK, tasks)
}

// =======================
// CREATE TASK (auto reminder H-1 day & H-1 hour)
// =======================
func (h *TaskHandler) CreateTask(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	subCardID := c.Param("subcard_id")

	var req struct {
		Title      string `json:"title" binding:"required"`
		DeadlineAt string `json:"deadline_at" binding:"required"` // RFC3339
		Note       string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and deadline_at required"})
		return
	}

	deadline, err := time.Parse(time.RFC3339, req.DeadlineAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "deadline_at must be RFC3339 format"})
		return
	}

	taskID := uuid.New().String()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	var notePtr *string
	if req.Note != "" {
		notePtr = &req.Note
	}

	// Insert Task
	_, err = h.db.Exec(ctx, `
		INSERT INTO tasks (id, sub_card_id, title, deadline_at, note, status, urgency)
		VALUES ($1, $2, $3, $4, $5, 'pending', 'normal')
	`, taskID, subCardID, req.Title, deadline, notePtr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed create task"})
		return
	}

	// ✅ Auto create reminders (Option B): H-1 day & H-1 hour
	remindTimes := []time.Time{
		deadline.Add(-24 * time.Hour),
		deadline.Add(-1 * time.Hour),
	}

	for _, rt := range remindTimes {
		// insert hanya kalau remind_at masih di masa depan
		if rt.After(time.Now()) {
			remID := uuid.New().String()
			_, _ = h.db.Exec(ctx, `
				INSERT INTO reminders (id, task_id, remind_at, channel, status)
				VALUES ($1, $2, $3, 'email', 'pending')
			`, remID, taskID, rt)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          taskID,
		"sub_card_id": subCardID,
		"title":       req.Title,
		"deadline_at": deadline,
		"urgency":     "normal",
	})
}

// =======================
// UPDATE TASK STATUS
// =======================
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"` // pending/done
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status required"})
		return
	}

	if req.Status != "pending" && req.Status != "done" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status must be pending or done"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := h.db.Exec(ctx, `
		UPDATE tasks SET status=$1 WHERE id=$2
	`, req.Status, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed update task"})
		return
	}

	if cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// =======================
// DELETE TASK
// =======================
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	_, ok := getUserID(c)
	if !ok {
		return
	}

	id := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd, err := h.db.Exec(ctx, `
		DELETE FROM tasks WHERE id=$1
	`, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed delete task"})
		return
	}

	if cmd.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
