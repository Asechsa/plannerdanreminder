package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReminderService struct {
	db          *pgxpool.Pool
	email       *EmailService
	pollSeconds int
}

func NewReminderService(db *pgxpool.Pool, email *EmailService, pollSeconds int) *ReminderService {
	return &ReminderService{
		db:          db,
		email:       email,
		pollSeconds: pollSeconds,
	}
}

// Start scheduler loop
func (s *ReminderService) Start() {
	ticker := time.NewTicker(time.Duration(s.pollSeconds) * time.Second)
	defer ticker.Stop()

	log.Println("[ReminderService] scheduler started...")

	for range ticker.C {
		// ✅ update urgency otomatis
		_ = s.updateTaskUrgency()

		// ✅ process reminder pending
		err := s.processPendingReminders()
		if err != nil {
			log.Println("[ReminderService] error:", err)
		}
	}
}

// =======================
// UPDATE TASK URGENCY
// =======================
func (s *ReminderService) updateTaskUrgency() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// overdue
	_, err := s.db.Exec(ctx, `
		UPDATE tasks
		SET urgency='overdue'
		WHERE status != 'done'
		  AND deadline_at < NOW()
	`)
	if err != nil {
		return err
	}

	// urgent (deadline <= 1 hour)
	_, err = s.db.Exec(ctx, `
		UPDATE tasks
		SET urgency='urgent'
		WHERE status != 'done'
		  AND deadline_at BETWEEN NOW() AND NOW() + INTERVAL '1 hour'
	`)
	if err != nil {
		return err
	}

	// normal (deadline > 1 hour)
	_, err = s.db.Exec(ctx, `
		UPDATE tasks
		SET urgency='normal'
		WHERE status != 'done'
		  AND deadline_at > NOW() + INTERVAL '1 hour'
	`)
	return err
}

// =======================
// PROCESS PENDING REMINDERS
// =======================
func (s *ReminderService) processPendingReminders() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := s.db.Query(ctx, `
		SELECT r.id, r.task_id, r.channel
		FROM reminders r
		WHERE r.status='pending' AND r.remind_at <= NOW()
		ORDER BY r.remind_at ASC
		LIMIT 20
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var reminderID, taskID, channel string
		_ = rows.Scan(&reminderID, &taskID, &channel)

		err := s.sendReminder(reminderID, taskID, channel)
		if err != nil {
			log.Println("[ReminderService] failed send reminder:", err)
			_ = s.updateReminderStatus(reminderID, "failed")
		} else {
			_ = s.updateReminderStatus(reminderID, "sent")
		}
	}

	return nil
}

// =======================
// SEND REMINDER (EMAIL)
// =======================
func (s *ReminderService) sendReminder(reminderID, taskID, channel string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var title, email string
	var deadline time.Time

	err := s.db.QueryRow(ctx, `
		SELECT t.title, t.deadline_at, u.email
		FROM tasks t
		JOIN sub_cards sc ON t.sub_card_id = sc.id
		JOIN cards c ON sc.card_id = c.id
		JOIN users u ON c.user_id = u.id
		WHERE t.id=$1
	`, taskID).Scan(&title, &deadline, &email)

	if err != nil {
		return err
	}

	subject := "Reminder: " + title
	body := fmt.Sprintf("Hai!\n\nIni reminder untuk tugas:\n- %s\nDeadline: %s\n\nSemangat ya! ✅",
		title, deadline.Format("02 Jan 2006 15:04"),
	)

	// EMAIL (default)
	if channel == "email" || channel == "both" {
		if s.email == nil {
			return fmt.Errorf("email service not configured")
		}
		if err := s.email.Send(email, subject, body); err != nil {
			return err
		}
	}

	return nil
}

// =======================
// UPDATE REMINDER STATUS
// =======================
func (s *ReminderService) updateReminderStatus(reminderID string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.db.Exec(ctx, `
		UPDATE reminders
		SET status=$1,
		    sent_at=CASE WHEN $1='sent' THEN NOW() ELSE sent_at END
		WHERE id=$2
	`, status, reminderID)

	return err
}
