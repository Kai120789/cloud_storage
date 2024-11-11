package models

import "time"

type File struct {
	ID        uint      `json:"id"`
	FileName  string    `json:"file_name"`
	Path      string    `json:"path"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
