package users

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
)

func (r *repository) CheckExistingUser(ctx context.Context, input CheckExistingUserInput) (exists bool, err error) {
	query := r.sq.Select("1").
		Prefix("SELECT EXISTS(").
		Suffix(")").
		From("users")

	if input.Email.Valid {
		query = query.Where(squirrel.Eq{"email": input.Email.String})
	}
	if input.ID.Valid {
		query = query.Where(squirrel.Eq{"id": input.ID.Int64})
	}
	if input.IsEmailVerified.Valid {
		query = query.Where(squirrel.Eq{"is_email_verified": input.IsEmailVerified.Bool})
	}

	err = r.rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeDefault, &exists)
	if err != nil {
		return exists, collection.Err(err)
	}
	return
}
