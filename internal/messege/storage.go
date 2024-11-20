package messege

import (
	"context"
)

type Repository interface {
	CreateMessege(ctx context.Context, mess *Messege) error
	FindAllMesseges(ctx context.Context) (messeges []Messege, err error)
	AddViewer(ctx context.Context, mess *Messege) error
}
