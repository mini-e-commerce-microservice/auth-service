package users_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func TestRepository_CreateUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := wsqlx.NewRdbms(sqlxDB)

	r := users.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedID := rand.Int63()
		expectedInput := users.CreateUserInput{
			Payload: model.User{
				Email:           faker.Email(),
				Password:        faker.Password(),
				IsEmailVerified: true,
				//RegisterAs:      int8(primitive.EnumRegisterAsMerchant),
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`INSERT INTO users (email,password,is_email_verified,register_as,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		)).WithArgs(
			expectedInput.Payload.Email,
			expectedInput.Payload.Password,
			expectedInput.Payload.IsEmailVerified,
			expectedInput.Payload.RegisterAs,
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).WillReturnRows(sqlmock.NewRows([]string{
			"id",
		}).AddRow(expectedID))

		output, err := r.CreateUser(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedID, output.ID)
	})
}
