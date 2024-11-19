package main

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func writeMesseges(conn *websocket.Conn) {
	for {
		var msg string
		fmt.Scan(&msg)
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			fmt.Print(fmt.Errorf("ошибка отправки сообщения: %w", err))
			return
		}
	}
}

func printMesseges(conn *websocket.Conn) { // отрисовка сообщений для пользователей
	for {
		_, byte_msg, err := conn.ReadMessage()
		msg := string(byte_msg)
		if err != nil {
			fmt.Print(fmt.Errorf("ошибка чтения сообщения: %w", err))
		}
		fmt.Println(msg)
	}
}

func main() {
	addr := "ws://192.168.0.104:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось подключиться к серверу: %w", err))
	}
	defer conn.Close()
	go writeMesseges(conn)
	printMesseges(conn)
}
