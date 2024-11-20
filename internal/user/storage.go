package user

import (
	"context"
)

type Repository interface {
	CreateUser(ctx context.Context, usr *User) error
	FindAllUsers(ctx context.Context) (users []User, err error)
	FindUserByID(ctx context.Context, id string) (User, error)
	//Update(ctx context.Context, user User) error
	//Delete(ctx context.Context, id string) error
}
