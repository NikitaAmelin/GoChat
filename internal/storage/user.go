package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"goydamess/internal/domain"
)

type UserStorage interface {
	CreateUser(ctx context.Context, usr *domain.User) error
	FindAllUsers(ctx context.Context) (users []domain.User, err error)
	FindUserByID(ctx context.Context, id string) (domain.User, error)
	CheckIfExist(ctx context.Context, Login string) (bool, error)
	FindUserByLogin(ctx context.Context, Login string) (domain.User, error)
}

func (s *storage) CreateUser(ctx context.Context, usr *domain.User) error {
	q := `INSERT INTO "Users" ("Login", "Password") VALUES ($1, $2) RETURNING id`
	if err := s.Client.QueryRow(ctx, q, usr.Login, usr.Password).Scan(&usr.ID); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			newerr := fmt.Errorf(fmt.Sprintf("SQL error: %s, detail: %s, where: %s, code:%s", pgerr.Message, pgerr.Detail, pgerr.Where, pgerr.Code))
			fmt.Println(newerr)
			return newerr
		}
		return err
	}
	return nil
}

func (s *storage) FindAllUsers(ctx context.Context) (users []domain.User, err error) {
	q := `SELECT ID, "Login", "Password" from "Users"`
	rows, err := s.Client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	users = make([]domain.User, 0)
	for rows.Next() {
		var usr domain.User
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

func (s *storage) FindUserByID(ctx context.Context, id string) (domain.User, error) {
	q := `SELECT id, "Login", "Password" from "Users" WHERE id = $1`
	var usr domain.User
	err := s.Client.QueryRow(ctx, q, id).Scan(&usr.ID, &usr.Login, &usr.Password)
	if err != nil {
		return domain.User{}, err
	}
	return usr, nil
}

func (s *storage) CheckIfExist(ctx context.Context, Login string) (bool, error) {
	q := `SELECT id FROM "Users" WHERE "Login"=$1`

	var flag string
	err := s.Client.QueryRow(ctx, q, Login).Scan(&flag)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if len(flag) > 0 {
		return true, nil
	}
	return false, nil
}

func (s *storage) FindUserByLogin(ctx context.Context, Login string) (domain.User, error) {
	q := `SELECT id, "Login", "Password" from "Users" WHERE "Login" = $1`
	var usr domain.User
	err := s.Client.QueryRow(ctx, q, Login).Scan(&usr.ID, &usr.Login, &usr.Password)
	if err != nil {
		return domain.User{}, err
	}
	return usr, nil
}
