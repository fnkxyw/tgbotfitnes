package main

import (
	"log"
	"strconv"
	"sync"
	"tgbotfitnes/db"
	"tgbotfitnes/handler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Мои данные"),
			tgbotapi.NewKeyboardButton("Изменить данные"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Создать план на тренировку"),
			tgbotapi.NewKeyboardButton("Мои тренировки"),
		),
	)

	keyboard2 = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Упражнения на грудные мышцы"),
			tgbotapi.NewKeyboardButton("Упражнения на спину"),
			tgbotapi.NewKeyboardButton("Упражнения на ноги"),
			tgbotapi.NewKeyboardButton("Упражнения на плечи"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Упражнения на бицепс и трицепс"),
			tgbotapi.NewKeyboardButton("Упражнения на пресс"),
			tgbotapi.NewKeyboardButton("Упражнения на кардио"),
			tgbotapi.NewKeyboardButton("Назад"),
		),
	)
	previouskeyboard tgbotapi.ReplyKeyboardMarkup
	newUserMu        sync.Mutex // Мьютекс для newUser
	currentUserMu    sync.Mutex // Мьютекс для currentUser
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6775953510:AAFvm_iNNVDm-38eE2iaQdlkXyWydPUfgbY")
	if err != nil {
		log.Panic(err)
	}
	db := database.DbConnectin()
	defer db.Close()

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	var newUser helper.User
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() && update.Message.Command() == "start" {

			newUserMu.Lock() // Блокировка доступа к newUser
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
				newUserMu.Unlock() // Разблокировка доступа к newUser
				continue
			}
			newUser.Height, _ = strconv.Atoi(update.Message.Text)

			err := database.InsertUser(db, &newUser)
			if err != nil {
				log.Fatal(err)
			}
			newUserMu.Unlock() // Разблокировка доступа к newUser

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите кнопку:")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}

		switch update.Message.Text {
		case "Мои данные":
			currentUserMu.Lock() // Блокировка доступа к currentUser
			currentUser, err := database.InfoAboutUser(db, update.Message.From.ID)
			currentUserMu.Unlock() // Разблокировка доступа к currentUser
			if err != nil {
				log.Println(err)
				break
			}
			database.PrintTable(db)
			calorie, quantity := database.SummaryInfo(db, currentUser.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, helper.CreateMessageAboutNameHeightWeigth((&currentUser), calorie, quantity))
			bot.Send(msg)
		case "Изменить данные":
			currentUserMu.Lock() // Блокировка доступа к currentUser
			currentUser, err := database.InfoAboutUser(db, update.Message.From.ID)
			if err != nil {
				currentUserMu.Unlock() // Разблокировка доступа к currentUser
				log.Println(err)
				break
			}
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

			err = database.UpdateUser(db, &currentUser)
			if err != nil {
				log.Fatal(err)
			}
			currentUserMu.Unlock() // Разблокировка доступа к currentUser

			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Ваши данные успешно обновлены.")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		case "Создать план на тренировку":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите группу тренировок: ")
			msg.ReplyMarkup = keyboard2
			bot.Send(msg)
			helper.CreateTrainHandler(bot, update, keyboard2, keyboard, updates, db)
		case "Мои тренировки":
			helper.MyTrain(db, bot, update, updates)
			fallthrough
		case "Назад":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите опцию: ")
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
		}
	}
}
