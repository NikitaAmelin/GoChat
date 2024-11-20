package chat

import (
	"context"
)

type Repository interface {
	CreateChat(ctx context.Context, chat *Chat) error
	FindMembers(ctx context.Context) (members []string, err error)
	FindChatByID(ctx context.Context, id string) (Chat, error)
	AddMember(ctx context.Context, chat *Chat) error
	CreateMessegesDB(ctx context.Context, chat *Chat)
}
