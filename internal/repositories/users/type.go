package users

import (
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
)

type FindOneUserInput struct {
	ID    null.Int
	Email null.String
}

type FindOneUserOutput struct {
	Data model.User
}

type UpSertUserInput struct {
	Tx      wsqlx.Rdbms
	Payload model.User
}

type UpSertUserOutput struct {
}

type CheckExistingUserInput struct {
	ID              null.Int
	Email           null.String
	IsEmailVerified null.Bool
}

type UpdateUserInput struct {
	Tx      wsqlx.Rdbms
	ID      null.Int
	Email   null.String
	Payload UpdateUserInputPayload
}

type UpdateUserInputPayload struct {
	IsEmailVerified null.Bool
	Password        null.String
}
