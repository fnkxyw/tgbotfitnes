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

func CreateMessageAboutNameHeightWeigth(u *User, calorie float64, quantity int) string {
	msg := fmt.Sprintf("Пользователь: %s \nВес %d кг \nРост %d см\nСуммарно потрачено колорий %0.2f\nСуммарное количество тренировок %d", u.Name, u.Weight, u.Height, calorie, quantity)
	return msg
}
