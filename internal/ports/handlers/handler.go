package handlers

import (
	"context"
	"fmt"
	"goydamess/internal/domain"
	"goydamess/internal/domain/response"
	"goydamess/internal/storage"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHandler struct {
	Storage     storage.Storage
	Upgrader    websocket.Upgrader
	OnlineUsers map[*websocket.Conn]bool
}

func NewHandler(s *storage.Storage, u *websocket.Upgrader) *WSHandler {
	return &WSHandler{
		Storage:  *s,
		Upgrader: *u,
	}
}

/*func (h *WSHandler) NewMessage(w http.ResponseWriter, r *http.Request) {
	var mess domain.Message
	conn.ReadJSON(&mess)
	if err = h.Storage.CreateMessage(context.TODO(), &mess, mess.NameMessagesDB); err != nil {
		fmt.Println(fmt.Errorf("ошибка записи сообщения в БД: %w", err))
	}
}*/

func (h *WSHandler) Login(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка преобразования протокола: %w", err))
	}
	fmt.Println("Пользователь подключён")
	h.OnlineUsers[conn] = true
	var user domain.User
	err = conn.ReadJSON(&user)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка входа на сайт: %w", err))
		conn.Close()
		delete(h.OnlineUsers, conn)
	}
	err = h.Storage.CreateUser(context.TODO(), &user)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка добавления пользователя: %w", err))
		conn.Close()
		delete(h.OnlineUsers, conn)
	}
	id := response.ID{
		ID: user.ID,
	}
	conn.WriteJSON(&id)
}
