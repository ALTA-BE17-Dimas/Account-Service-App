package main

import (
	"alta/account-service-app/controllers"
	"alta/account-service-app/database"
	//"alta/account-service-app/models"
	"fmt"
	"log"
	"os"
)

func main() {
	// connect to database
	connStr := os.Getenv("DB_CONNECTION")
	db, err := database.DBConnect(connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err.Error())
	}

	// create new user
	// newUser := models.User{
	// 	FullName:       "Dimas Pradana",
	// 	IdentityNumber: "1234567890",
	// 	BirthDate:      "07-07-1990",
	// 	Address:        "23 Brookmead, Thornbury England",
	// 	Email:          "dimasprd11@gmail.com",
	// 	PhoneNumber:    "089978234511",
	// 	Password:       "fsaghdRTUbjkWhs1@&",
	// 	Balance:        10000,
	// }

	// // register new user
	// str, err := controllers.UserRegister(db, newUser)
	// if err != nil {
	// 	log.Printf("Error: %s\n", err.Error())
	// } else {
	// 	fmt.Println(str)
	// }

	defer db.Close()
	// login user
	user, err := controllers.LoginUsers(db, "089978234511", "fsaghdRTUbjkWhs1@&")
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	} else {
		fmt.Println(user)
	}
}
