package helper

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

var (
	keyboardbreast = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Жим штанги лежа"),
			tgbotapi.NewKeyboardButton("Жим гантелей лежа"),
			tgbotapi.NewKeyboardButton("К другим упражнениям"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сведение рук в тренажере Баттерфляй."),
			tgbotapi.NewKeyboardButton("Пуловер"),
		),
	)
	keyboardback = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Тяга штанги к поясу"),
			tgbotapi.NewKeyboardButton("Подтягивания"),
			tgbotapi.NewKeyboardButton("К другим упражнениям"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Тяга гантелей в наклоне"),
			tgbotapi.NewKeyboardButton("Тяга к груди в каблучке"),
		),
	)
	keyboardlegs = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Приседания со штангой"),
			tgbotapi.NewKeyboardButton("Выпады"),
			tgbotapi.NewKeyboardButton("К другим упражнениям"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Жим ногами в тренажере"),
			tgbotapi.NewKeyboardButton("Разгибание ног в тренажере"),
		),
	)
	keyboardovens = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Армейский жим"),
			tgbotapi.NewKeyboardButton("Махи гантелей в стороны"),
			tgbotapi.NewKeyboardButton("К другим упражнениям"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Подъем гантелей на бицепс"),
			tgbotapi.NewKeyboardButton("Вертикальная тяга в тренажерее"),
		),
	)
	keyboardbiceps = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сгибание рук со штангой"),
			tgbotapi.NewKeyboardButton("Сгибание рук с гантелями"),
			tgbotapi.NewKeyboardButton("К другим упражнениям"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Французский жим"),
			tgbotapi.NewKeyboardButton("Трицепсовый жим"),
		),
	)
	keyboardpress = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Скручивания"),
			tgbotapi.NewKeyboardButton("Подъем ног в висе"),
			tgbotapi.NewKeyboardButton("К другим упражнениям"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Планка"),
			tgbotapi.NewKeyboardButton("Боковые наклоны"),
		),
	)
	keyboardcardio = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Бег на беговой дорожке"),
			tgbotapi.NewKeyboardButton("Велотренажер"),
			tgbotapi.NewKeyboardButton("К другим упражнениям"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Эллиптический тренажер"),
			tgbotapi.NewKeyboardButton("Степпер"),
		),
	)
	keyboardaccept = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Да"),
			tgbotapi.NewKeyboardButton("Нет"),
		),
	)
	counter int = 1
)

func CreateTrainHandler(bot *tgbotapi.BotAPI, update tgbotapi.Update, keyboard1, keyboard2 tgbotapi.ReplyKeyboardMarkup, updates <-chan tgbotapi.Update, db *sql.DB) {
	AddId(db, counter)
	for {

		update = <-updates
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите упражнение и введите количество")
		switch update.Message.Text {
		case "Упражнения на грудные мышцы":
			msg.ReplyMarkup = keyboardbreast
			bot.Send(msg)
		case "Упражнения на спину":
			msg.ReplyMarkup = keyboardback
			bot.Send(msg)
		case "Упражнения на ноги":
			msg.ReplyMarkup = keyboardlegs
			bot.Send(msg)
		case "Упражнения на плечи":
			msg.ReplyMarkup = keyboardovens
			bot.Send(msg)
		case "Упражнения на бицепс и трицепс":
			msg.ReplyMarkup = keyboardbiceps
			bot.Send(msg)
		case "Упражнения на пресс":
			msg.ReplyMarkup = keyboardpress
			bot.Send(msg)
		case "Упражнения на кардио":
			msg.ReplyMarkup = keyboardcardio
			bot.Send(msg)
		case "К другим упражнениям":
			msg.ReplyMarkup = keyboard1
			bot.Send(msg)
		case "Назад":
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Желаете сохранить ваш тренировочный план?")
			msg.ReplyMarkup = keyboardaccept
			bot.Send(msg)
			update = <-updates
			switch update.Message.Text {
			case "Да":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название для вашей тренировки")
				bot.Send(msg)
				update = <-updates
				err := AddName(db, update.Message.Text, counter)
				if err != nil {
					log.Println(err)
				}
				err = AddIdInUsers(db, counter, update)
				if err != nil {
					log.Println(err)
				}
				counter++
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите кнопку:")
				msg.ReplyMarkup = keyboard2
				bot.Send(msg)
				return
			case "Нет":
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите кнопку:")
				msg.ReplyMarkup = keyboard2
				bot.Send(msg)
				return
			}
		case "Жим штанги лежа":
			handleExercise(bot, update, db, "Жим штанги лежа", counter, updates)
		case "Жим гантелей лежа":
			handleExercise(bot, update, db, "Жим гантелей лежа", counter, updates)
		case "Сведение рук в тренажере Баттерфляй.":
			handleExercise(bot, update, db, "Сведение рук в тренажере Баттерфляй.", counter, updates)
		case "Пуловер":
			handleExercise(bot, update, db, "Пуловер", counter, updates)
		case "Тяга штанги к поясу":
			handleExercise(bot, update, db, "Тяга штанги к поясу", counter, updates)
		case "Подтягивания":
			handleExercise(bot, update, db, "Подтягивания", counter, updates)
		case "Тяга гантелей в наклоне":
			handleExercise(bot, update, db, "Тяга гантелей в наклоне", counter, updates)
		case "Тяга к груди в каблучке":
			handleExercise(bot, update, db, "Тяга к груди в каблучке", counter, updates)
		case "Приседания со штангой":
			handleExercise(bot, update, db, "Приседания со штангой", counter, updates)
		case "Выпады":
			handleExercise(bot, update, db, "Выпады", counter, updates)
		case "Жим ногами в тренажере":
			handleExercise(bot, update, db, "Жим ногами в тренажере", counter, updates)
		case "Разгибание ног в тренажере":
			handleExercise(bot, update, db, "Разгибание ног в тренажере", counter, updates)
		case "Армейский жим":
			handleExercise(bot, update, db, "Армейский жим", counter, updates)
		case "Махи гантелей в стороны":
			handleExercise(bot, update, db, "Махи гантелей в стороны", counter, updates)
		case "Подъем гантелей на бицепс":
			handleExercise(bot, update, db, "Подъем гантелей на бицепс", counter, updates)
		case "Вертикальная тяга в тренажерее":
			handleExercise(bot, update, db, "Вертикальная тяга в тренажерее", counter, updates)
		case "Сгибание рук со штангой":
			handleExercise(bot, update, db, "Сгибание рук со штангой", counter, updates)
		case "Сгибание рук с гантелями":
			handleExercise(bot, update, db, "Сгибание рук с гантелями", counter, updates)
		case "Французский жим":
			handleExercise(bot, update, db, "Французский жим", counter, updates)
		case "Трицепсовый жим":
			handleExercise(bot, update, db, "Трицепсовый жим", counter, updates)
		case "Скручивания":
			handleExercise(bot, update, db, "Скручивания", counter, updates)
		case "Подъем ног в висе":
			handleExercise(bot, update, db, "Подъем ног в висе", counter, updates)
		case "Планка":
			handleTimeExercise(bot, update, db, "Планка", counter, updates)
		case "Боковые наклоны":
			handleExercise(bot, update, db, "Боковые наклоны", counter, updates)
		case "Бег на беговой дорожке":
			handleTimeExercise(bot, update, db, "Бег на беговой дорожке", counter, updates)
		case "Велотренажер":
			handleTimeExercise(bot, update, db, "Велотренажер", counter, updates)
		case "Эллиптический тренажер":
			handleTimeExercise(bot, update, db, "Эллиптический тренажер", counter, updates)
		case "Степпер":
			handleTimeExercise(bot, update, db, "Степпер", counter, updates)
		}
	}
}

func AddTrainers(db *sql.DB, name string, counts int, id int) error {
	_, err := db.Exec("UPDATE public.trainers SET trainer = array_append(trainer, $1), count = array_append(count, $2) WHERE id = $3", name, counts, id)
	if err != nil {
		log.Println("Error adding trainers to database:", err)
		return err
	}
	log.Println("Trainers added successfully")
	return nil
}

func AddName(db *sql.DB, name string, id int) error {
	_, err := db.Exec("UPDATE public.trainers SET name = $1 WHERE id =$2", name, id)
	if err != nil {
		log.Println("Error adding name to database:", err)
		return err
	}
	log.Println("Name added successfully")
	return nil
}

func AddIdInUsers(db *sql.DB, id int, update tgbotapi.Update) error {
	_, err := db.Exec("UPDATE public.users SET trainersid = array_append(trainersid, $1) WHERE id=$2", id, update.Message.Chat.ID)
	if err != nil {
		log.Println("Error adding id in users table:", err)
		return err
	}
	log.Println("Id in users added successfully")
	return nil
}

func AddId(db *sql.DB, id int) error {
	_, err := db.Exec("INSERT INTO trainers (id) VALUES ($1) ON CONFLICT (id) DO NOTHING", id)
	return err
}

func handleExercise(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB, exerciseName string, counter int, updates tgbotapi.UpdatesChannel) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите количество повторений: ")
	bot.Send(msg)
	update = <-updates
	count := update.Message.Text
	value, err := strconv.Atoi(count)
	if err != nil {
		log.Println("Atoi error")
		return
	}
	err = AddTrainers(db, exerciseName, value, counter)
	if err == nil {
		SendTrainersToUser(bot, update, db, counter)
	} else {
		log.Println("Add error")
	}
}

func handleTimeExercise(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB, exerciseName string, counter int, updates tgbotapi.UpdatesChannel) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите время (в минутах): ")
	bot.Send(msg)
	update = <-updates
	timeInput := update.Message.Text
	duration, err := strconv.Atoi(timeInput)
	if err != nil {
		log.Println("Atoi error")
		return
	}
	err = AddTrainers(db, exerciseName, duration, counter)
	if err == nil {
		SendTrainersToUser(bot, update, db, counter)
	} else {
		log.Println("Add error")
	}
}

func DeleteFromTrainers(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM public.trainers WHERE id = $1", id)
	if err != nil {
		log.Println("Error adding id in users table:", err)
		return err
	}
	log.Println("Id in users added successfully")
	return nil
}

func SendTrainersToUser(bot *tgbotapi.BotAPI, update tgbotapi.Update, db *sql.DB, userID int) error {
	rows, err := db.Query("SELECT trainer, count FROM public.trainers WHERE id = $1", userID)
	if err != nil {
		log.Println("Error querying trainers:", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var trainers string
		var counts string
		if err := rows.Scan(&trainers, &counts); err != nil {
			log.Println("Error scanning trainers row:", err)
			continue
		}

		trainersArr := strings.Split(trainers[1:len(trainers)-1], ",") // Убираем квадратные скобки
		countsArr := strings.Split(counts[1:len(counts)-1], ",")

		if len(trainersArr) != len(countsArr) {
			log.Println("Mismatch between trainers and counts arrays")
			continue
		}

		var output strings.Builder
		for i := range trainersArr {
			output.WriteString(fmt.Sprintf("%s: %s\n", trainersArr[i], countsArr[i]))
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, output.String())
		if _, err := bot.Send(msg); err != nil {
			log.Println("Error sending message:", err)
			continue
		}
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over trainers rows:", err)
		return err
	}

	return nil
}
