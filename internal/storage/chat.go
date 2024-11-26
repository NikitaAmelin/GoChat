package storage

import (
	"context"
	"fmt"
	"goydamess/internal/TablesDB"
	"goydamess/internal/domain"
	postgresql "goydamess/pkg/data_base"

	"github.com/jackc/pgconn"
)

type ChatStorage interface {
	CreateChat(ctx context.Context, chat *domain.Chat, pgclient *postgresql.Client) error
	FindChatByID(ctx context.Context, id string) (domain.Chat, error)
	AddMember(ctx context.Context, chat *domain.Chat, usrID string) error
}

func (s *storage) CreateChat(ctx context.Context, chat *domain.Chat, pgclient *postgresql.Client) error {
	q := `INSERT INTO "Chats" ("Name", "Members") VALUES ($1, $2) RETURNING id`
	if err := s.Client.QueryRow(ctx, q, chat.Name, chat.Members).Scan(&chat.ID); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			newerr := fmt.Errorf(fmt.Sprintf("SQL error: %s, detail: %s, where: %s, code:%s", pgerr.Message, pgerr.Detail, pgerr.Where, pgerr.Code))
			fmt.Println(newerr)
			return newerr
		}
		return err
	}

	nameMessagesDB := fmt.Sprintf("History_%s", chat.ID)
	q2 := `INSERT INTO "Chats" ("NameMessegesDB") VALUES ($1)`
	_, err := s.Client.Exec(ctx, q2, nameMessagesDB)
	if err != nil {
		return err
	}
	chat.NameMessagesDB = nameMessagesDB
	rep := TablesDB.Repository{Client: *pgclient}
	err = rep.CreateMessegesTable(context.TODO(), nameMessagesDB)
	if err != nil {
		fmt.Println(fmt.Errorf("не удалось создать таблицу Messeges для чата %s: %w", chat.Name, err))
		return err
	}

	return nil
}

func (s *storage) FindChatByID(ctx context.Context, id string) (domain.Chat, error) {
	q := `SELECT id, "Name", "Members", "NameMessegesDB" from "Users" WHERE id = $1`
	var cht domain.Chat
	err := s.Client.QueryRow(ctx, q, id).Scan(&cht.ID, &cht.Name, &cht.Members, &cht.NameMessagesDB)
	if err != nil {
		return domain.Chat{}, err
	}
	return cht, nil
}

func (s *storage) AddMember(ctx context.Context, chat *domain.Chat, usrID string) error {
	chat.Members = append(chat.Members, usrID)
	q := `INSERT INTO "Chats" ("Members") VALUES ($1)`
	_, err := s.Client.Exec(ctx, q, chat.Members)
	if err != nil {
		return err
	}
	return nil
}
