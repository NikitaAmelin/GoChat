package storage

import (
	postgresql "goydamess/GoChat/pkg/data_base"
)

type Storage interface {
	UserStorage
	ChatStorage
	MessegeStorage
}

type storage struct {
	Client  postgresql.Client
	User    UserStorage
	Chat    ChatStorage
	Messege MessegeStorage
}

func NewStorage(client postgresql.Client, user UserStorage, chat ChatStorage, messege MessegeStorage) Storage {
	return &storage{
		Client:  client,
		User:    user,
		Chat:    chat,
		Messege: messege,
	}
}
