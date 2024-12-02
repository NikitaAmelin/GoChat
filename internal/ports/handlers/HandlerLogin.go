package handlers

import (
	"context"
	"fmt"
	"goydamess/internal/domain"
	"goydamess/internal/domain/response"
	"net/http"
)

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
				answ = "200 " + user.Login
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
