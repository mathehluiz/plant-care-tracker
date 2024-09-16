package domain

import "context"

type UserStorer interface {
	CreateUser(ctx context.Context, user *User) (string, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByExternalId(ctx context.Context, id string) (*User, error)
	AddRolesToUser(ctx context.Context, id string, roles []string) error
	RemoveRolesFromUser(ctx context.Context, id string, roles []string) error
	VerifyUser(ctx context.Context, id string) error
	UpdateActiveUserStatus(ctx context.Context, id string, active bool) error
	DeleteUser(ctx context.Context, id string) error
	UpdatePassword(ctx context.Context, id, password string) error
}
