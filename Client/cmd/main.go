package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

type Messege struct {
	Data     string `json:"data"`
	Username string `json:"username"`
}

type User struct {
	name string
}

func writeMesseges(conn *websocket.Conn, name string) {
	for {
		msg := Messege{Username: name}
		var err error
		r := bufio.NewReader(os.Stdin)
		msg.Data, err = r.ReadString('\n')
		if err != nil {
			fmt.Print(fmt.Errorf("ошибка считывания сообщения: %w", err))
			return
		}
		err = conn.WriteJSON(msg)
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
		fmt.Print(msg)
	}
}

func main() {
	addr := "ws://localhost:8080/ws"
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось подключиться к серверу: %w", err))
	}
	defer conn.Close()
	var name string
	fmt.Print("Введите ваше имя: ")
	fmt.Scan(&name)
	go writeMesseges(conn, name)
	printMesseges(conn)
}
