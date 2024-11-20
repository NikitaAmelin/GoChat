package db

import (
	"context"
	"goydamess/GoChat/internal/messege"
	"goydamess/GoChat/pkg/data_base"
)

type repository struct {
	client postgresql.Client
}

func (r repository) CreateMessege(ctx context.Context, mess *messege.Messege) error {
	//TODO implement me
	panic("implement me")
}

func (r repository) FindAllMesseges(ctx context.Context) (messeges []messege.Messege, err error) {
	//TODO implement me
	panic("implement me")
}

func (r repository) AddViewer(ctx context.Context, mess *messege.Messege) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postgresql.Client) messege.Repository {
	return &repository{
		client,
	}
}
