package users_test

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/mini-e-commerce-microservice/auth-service/internal/model"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func Test_repository_FindOneUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := wsqlx.NewRdbms(sqlxDB)

	r := users.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := users.FindOneUserInput{
			ID:    null.IntFrom(rand.Int63()),
			Email: null.StringFrom(faker.Email()),
		}
		expectedOutput := users.FindOneUserOutput{
			Data: model.User{
				ID:              expectedInput.ID.Int64,
				Email:           expectedInput.Email.String,
				IsEmailVerified: true,
			},
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, is_email_verified, email FROM users WHERE id = $1 AND email = $2`,
		)).WithArgs(expectedInput.ID.Int64, expectedInput.Email.String).WillReturnRows(sqlmock.NewRows([]string{
			"id", "is_email_verified", "email",
		}).AddRow(
			expectedOutput.Data.ID,
			expectedOutput.Data.IsEmailVerified,
			expectedOutput.Data.Email,
		))

		output, err := r.FindOneUser(ctx, expectedInput)
		require.NoError(t, err)
		require.Equal(t, expectedOutput, output)
	})

	t.Run("should be return error no record", func(t *testing.T) {
		expectedInput := users.FindOneUserInput{
			ID:    null.IntFrom(rand.Int63()),
			Email: null.StringFrom(faker.Email()),
		}
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT id, is_email_verified, email FROM users WHERE id = $1 AND email = $2`,
		)).WithArgs(expectedInput.ID.Int64, expectedInput.Email.String).
			WillReturnError(sql.ErrNoRows)

		_, err = r.FindOneUser(ctx, expectedInput)
		require.ErrorIs(t, err, repositories.ErrRecordNotFound)
	})
}
