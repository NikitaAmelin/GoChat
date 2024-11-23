package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"goydamess/GoChat/internal/messege"
	"goydamess/GoChat/pkg/data_base"
)

type repository struct {
	client postgresql.Client
}

func (r repository) CreateMessege(ctx context.Context, mess *messege.Messege, tableName string) error {
	q := `INSERT INTO $1 ("Author", "Text", "TimeOfSending", "Viewed") VALUES ($2, $3, $4, $5) RETURNING id`
	if err := r.client.QueryRow(ctx, q, tableName, mess.Author, mess.Text, mess.TimeOfSending, mess.Viewed).Scan(&mess.ID); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			newerr := fmt.Errorf(fmt.Sprintf("SQL error: %s, detail: %s, where: %s, code:%s", pgerr.Message, pgerr.Detail, pgerr.Where, pgerr.Code))
			fmt.Println(newerr)
			return newerr
		}
		return err
	}
	return nil
}

func (r repository) FindAllMesseges(ctx context.Context, tableName string) (messeges []messege.Messege, err error) {
	q := `SELECT ID, "Login", "Password" from $1`
	rows, err := r.client.Query(ctx, q, tableName)
	if err != nil {
		return nil, err
	}
	messeges = make([]messege.Messege, 0)
	for rows.Next() {
		var mess messege.Messege
		err = rows.Scan(&mess.ID, &mess.Author, &mess.Text, &mess.TimeOfSending, &mess.Viewed)
		if err != nil {
			return nil, err
		}
		messeges = append(messeges, mess)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return messeges, nil
}

func (r repository) AddViewer(ctx context.Context, mess *messege.Messege, tableName, usrID string) error {
	mess.Viewed = append(mess.Viewed, usrID)
	q := `INSERT INTO $1 ("Viewed") VALUES ($2)`
	_, err := r.client.Exec(ctx, q, tableName, mess.Viewed)
	if err != nil {
		return err
	}
	return nil
}

func NewRepository(client postgresql.Client) messege.Repository {
	return &repository{
		client,
	}
}
