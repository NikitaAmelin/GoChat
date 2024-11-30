package main

import (
	"context"
	"fmt"
	"goydamess/internal/ports/handlers"
	"goydamess/internal/storage"
	postgresql "goydamess/pkg/data_base"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	postgreSQLclient, err := postgresql.NewClient(context.TODO(), 3, "user", "qwerty", "localhost", "5433", "gochat")
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось запустить postgres: %w", err))
		return
	}
	//rep := TablesDB.Repository{postgreSQLclient}
	//err = rep.CreateUsersTable(context.TODO())
	//if err != nil {
	//	fmt.Print(fmt.Errorf("не удалось coздать таблицу users: %w", err))
	//	return
	//}
	var user storage.UserStorage
	var chat storage.ChatStorage
	var mess storage.MessageStorage
	s := storage.NewStorage(postgreSQLclient, user, chat, mess)
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}
	handler := handlers.NewHandler(s, upgrader)
	http.HandleFunc("/ws/user/auth/login/", handler.Login)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Print(fmt.Errorf("не удалось открыть сервер: %w", err))
		return
	}
}
