package models

import "time"

// file or folder
type Object struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
