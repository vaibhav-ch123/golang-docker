package models

import "time"

type Todo struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"userID" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	PendingAt   time.Time `json:"pendingAt" db:"pending_at"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}
