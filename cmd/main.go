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
	postgreSQLclient, err := postgresql.NewClient(context.TODO(), 3, "postgres", "123", "localhost", "8080", "Chat_db")
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось запустить postgres: %w", err))
		return
	}
	var user storage.UserStorage
	var chat storage.ChatStorage
	var mess storage.MessageStorage
	s := storage.NewStorage(postgreSQLclient, user, chat, mess)
	var upgrader websocket.Upgrader
	handler := handlers.NewHandler(&s, &upgrader)
	http.HandleFunc("/ws/send/message", handler.NewMessage)
	http.HandleFunc("/ws/user/auth/login", handler.Login)
}
