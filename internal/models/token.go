package models

type UserToken struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}
