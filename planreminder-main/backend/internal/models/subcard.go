package models

import "time"

type SubCard struct {
	ID        string    `json:"id"`
	CardID    string    `json:"card_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
