package models

import "time"


type User struct {
	ID    string `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}


