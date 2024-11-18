package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
},
}

var clients = make(map[*websocket.Conn]bool)

var channel = make(chan string)

type Messege struct {
	Data     string `json:"data"`
	Username string `json:"username"`
}

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
		var msg Messege
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
	go writeUsersMassages()
	http.HandleFunc("/ws", handler)
	fmt.Println("Сервер открыт на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Print(fmt.Errorf("не удалось открыть сервер: %w", err))
		return
	}
}


/*
БД
Точка входа (main)
Bind нужная вещь

internal   !
internal/domain (структуры)
intternal/app (сервисы с бизнес-логикой, как храним, что храним, где храним)
intternal/app/media-task.go 
GetMediaTask
разобраться с cookie
разобраться, как работать с постгрой


Во вторник встреча в 19:00

BaumanLegendsBeckend

Ubuntu, docker, git, go

Убираю Server, навожу порядок по аналогии с BaumanLegeds
папка Deployments (хранятся файлы докеркомпоуза (кусочек, позволяющий описывать, как будет развернуто приложение)) файл для конфигурации docker compous
dockercompous.yaml

internal/
domain (реквесты, респонзы)

ports (входные данные)   api.go (что происходит с запросом) глягуть легенды

app

storage

используем gorilla router

postgress на порту 5432 по умолчанию

goSQLmigrate
goose

папка миграции
/*
>>>>>>> 4d65f15d0c766a2ca2f351bdf4b2013d30fd29e5
