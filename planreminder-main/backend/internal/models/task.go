package models

import "time"

type Task struct {
	ID         string    `json:"id"`
	SubCardID  string    `json:"sub_card_id"`
	Title      string    `json:"title"`
	DeadlineAt time.Time `json:"deadline_at"`
	Note       *string   `json:"note"`
	Status     string    `json:"status"`
	Urgency    string    `json:"urgency"`
	CreatedAt  time.Time `json:"created_at"`
}
