package chat

import (
	"context"
	postgresql "goydamess/GoChat/pkg/data_base"
)

type Repository interface {
	CreateChat(ctx context.Context, chat *Chat, pgclient *postgresql.Client) error
	FindMembers(ctx context.Context) (members []string, err error)
	FindChatByID(ctx context.Context, id string) (Chat, error)
	AddMember(ctx context.Context, chat *Chat) error
	CreateMessegesDB(ctx context.Context, chat *Chat)
}
