package model

type TokenCache struct {
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
}
