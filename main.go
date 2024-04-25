package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"tgbotfitnes/handler"
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
	user := handler.CreateUser()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваше имя:")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Отмена"),
			),
		)

		bot.Send(msg)

		update = <-updates

		if update.Message.Text == "Отмена" {
			continue
		}

		user.Name = update.Message.Text

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш вес:")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Отмена"),
			),
		)

		bot.Send(msg)

		update = <-updates

		if update.Message.Text == "Отмена" {
			continue
		}

		user.Weight, _ = strconv.Atoi(update.Message.Text)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите ваш рост:")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Отмена"),
			),
		)

		bot.Send(msg)

		update = <-updates

		if update.Message.Text == "Отмена" {
			continue
		}

		user.Height, _ = strconv.Atoi(update.Message.Text)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, handler.CreateMessageAboutNameHeightWeigth(&user))
		bot.Send(msg)
	}
}
