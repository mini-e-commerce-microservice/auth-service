package users

import "context"

type Repository interface {
	FindOneUser(ctx context.Context, input FindOneUserInput) (output FindOneUserOutput, err error)

	// CheckExistingUser if return true, data available
	CheckExistingUser(ctx context.Context, input CheckExistingUserInput) (exists bool, err error)

	UpSertUser(ctx context.Context, input UpSertUserInput) (output UpSertUserOutput, err error)
}
