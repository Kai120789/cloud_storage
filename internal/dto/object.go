package dto

// file or folder
type Object struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	UserID uint   `json:"user_id"`
}
