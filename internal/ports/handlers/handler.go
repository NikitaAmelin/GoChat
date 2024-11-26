package handlers

import (
	"context"
	"fmt"
	"goydamess/internal/domain"
	"goydamess/internal/storage"
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHandler struct {
	Storage     storage.Storage
	Upgrader    websocket.Upgrader
	OnlineUsers map[*websocket.Conn]bool
}

func NewHandler(s *storage.Storage, upgrader *websocket.Upgrader) *WSHandler {
	return &WSHandler{
		Storage:  *s,
		Upgrader: *upgrader,
	}
}

func (h *WSHandler) NewMessage(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка преобразования протокола: %w", err))
	}
	var mess domain.Message
	conn.ReadJSON(&mess)
	if err = h.Storage.CreateMessage(context.TODO(), &mess, mess.NameMessagesDB); err != nil {
		fmt.Println(fmt.Errorf("ошибка записи сообщения в БД: %w", err))
	}
}

func (h *WSHandler) Login(w http.ResponseWriter, r *http.Request) {

	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка преобразования протокола: %w", err))
	}
	defer conn.Close()
	h.OnlineUsers[conn] = true
	var mess domain.Message
	conn.ReadJSON(&mess)
	if err = h.Storage.CreateMessage(context.TODO(), &mess, mess.NameMessagesDB); err != nil {
		fmt.Println(fmt.Errorf("ошибка записи сообщения в БД: %w", err))
	}
}
