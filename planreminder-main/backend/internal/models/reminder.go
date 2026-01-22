package models

import "time"

type Reminder struct {
	ID        string     `json:"id"`
	TaskID    string     `json:"task_id"`
	RemindAt  time.Time  `json:"remind_at"`
	Channel   string     `json:"channel"` // push/email/both
	Status    string     `json:"status"`  // pending/sent/failed
	SentAt    *time.Time `json:"sent_at"`
	CreatedAt time.Time  `json:"created_at"`
}
