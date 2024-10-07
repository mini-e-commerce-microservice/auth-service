package model

import (
	"time"
)

type User struct {
	ID              int64      `db:"id" json:"id"`
	Email           string     `db:"email" json:"email"`
	Password        string     `db:"password" json:"password"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	IsEmailVerified bool       `db:"is_email_verified" json:"is_email_verified"`
	RegisterAs      int8       `db:"register_as" json:"register_as"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
	TraceParent     *string    `db:"trace_parent" json:"trace_parent"`
}
