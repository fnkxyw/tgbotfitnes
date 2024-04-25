package handler

import (
	"strconv"
	"strings"
)

type User struct {
	Name   string
	Weight int
	Height int
}

func CreateUser() User {
	user := User{
		Name:   "",
		Weight: 0,
		Height: 0,
	}
	return user
}

func CreateMessageAboutNameHeightWeigth(u *User) string {
	res := []string{"Пользователь", u.Name, "весит", strconv.Itoa(u.Weight), "кг и его рост", strconv.Itoa(u.Height), "см"}
	msg := strings.Join(res, " ")
	return msg
}
