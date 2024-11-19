package app

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func WriteUsersMassages(clients map[*websocket.Conn]bool, channel chan string) {
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
