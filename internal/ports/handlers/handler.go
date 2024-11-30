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

func NewHandler(s storage.Storage, u websocket.Upgrader) *WSHandler {
	return &WSHandler{
		Storage:     s,
		Upgrader:    u,
		OnlineUsers: make(map[*websocket.Conn]bool),
	}
}

/*func (h *WSHandler) NewMessage(w http.ResponseWriter, r *http.Request) {
	var mess domain.Message
	conn.ReadJSON(&mess)
	if err = h.Storage.CreateMessage(context.TODO(), &mess, mess.NameMessagesDB); err != nil {
		fmt.Println(fmt.Errorf("ошибка записи сообщения в БД: %w", err))
	}
}*/

func (h *WSHandler) Register(w http.ResponseWriter, r *http.Request) {
	h.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка преобразования протокола: %w", err))
		return
	}
	fmt.Println("Пользователь подключён на регистрацию")
	//h.OnlineUsers[conn] = true
	var user domain.User
	for {
		err = conn.ReadJSON(&user)
		if err != nil {
			fmt.Println(fmt.Errorf("ошибка входа на сайт: %w", err))
			return
		}
		// если пользователь существует -> предупреждение, иначе -> пишем в базу
		flag, err := h.Storage.CheckIfExist(context.TODO(), user.Login)
		if err != nil {
			fmt.Println(fmt.Errorf("ошибка проверки наличия пользователя: %w", err))
			return
		}
		if flag {
			answ := fmt.Sprintf("Пользователь %s уже существует", user.Login)
			fmt.Println(answ)
			a := response.Answer{
				Answer: answ,
			}
			err = conn.WriteJSON(&a)
			if err != nil {
				fmt.Println(fmt.Errorf("ошибка связи с сайтом: %w", err))
				return
			}
		} else {

			err = h.Storage.CreateUser(context.TODO(), &user)
			if err != nil {
				fmt.Println(fmt.Errorf("ошибка добавления пользователя: %w", err))
				return
			}
			a := response.Answer{
				Answer: user.ID,
			}
			err = conn.WriteJSON(&a)
			if err != nil {
				fmt.Println(fmt.Errorf("ошибка связи с сайтом: %w", err))
				return
			}
		}
	}
}

func (h *WSHandler) Login(w http.ResponseWriter, r *http.Request) {
	h.Upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Errorf("ошибка преобразования протокола: %w", err))
		return
	}
	fmt.Println("Пользователь подключён на вход")
	h.OnlineUsers[conn] = true
	var user domain.User
	for {
		err = conn.ReadJSON(&user)
		if err != nil {
			fmt.Println(fmt.Errorf("ошибка входа на сайт: %w", err))
			return
		}
		// если пользователь существует -> вход в его аккаунт, иначе -> предупреждение
		flag, err := h.Storage.CheckIfExist(context.TODO(), user.Login)
		if err != nil {
			fmt.Println(fmt.Errorf("ошибка проверки наличия пользователя: %w", err))
			return
		}
		if !flag {
			answ := fmt.Sprintf("Пользователь %s не существует", user.Login)
			fmt.Println(answ)
			a := response.Answer{
				Answer: answ,
			}
			err = conn.WriteJSON(&a)
			if err != nil {
				fmt.Println(fmt.Errorf("ошибка связи с сайтом: %w", err))
				return
			}
		} else {
			var userFromDb domain.User
			userFromDb, err = h.Storage.FindUserByLogin(context.TODO(), user.Login)
			if err != nil {
				fmt.Println(fmt.Errorf("ошибка добавления пользователя: %w", err))
				return
			}
			var answ string
			if user.Password != userFromDb.Password {
				answ = fmt.Sprintf("Неверный пароль для %s", user.Login)
			} else {
				answ = fmt.Sprintf("Вы успешно вошли в аккаунт %s", user.Login)
			}
			fmt.Println(answ)
			answer := response.Answer{
				Answer: answ,
			}
			err = conn.WriteJSON(&answer)
			if err != nil {
				fmt.Println(fmt.Errorf("ошибка связи с сайтом: %w", err))
				return
			}
		}
	}
}
