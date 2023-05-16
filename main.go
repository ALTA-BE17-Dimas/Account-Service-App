package main

import (
	"alta/account-service-app/controllers"
	"alta/account-service-app/database"
	"alta/account-service-app/helpers"
	"alta/account-service-app/models"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	// connect to database
	connStr := os.Getenv("DB_CONNECTION")
	db, err := database.DBConnect(connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err.Error())
	}

	defer db.Close()

	menu := `
	---------------------------------
	Menu:

	[1].	Register New Account
	[2].	Login
	[3].	Read Account
	[4].	Update Account
	[5].	Delete Account
	[6].	Top-Up
	[7].	Transfer
	[8].	Display Top-Up History
	[9].	Display Transfer History
	[10].	Read Other Account
	[11].	Exit
	---------------------------------
	`
	var phoneNumber, password string

	loop := true
	for loop {
		fmt.Println(menu)
		fmt.Print("\nEnter menu option: ")
		var option int
		fmt.Scanln(&option)

		switch option {
		case 1:
			// create new user
			newUser := models.User{Balance: 0}

			fmt.Println("\nEnter the data below:")

			// Entering Full Name
			fmt.Print("\nFull Name\t: ")
			newUser.FullName, err = ReadLine()
			checkError(err)

			// Entering Identity Number
			fmt.Print("\nIdentity Number\t: ")
			_, err = fmt.Scanln(&newUser.IdentityNumber)
			checkError(err)

			// Entering Address
			fmt.Print("\nAddress\t\t: ")
			newUser.Address, err = ReadLine()
			checkError(err)

			// Entering Birth Date
			birthDateLoop := true
			for birthDateLoop {
				fmt.Print("\nBirth Date (DD-MM-YYYY)\t: ")
				_, err = fmt.Scanln(&newUser.BirthDate)
				checkError(err)
				birthDateIsValid, _, err := helpers.ValidateDate(newUser.BirthDate)
				if birthDateIsValid {
					birthDateLoop = false
				} else {
					log.Println("Error:", err.Error())
				}
			}

			// Entering Email
			mailLoop := true
			for mailLoop {
				fmt.Print("\nEmail\t\t: ")
				_, err = fmt.Scanln(&newUser.Email)
				checkError(err)
				emailIsValid, err := helpers.ValidateEmail(newUser.Email)
				if emailIsValid {
					mailLoop = false
				} else {
					log.Println("Error:", err.Error())
				}
			}

			// Entering Phone Number
			phoneNumberLoop := true
			for phoneNumberLoop {
				fmt.Print("\nPhone Number\t: ")
				_, err = fmt.Scanln(&newUser.PhoneNumber)
				checkError(err)
				isPhoneNumberValid, err := helpers.ValidatePhoneNumber(newUser.PhoneNumber)
				if isPhoneNumberValid {
					phoneNumberLoop = false
				} else {
					log.Println("Error:", err.Error())
				}
			}

			// Entering Password
			passLoop := true
			for passLoop {
				fmt.Print("\nPassword\t: ")
				_, err = fmt.Scanln(&newUser.Password)
				checkError(err)
				isPassValid, err := helpers.ValidatePassword(newUser.Password)
				if isPassValid {
					passLoop = false
				} else {
					log.Println("Error:", err.Error())
				}
			}

			// registering new user
			str, err := controllers.RegisterAccount(db, newUser)
			if err != nil {
				fmt.Printf("\n")
				log.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Printf("\n")
				fmt.Printf("%s\n", str)
			}
		case 2:
			fmt.Print("\nEnter phone number: ")
			fmt.Scanln(&phoneNumber)
			fmt.Print("\nEnter password: ")
			fmt.Scanln(&password)

			str, err := controllers.LoginAccount(db, phoneNumber, password)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			} else {
				fmt.Println(str)
			}
		case 5:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				str, err := controllers.DeleteAccount(db, phoneNumber, password)
				if err != nil {
					fmt.Printf("\n")
					log.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Printf("\n")
					fmt.Printf("%s\n", str)
				}
			}

		case 10:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				var phoneNumber string
				fmt.Print("Enter other user's phone number\t: ")
				fmt.Scanln(&phoneNumber)
				otherAccount, err := controllers.ReadOtherAccount(db, phoneNumber)
				if err != nil {
					fmt.Printf("\n")
					log.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Println("------------------------------------------")
					fmt.Printf("User Information\n")
					fmt.Println("------------------------------------------")
					fmt.Printf("ID\t\t: %d\n", otherAccount.ID)
					fmt.Printf("Full Name\t: %s\n", otherAccount.FullName)
					fmt.Printf("Birth Date\t: %s\n", otherAccount.BirthDate)
					fmt.Printf("Address\t\t: %s\n", otherAccount.Address)
					fmt.Printf("Email\t\t: %s\n", otherAccount.Email)
					fmt.Printf("Phone Number\t: %s\n", otherAccount.PhoneNumber)
					fmt.Printf("Balance\t\t: %.2f\n", otherAccount.Balance)
					fmt.Println("------------------------------------------")
				}
			}
		case 11:
			loop = false
			fmt.Printf("\nExit program\n")
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
	}
}

func ReadLine() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	str = strings.TrimSuffix(str, "\n")
	return str, nil
}
