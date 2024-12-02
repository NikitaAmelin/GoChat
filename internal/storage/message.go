package storage

import (
	"context"
	"fmt"
	"goydamess/internal/domain"

	"github.com/jackc/pgconn"
)

type MessageStorage interface {
	CreateMessage(ctx context.Context, mess *domain.Message, tableName string) error
	GetAllMessages(ctx context.Context, tableName string) (messages []domain.Message, err error)
	AddViewer(ctx context.Context, mess *domain.Message, tableName, usrID string) error
}

func (s *storage) CreateMessage(ctx context.Context, mess *domain.Message, tableName string) error {
	q := `INSERT INTO $1 ("Author", "Text", "TimeOfSending", "Viewed") VALUES ($2, $3, $4, $5) RETURNING id`
	if err := s.Client.QueryRow(ctx, q, tableName, mess.Author, mess.Text, mess.TimeOfSending, mess.Viewed).Scan(&mess.ID); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			newerr := fmt.Errorf(fmt.Sprintf("SQL error: %s, detail: %s, where: %s, code:%s", pgerr.Message, pgerr.Detail, pgerr.Where, pgerr.Code))
			fmt.Println(newerr)
			return newerr
		}
		return err
	}
	return nil
}

func (s *storage) GetAllMessages(ctx context.Context, tableName string) (messages []domain.Message, err error) {
	q := `SELECT ID, "Login", "Password" from $1`
	rows, err := s.Client.Query(ctx, q, tableName)
	if err != nil {
		return nil, err
	}
	messages = make([]domain.Message, 0)
	for rows.Next() {
		var mess domain.Message
		err = rows.Scan(&mess.ID, &mess.Author, &mess.Text, &mess.TimeOfSending, &mess.Viewed)
		if err != nil {
			return nil, err
		}
		messages = append(messages, mess)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func (s *storage) AddViewer(ctx context.Context, mess *domain.Message, tableName, usrID string) error {
	mess.Viewed = append(mess.Viewed, usrID)
	q := `INSERT INTO $1 ("Viewed") VALUES ($2)`
	_, err := s.Client.Exec(ctx, q, tableName, mess.Viewed)
	if err != nil {
		return err
	}
	return nil
}
