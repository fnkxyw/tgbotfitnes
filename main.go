package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"tgbotfitnes/handler"
)

var keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Мои данные"),
		tgbotapi.NewKeyboardButton("Изменить данные"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6775953510:AAFvm_iNNVDm-38eE2iaQdlkXyWydPUfgbY")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	var users map[int64]handler.User = make(map[int64]handler.User)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() && update.Message.Command() == "start" {
			if _, exists := users[update.Message.From.ID]; !exists {
				var newUser handler.User
				newUser.ID = update.Message.From.ID

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваше имя:")

				bot.Send(msg)

				update = <-updates

				newUser.Name = update.Message.Text

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш вес:")
				bot.Send(msg)

				update = <-updates

				newUser.Weight, _ = strconv.Atoi(update.Message.Text)

				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш рост:")
				bot.Send(msg)

				update = <-updates
				if update.Message.Text == "Отмена" {
					continue
				}
				newUser.Height, _ = strconv.Atoi(update.Message.Text)

				users[update.Message.From.ID] = newUser
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите кнопку:")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}

		switch update.Message.Text {
		case "Мои данные":
			currentUser := users[update.Message.From.ID]
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, handler.CreateMessageAboutNameHeightWeigth(&currentUser))
			bot.Send(msg)
		case "Изменить данные":
			currentUser := users[update.Message.From.ID]

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваше имя:")
			bot.Send(msg)

			update = <-updates

			currentUser.Name = update.Message.Text

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш вес:")
			bot.Send(msg)

			update = <-updates

			currentUser.Weight, _ = strconv.Atoi(update.Message.Text)

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш рост:")
			bot.Send(msg)

			update = <-updates

			currentUser.Height, _ = strconv.Atoi(update.Message.Text)

			users[update.Message.From.ID] = currentUser

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваши данные успешно обновлены.")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}
	}
}
