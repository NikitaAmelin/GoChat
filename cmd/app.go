package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"goydamess/GoChat/internal/domain"
	"goydamess/GoChat/internal/user/db"
	"goydamess/GoChat/pkg/data_base"
	"net/http"
)

var clients = make(map[*websocket.Conn]bool)

var channel = make(chan string)

func writeUsersMassages() {
	for {
		txt_msg := <-channel
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(txt_msg))
			if err != nil {
				fmt.Print(fmt.Errorf("ошибка отправки сообщения %w", err))
				client.Close()
				delete(clients, client)
				continue
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
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

	postgreSQLclient, err := postgresql.NewClient(context.TODO(), 3, "postgres", "post1212", "localhost", "8080", "Chat_db")
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось запустить postgres: %w", err))
		return
	}
	rep := user.NewRepository(postgreSQLclient)
	all, err := rep.FindAll(context.TODO())
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось найти пользователей: %w", err))
		return
	}
	for _, usr := range all {
		fmt.Println(usr)
	}

	go writeUsersMassages()
	http.HandleFunc("/ws", handler)
	fmt.Println("Сервер открыт на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Print(fmt.Errorf("не удалось открыть сервер: %w", err))
		return
	}
}
