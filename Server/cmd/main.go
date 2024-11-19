package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	postgresql "goydamess/GoChat/pkg/data_base"
	"net/http"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
},
}

type clnt struct {
	status bool
	name   string `default:"void"`
}

var clients = make(map[*websocket.Conn]clnt)

var channel = make(chan string)
var sys_channel = make(chan string)

func writeUsersMassages() {
	for {
		msg := <-channel
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s: %s", clients[client].name, msg)))
			if err != nil {
				fmt.Print(fmt.Errorf("ошибка отправки сообщения %w", err))
				client.Close()
				delete(clients, client)
				continue
			}
		}
	}
}

func writeServerMassages() {
	for {
		msg := <-sys_channel
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", msg)))
			if err != nil {
				fmt.Print(fmt.Errorf("ошибка отправки серверного сообщения %w", err))
				continue
			}
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	err2 := conn.WriteMessage(websocket.TextMessage, []byte("VVEDITE NAME"))
	if err2 != nil {
		fmt.Print(fmt.Errorf("ошибка отправки сообщения %w", err))
		conn.Close()
		return
	}
	msg_name := ""
	for {
		_, byte_msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Print(fmt.Errorf("ошибка чтения сообщения: %w", err))
			delete(clients, conn)
			return
		}
		msg_name = string(byte_msg)
		break
	}
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка преобразования протокола: %w", err))
		return
	}
	defer conn.Close()

	fmt.Printf("New user %s join the chat!", msg_name)
	sys_channel <- fmt.Sprintf("New user, %s, join the chat!", msg_name)
	clients[conn] = clnt{true, msg_name}
	for {
		_, byte_msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Print(fmt.Errorf("ошибка чтения сообщения: %w", err))
			delete(clients, conn)
			return
		}
		msg := string(byte_msg)
		channel <- msg
	}
}

func main() {
	client, err := postgresql.NewClient()(context.TODO(), 3, "postgres", "post1212", "localhost", "8080", "Chat_db")
	go writeUsersMassages()
	go writeServerMassages()
	http.HandleFunc("/ws", handler)
	fmt.Println("Сервер открыт на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Print(fmt.Errorf("не удалось открыть сервер: %w", err))
		return
	}
}
