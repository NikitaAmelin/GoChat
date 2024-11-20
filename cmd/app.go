package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"goydamess/GoChat/internal/TablesDB"
	"goydamess/GoChat/internal/domain"
	postgresql "goydamess/GoChat/pkg/data_base"
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
		var msg domain.Message
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
	rep := TablesDB.Repository{Client: postgreSQLclient}
	err = rep.CreateUsersTable(context.TODO())
	if err != nil {
		fmt.Println(fmt.Errorf("не удалось создать таблицу Users в postgres: %w", err))
		return
	}
	/*rep := userFunc.NewRepository(postgreSQLclient) //TODO репозиторий пользователя в постгресе

	u1 := user.User{
		Login:    "Ruslan",
		Password: "1wq2",
	}
	err = rep.CreateUser(context.TODO(), &u1) //пример создания, id записывается в u1.ID
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось создать пользователя: %w", err))
		return
	}
	fmt.Println(u1.ID)

	all, err := rep.FindAllUsers(context.TODO()) //пример поиска всех пользователей в базе
	if err != nil {
		fmt.Print(fmt.Errorf("не удалось найти пользователей: %w", err))
		return
	}
	for _, usr := range all {
		fmt.Println(usr)
	}
	u2, err := rep.FindUserByID(context.TODO(), "de618c38-9cf7-498b-9d98-eb6433ef09e9") //пример поиска пользователей в базе по id
	fmt.Printf("id: %s, log: %s, pass: %s", u2.ID, u2.Login, u2.Password)*/

	go writeUsersMassages()
	http.HandleFunc("/ws", handler)
	fmt.Println("Сервер открыт на порту :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Print(fmt.Errorf("не удалось открыть сервер: %w", err))
		return
	}
}
