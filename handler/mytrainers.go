package helper

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	keyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Назад"),
			tgbotapi.NewKeyboardButton("Начать тренировку"),
		),
	)
)

type Plan struct {
	name     string
	count    []string
	trainers []string
}

func MyTrain(db *sql.DB, bot *tgbotapi.BotAPI, update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	ShowTrainers(db, bot, update)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите кнопку:")
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
	update = <-updates
	switch update.Message.Text {

	case "Начать тренировку":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите название тренировки, которую хотите начать:")
		bot.Send(msg)
		update = <-updates
		plan, err := GetTrainFromDb(db, update.Message.Text)
		if err != nil {
			log.Println("GetTrainFromDb error")
		}
		if len(plan.trainers) == 0 {
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Невереное название")
			bot.Send(msg)
			return
		}
		weight, err := GetWeightFromDb(db, update)
		if err != nil {
			log.Println("GetTrainFromDb error")
		}
		calorie := ParsePlanAndCalculate(&plan, weight)
		calorieAsString := fmt.Sprintf("%.2f", calorie) // Преобразование float64 в строку с двумя знаками после запятой
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "За эту тренировку вы потратите: "+calorieAsString)
		bot.Send(msg)
		err = UpdateSummaryTable(db, calorie, update)
		if err != nil {
			log.Println(err)
		}

	case "Назад":
		return
	}

}

func GetProcessedSlice(text []string) []string {
	word := strings.Join(text, " ")
	word = strings.ReplaceAll(word, `"`, "")
	text = strings.Split(word, `,`)
	return text
}

func ParsePlanAndCalculate(plan *Plan, weight int) float64 {
	var CalorieSum float64
	var TempSum float64
	trainersslice := GetProcessedSlice(plan.trainers)
	countslice := GetProcessedSlice(plan.count)
	for i := 0; i < len(trainersslice); i++ {
		repetions, err := strconv.Atoi(countslice[i])
		if err != nil {
			log.Println("Strconv error in mytrainers")
		}
		switch trainersslice[i] {
		case "Жим штанги лежа":
			TempSum = CalculateCaloriesBurnedBenchPress(weight, repetions)
		case "Жим гантелей лежа":
			TempSum = CalculateCaloriesBurnedDumbbellPress(weight, repetions)
		case "Сведение рук в тренажере Баттерфляй":
			TempSum = CalculateCaloriesBurnedButterfly(weight, repetions)
		case "Пуловер":
			TempSum = CalculateCaloriesBurnedPullover(weight, repetions)
		case "Тяга штанги к поясу":
			TempSum = CalculateCaloriesBurnedBarbellRow(weight, repetions)
		case "Подтягивания":
			TempSum = CalculateCaloriesBurnedPullUp(weight, repetions)
		case "Тяга гантелей в наклоне":
			TempSum = CalculateCaloriesBurnedDumbbellRow(weight, repetions)
		case "Тяга к груди в каблучке":
			TempSum = CalculateCaloriesBurnedCableChestPull(weight, repetions)
		case "Приседания со штангой":
			TempSum = CalculateCaloriesBurnedSquatsWithBarbell(weight, repetions)
		case "Выпады":
			TempSum = CalculateCaloriesBurnedLunges(weight, repetions)
		case "Жим ногами в тренажере":
			TempSum = CalculateCaloriesBurnedLegPress(weight, repetions)
		case "Разгибание ног в тренажере":
			TempSum = CalculateCaloriesBurnedLegExtension(weight, repetions)
		case "Армейский жим":
			TempSum = CalculateCaloriesBurnedMilitaryPress(weight, repetions)
		case "Махи гантелей в стороны":
			TempSum = CalculateCaloriesBurnedLateralRaises(weight, repetions)
		case "Подъем гантелей на бицепс":
			TempSum = CalculateCaloriesBurnedBicepCurl(weight, repetions)
		case "Вертикальная тяга в тренажерее":
			TempSum = CalculateCaloriesBurnedLatPullDown(weight, repetions)
		case "Сгибание рук со штангой":
			TempSum = CalculateCaloriesBurnedBarbellBicepCurl(weight, repetions)
		case "Сгибание рук с гантелями":
			TempSum = CalculateCaloriesBurnedDumbbellBicepCurl(weight, repetions)
		case "Французский жим":
			TempSum = CalculateCaloriesBurnedFrenchPress(weight, repetions)
		case "Трицепсовый жим":
			TempSum = CalculateCaloriesBurnedTricepPress(weight, repetions)
		case "Скручивания":
			TempSum = CalculateCaloriesBurnedCrunches(weight, repetions)
		case "Подъем ног в висе":
			TempSum = CalculateCaloriesBurnedLegRaises(weight, repetions)
		case "Планка":
			TempSum = CalculateCaloriesBurnedPlank(repetions)
		case "Боковые наклоны":
			TempSum = CalculateCaloriesBurnedSidePlank(repetions)
		case "Бег на беговой дорожке":
			TempSum = CalculateCaloriesBurnedTreadmillRunning(repetions)
		case "Велотренажер":
			TempSum = CalculateCaloriesBurnedStationaryBike(repetions)
		case "Эллиптический тренажер":
			TempSum = CalculateCaloriesBurnedElliptical(repetions)
		case "Степпер":
			TempSum = CalculateCaloriesBurnedStairStepper(repetions)
		}
		CalorieSum += TempSum
	}

	return CalorieSum
}

func UpdateSummaryTable(db *sql.DB, sum float64, update tgbotapi.Update) error {
	// Обновляем столбец `calorie`
	_, err := db.Exec("UPDATE public.summary SET calorie = calorie + $1 WHERE id = $2", sum, update.Message.Chat.ID)
	if err != nil {
		log.Println("Update summary error:", err)
		return err
	}

	// Увеличиваем значение столбца `quantity` на 1
	_, err = db.Exec("UPDATE public.summary SET quantity = quantity + 1 WHERE id = $1", update.Message.Chat.ID)
	if err != nil {
		log.Println("Update summary2 error:", err)
		return err
	}
	log.Println("UpdateSummaryTable corrceted")
	return nil
}

// жим штанги лежа
func CalculateCaloriesBurnedBenchPress(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.05
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// жим гантелей лежа
func CalculateCaloriesBurnedDumbbellPress(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// батерфляй
func CalculateCaloriesBurnedButterfly(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.03
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// пуловек
func CalculateCaloriesBurnedPullover(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.03
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// тяги штанги к поясу
func CalculateCaloriesBurnedBarbellRow(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// подтягивания
func CalculateCaloriesBurnedPullUp(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.05
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// тяга гантелей в наклоне
func CalculateCaloriesBurnedDumbbellRow(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// к груди в каблучке
func CalculateCaloriesBurnedCableChestPull(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// приседания со штангой
func CalculateCaloriesBurnedSquatsWithBarbell(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.08
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// выпады
func CalculateCaloriesBurnedLunges(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.07
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// жим ногами в тренажере
func CalculateCaloriesBurnedLegPress(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.06
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// разгибание ног в тренажере
func CalculateCaloriesBurnedLegExtension(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.05
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// армейский жим
func CalculateCaloriesBurnedMilitaryPress(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.05
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// махи гантелей в стороны
func CalculateCaloriesBurnedLateralRaises(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// подьем гантелей на бицепс
func CalculateCaloriesBurnedBicepCurl(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// вертикальная тяга в тренажере
func CalculateCaloriesBurnedLatPullDown(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.05
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// сгибание рук со штангой
func CalculateCaloriesBurnedBarbellBicepCurl(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// сгибанеие рук с гантелями
func CalculateCaloriesBurnedDumbbellBicepCurl(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// французский жим
func CalculateCaloriesBurnedFrenchPress(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.05
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// трицепсовый жим
func CalculateCaloriesBurnedTricepPress(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.05
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// скручивания
func CalculateCaloriesBurnedCrunches(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.03
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// подьем ног в висе
func CalculateCaloriesBurnedLegRaises(weight, repetitions int) float64 {
	caloriesBurnedPerRepetition := 0.04
	totalCaloriesBurned := float64(repetitions) * caloriesBurnedPerRepetition * float64(weight)
	return totalCaloriesBurned
}

// планка
func CalculateCaloriesBurnedPlank(durationMinutes int) float64 {
	caloriesBurnedPerMinute := 0.05
	totalCaloriesBurned := caloriesBurnedPerMinute * float64(durationMinutes)
	return totalCaloriesBurned
}

// боковые наклоны
func CalculateCaloriesBurnedSidePlank(durationMinutes int) float64 {
	caloriesBurnedPerMinute := 0.06
	totalCaloriesBurned := caloriesBurnedPerMinute * float64(durationMinutes)
	return totalCaloriesBurned
}

// бег на беговой дорожке
func CalculateCaloriesBurnedTreadmillRunning(durationMinutes int) float64 {
	caloriesBurnedPerMinute := 0.12
	totalCaloriesBurned := caloriesBurnedPerMinute * float64(durationMinutes)
	return totalCaloriesBurned
}

// велотренажер
func CalculateCaloriesBurnedStationaryBike(durationMinutes int) float64 {
	caloriesBurnedPerMinute := 0.09
	totalCaloriesBurned := caloriesBurnedPerMinute * float64(durationMinutes)
	return totalCaloriesBurned
}

// эллиптический тренажер
func CalculateCaloriesBurnedElliptical(durationMinutes int) float64 {
	caloriesBurnedPerMinute := 0.1
	totalCaloriesBurned := caloriesBurnedPerMinute * float64(durationMinutes)
	return totalCaloriesBurned
}

// степпер
func CalculateCaloriesBurnedStairStepper(durationMinutes int) float64 {
	caloriesBurnedPerMinute := 0.11
	totalCaloriesBurned := caloriesBurnedPerMinute * float64(durationMinutes)
	return totalCaloriesBurned
}

func GetTrainFromDb(db *sql.DB, name string) (Plan, error) {
	// Выполнение SQL-запроса к базе данных
	rows, err := db.Query("SELECT name, trainer, count FROM public.trainers WHERE name = $1", name)
	if err != nil {
		log.Println("query gettrain error")
	}
	defer rows.Close()

	var plan Plan

	for rows.Next() {
		// Создание временных переменных для сканирования значений из строки
		var name string
		var count string
		var trainers string

		// Сканирование значений из строки
		err := rows.Scan(&name, &trainers, &count)
		if err != nil {
			log.Println(err)
			return Plan{}, err
		}
		trainersArr := strings.Fields(trainers[1 : len(trainers)-1])
		countArr := strings.Fields(count[1 : len(count)-1])

		plan.name = name
		plan.count = countArr
		plan.trainers = trainersArr
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return Plan{}, err
	}

	return plan, nil
}

func GetWeightFromDb(db *sql.DB, update tgbotapi.Update) (int, error) {
	// Выполняем запрос к базе данных для получения веса пользователя
	var weight int
	err := db.QueryRow("SELECT weight FROM public.users WHERE id = $1", update.Message.Chat.ID).Scan(&weight)
	if err != nil {
		log.Println("db query error:", err)
		return 0, err
	}

	return weight, nil
}
