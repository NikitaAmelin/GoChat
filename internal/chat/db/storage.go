package db

import (
	"context"
	"goydamess/GoChat/internal/chat"
	"goydamess/GoChat/pkg/data_base"
)

type repository struct {
	client postgresql.Client
}

func (r repository) CreateChat(ctx context.Context, chat *chat.Chat) error {
	//TODO implement me
	panic("implement me")
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
