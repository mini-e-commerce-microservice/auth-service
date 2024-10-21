package model

type TokenCache struct {
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
	RegisterAs      int8   `json:"register_as"`
}