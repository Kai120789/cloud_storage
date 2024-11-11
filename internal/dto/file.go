package dto

type File struct {
	FileName string `json:"file_name"`
	Path     string `json:"path"`
	UserID   uint   `json:"user_id"`
}
