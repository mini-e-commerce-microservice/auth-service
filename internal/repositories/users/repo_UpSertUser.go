package users

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
)

func (r *repository) UpSertUser(ctx context.Context, input UpSertUserInput) (output UpSertUserOutput, err error) {
	query := r.sq.Insert("users").Columns(
		"id", "email", "password", "is_email_verified", "created_at", "updated_at", "deleted_at", "trace_parent",
	).Values(
		input.Payload.ID,
		input.Payload.Email,
		input.Payload.Password,
		input.Payload.IsEmailVerified,
		input.Payload.CreatedAt,
		input.Payload.UpdatedAt,
		input.Payload.DeletedAt,
		input.Payload.TraceParent,
	).Suffix(
		`ON CONFLICT (id) DO UPDATE
			  	SET email = EXCLUDED.email,
				 password = EXCLUDED.password,
				 is_email_verified = EXCLUDED.is_email_verified,
				 updated_at = EXCLUDED.updated_at,
				 trace_parent = EXCLUDED.trace_parent,
				 created_at = EXCLUDED.created_at,
				 deleted_at = EXCLUDED.deleted_at`,
	)

	rdbms := r.rdbms
	if input.Tx != nil {
		rdbms = input.Tx
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return output, collection.Err(err)
	}

	return
}
