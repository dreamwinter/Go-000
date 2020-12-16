package biz

import "context"

// UserDO is the domain object to store user info
type UserDO struct {
	ID     string
	Name   string
	Gender string
	Age    int
}

// UserRepository defines the interface functions
type UserRepository interface {
	CreateUser(ctx context.Context, user UserDO) (UserDO, error)
	GetUser(ctx context.Context, id string) (UserDO, error)
	UpdateUser(ctx context.Context, user UserDO) (UserDO, error)
	DeleteUser(ctx context.Context, id string) error
	ListUsers(ctx context.Context) ([]UserDO, error)
}
