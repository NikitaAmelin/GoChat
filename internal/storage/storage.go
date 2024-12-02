package storage

import postgresql "goydamess/pkg/data_base"

type Storage interface {
	UserStorage
	ChatStorage
	MessageStorage
}

type storage struct {
	Client  postgresql.Client
	User    UserStorage
	Chat    ChatStorage
	Messege MessageStorage
}

func NewStorage(client postgresql.Client, user UserStorage, chat ChatStorage, messege MessageStorage) Storage {
	return &storage{
		Client:  client,
		User:    user,
		Chat:    chat,
		Messege: messege,
	}
}
