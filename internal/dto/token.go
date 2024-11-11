package dto

type UserToken struct {
	UserID       uint   `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}
