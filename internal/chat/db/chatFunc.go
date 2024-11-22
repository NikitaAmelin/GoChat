package chatFunc

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"goydamess/GoChat/internal/TablesDB"
	"goydamess/GoChat/internal/chat"
	"goydamess/GoChat/pkg/data_base"
)

type repository struct {
	client postgresql.Client
}

func (r repository) CreateChat(ctx context.Context, chat *chat.Chat, pgclient *postgresql.Client) error {
	q := `INSERT INTO "Chats" ("Name", "Members") VALUES ($1, $2) RETURNING id`
	if err := r.client.QueryRow(ctx, q, chat.Name, chat.Members).Scan(&chat.ID); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			newerr := fmt.Errorf(fmt.Sprintf("SQL error: %s, detail: %s, where: %s, code:%s", pgerr.Message, pgerr.Detail, pgerr.Where, pgerr.Code))
			fmt.Println(newerr)
			return newerr
		}
		return err
	}

	nameMessegesDB := fmt.Sprintf("History_%s", chat.ID)
	q2 := `INSERT INTO "Chats" ("NameMessegesDB") VALUES ($1)`
	_, err := r.client.Exec(ctx, q2, nameMessegesDB)
	if err != nil {
		return err
	}
	chat.NameMessegesDB = nameMessegesDB
	rep := TablesDB.Repository{Client: *pgclient}
	err = rep.CreateMessegesTable(context.TODO(), nameMessegesDB)
	if err != nil {
		fmt.Println(fmt.Errorf("не удалось создать таблицу Messeges для чата %s: %w", chat.Name, err))
		return err
	}

	return nil
}

func (r repository) FindMembers(ctx context.Context) (members []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindChatByID(ctx context.Context, id string) (chat.Chat, error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) AddMember(ctx context.Context, chat *chat.Chat) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) CreateMessegesDB(ctx context.Context, chat *chat.Chat) {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client) chat.Repository {
	return &repository{
		client,
	}
}
