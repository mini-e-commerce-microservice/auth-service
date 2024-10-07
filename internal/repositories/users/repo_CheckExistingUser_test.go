package users_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/go-faker/faker/v4"
	"github.com/guregu/null/v5"
	"github.com/jmoiron/sqlx"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories/users"
	"github.com/stretchr/testify/require"
	"math/rand"
	"regexp"
	"testing"
)

func TestRepository_CheckExistingUser(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	ctx := context.TODO()
	sqlxDB := sqlx.NewDb(dbMock, "sqlmock")

	sqlxx := wsqlx.NewRdbms(sqlxDB)

	r := users.NewRepository(sqlxx)

	t.Run("should be return correct", func(t *testing.T) {
		expectedInput := users.CheckExistingUserInput{
			ID:              null.IntFrom(rand.Int63()),
			Email:           null.StringFrom(faker.Email()),
			IsEmailVerified: null.BoolFrom(true),
		}

		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT EXISTS( SELECT 1 FROM users WHERE email = $1 AND id = $2 AND is_email_verified = $3 )`,
		)).WithArgs(expectedInput.Email.String, expectedInput.ID.Int64, expectedInput.IsEmailVerified.Bool).WillReturnRows(
			sqlmock.NewRows([]string{"exists"}).AddRow(true),
		)

		exists, err := r.CheckExistingUser(ctx, expectedInput)
		require.NoError(t, err)
		require.True(t, exists)
	})
}
