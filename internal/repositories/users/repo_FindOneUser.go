package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/mini-e-commerce-microservice/auth-service/internal/repositories"
)

func (r *repository) FindOneUser(ctx context.Context, input FindOneUserInput) (output FindOneUserOutput, err error) {
	query := r.sq.Select("id", "is_email_verified", "email", "password", "created_at").From("users")
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}
	if input.Email.Valid {
		query = query.Where(squirrel.Eq{"email": input.Email.String})
	}

	err = r.rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeStruct, &output.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = repositories.ErrRecordNotFound
		}

		return output, collection.Err(err)
	}

	return
}
