package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	helper "tgbotfitnes/handler"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "    " // Четыре пробела
	dbname   = "tgbotfitn"
)

func DbConnectin() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db

}

func InsertUser(db *sql.DB, newUser *helper.User) error {
	_, err := db.Exec("INSERT INTO public.users (id, name, weight, height) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO NOTHING", newUser.ID, newUser.Name, newUser.Weight, newUser.Height)
	_, err = db.Exec("INSERT INTO public.summary (id,quantity, calorie) VALUES ($1, $2, $3) ON CONFLICT (id) DO NOTHING", newUser.ID, 0, 0)
	return err
}

// Метод для обновления существующего пользователя в базе данных
func UpdateUser(db *sql.DB, currentUser *helper.User) error {
	_, err := db.Exec("UPDATE users SET name = $1, weight = $2, height = $3 WHERE id = $4", currentUser.Name, currentUser.Weight, currentUser.Height, currentUser.ID)
	return err
}

func PrintTable(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, weight, height FROM public.users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user helper.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Weight, &user.Height); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s, Weight: %d, Height: %d\n", user.ID, user.Name, user.Weight, user.Height)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func InfoAboutUser(db *sql.DB, id int64) (helper.User, error) {
	row, err := db.Query("SELECT id, name, weight, height FROM public.users WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return helper.User{}, err
	}
	defer row.Close()

	var user helper.User
	if row.Next() {
		if err := row.Scan(&user.ID, &user.Name, &user.Weight, &user.Height); err != nil {
			log.Println(err)
			return helper.User{}, err
		}
	} else {
		return helper.User{}, errors.New("User not found")
	}

	return user, nil
}

func SummaryInfo(db *sql.DB, id int64) (float64, int) {
	row, err := db.Query("SELECT calorie, quantity FROM public.summary WHERE id = $1", id)
	if err != nil {
		log.Println(err)
		return 0, 0
	}
	defer row.Close()

	var colorie float64
	var quantity int
	if row.Next() {
		if err := row.Scan(&colorie, &quantity); err != nil {
			log.Println(err)
			return colorie, quantity
		}
	} else {
		return 0, 0
	}
	log.Println("SummaryInfo corrected")
	return colorie, quantity
}
