package helper

import (
	"fmt"
)

type User struct {
	ID     int64
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
	msg := fmt.Sprintf("Пользователь: %s \nВес %d кг \nРост %d см", u.Name, u.Weight, u.Height)
	return msg
}
