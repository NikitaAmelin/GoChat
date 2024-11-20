package userFunc

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"goydamess/GoChat/internal/user"
	"goydamess/GoChat/pkg/data_base"
)

type repository struct {
	client postgresql.Client
}

func (r repository) CreateUser(ctx context.Context, usr *user.User) error {
	q := `INSERT INTO "Users" ("Login", "Password") VALUES ($1, $2) RETURNING id`
	if err := r.client.QueryRow(ctx, q, usr.Login, usr.Password).Scan(&usr.ID); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			newerr := fmt.Errorf(fmt.Sprintf("SQL error: %s, detail: %s, where: %s, code:%s", pgerr.Message, pgerr.Detail, pgerr.Where, pgerr.Code))
			fmt.Println(newerr)
			return newerr
		}
		return err
	}
	return nil
}

func (r repository) FindAllUsers(ctx context.Context) (users []user.User, err error) {
	q := `SELECT ID, "Login", "Password" from "Users"`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	users = make([]user.User, 0)
	for rows.Next() {
		var usr user.User
		err = rows.Scan(&usr.ID, &usr.Login, &usr.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r repository) FindUserByID(ctx context.Context, id string) (user.User, error) {
	q := `SELECT id, "Login", "Password" from "Users" WHERE id = $1`
	var usr user.User
	err := r.client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Login, &usr.Password)
	if err != nil {
		return user.User{}, err
	}
	return usr, nil
}

//func (r repository) Update(ctx context.Context, user user.User) error { }

//func (r repository) Delete(ctx context.Context, id string) error {}

func NewRepository(client postgresql.Client) user.Repository {
	return &repository{
		client,
	}
}
