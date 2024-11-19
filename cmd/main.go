package main

import (
	"My_local_chat/internal/app"
	"My_local_chat/internal/domain"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)

var channel = make(chan string)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка преобразования протокола: %w", err))
		return
	}
	defer conn.Close()
	clients[conn] = true
	fmt.Println("New user join the chat!")
	channel <- "New user join the chat!"
	for {
		var msg domain.Messege
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Print(fmt.Errorf("ошибка чтения сообщения: %w", err))
			delete(clients, conn)
			return
		}
		channel <- msg.Username + ": " + msg.Data
	}
}

func main() {
	go app.WriteUsersMassages(clients, channel)
	http.HandleFunc("/ws", handleConnections)
	fmt.Println("Сервер открыт на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Print(fmt.Errorf("не удалось открыть сервер: %w", err))
		return
	}
}
