package messege

import (
	"context"
)

type Repository interface {
	CreateMessege(ctx context.Context, mess *Messege, tableName string) error
	FindAllMesseges(ctx context.Context, tableName string) (messeges []Messege, err error)
	AddViewer(ctx context.Context, mess *Messege, tableName, usrID string) error
}
