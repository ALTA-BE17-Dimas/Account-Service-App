package main

import (
	"alta/account-service-app/controllers"
	"alta/account-service-app/database"
	"alta/account-service-app/helpers"
	"alta/account-service-app/models"
	"fmt"
	"log"
	"os"
	"time"
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
	[2].	Login Account
	[3].	Read Account
	[4].	Update Account
	[5].	Delete Account
	[6].	Top-Up Account
	[7].	Transfer
	[8].	Display Top-Up History
	[9].	Display Transfer History
	[10].	Read Other Account
	[11].	Log Out Account
	[12].	Exit
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
			if (phoneNumber != "") || (password != "") {
				fmt.Printf("\nYou have to log out first!\n")
			} else {
				// create new user
				newUser := models.User{Balance: 0}

				fmt.Println("\nEnter the data below:")

				// Entering Full Name
				fmt.Print("\nFull Name\t: ")
				newUser.FullName, err = helpers.ReadLine()
				checkError(err)

				// Entering Identity Number
				fmt.Print("\nIdentity Number\t: ")
				_, err = fmt.Scanln(&newUser.IdentityNumber)
				checkError(err)

				// Entering Address
				fmt.Print("\nAddress\t\t: ")
				newUser.Address, err = helpers.ReadLine()
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
			}

		case 2:
			fmt.Print("\nEnter phone number: ")
			fmt.Scanln(&phoneNumber)
			fmt.Print("\nEnter password: ")
			fmt.Scanln(&password)

			str, err := controllers.LoginAccount(db, phoneNumber, password)
			if err != nil {
				fmt.Println("")
				log.Printf("[FAIL] %s\n", err.Error())
			} else {
				fmt.Println("")
				log.Printf("%s\n", str)
			}

		case 3:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				user, err := controllers.ReadAccount(db, phoneNumber, password)
				fmt.Print("\n")
				if err != nil {
					log.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Printf("\n%s\n", user)
				}
			}
		case 4:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				updateMenu := `
				Select the section you want to update:

				[1].	Full Name
				[2].	Birth Date
				[3].	Address
				[4].	Email
				[5].	Phone Number
				[6].	Password
				[7].	Finish Update
				`

				updateLoop := true
			outer:
				for updateLoop {
					fmt.Println(updateMenu)
					fmt.Print("\nEnter update menu option: ")
					var option int
					fmt.Scanln(&option)

					switch option {
					case 1:
						fmt.Println("\nPress (x) to go back")
						fmt.Print("\nEnter new full name\t: ")
						var newValue string
						newValue, err = helpers.ReadLine()
						checkError(err)

						if newValue == "x" {
							continue outer
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Full name", "full_name", newValue)
						if err != nil {
							log.Printf("Error: %s\n", err.Error())
						} else {
							fmt.Printf("\n%s\n", str)
						}
					case 2:
						var newValue string
						var newBirthDate time.Time
						birthDateUpdateLoop := true
						for birthDateUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new birth date\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue outer
							}

							birthDateIsValid, birthDate, err := helpers.ValidateDate(newValue)
							if birthDateIsValid {
								newBirthDate = birthDate
								birthDateUpdateLoop = false
							} else {
								log.Println("Error:", err.Error())
							}
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Birth date", "birth_date", newBirthDate)
						if err != nil {
							log.Printf("Error: %s\n", err.Error())
						} else {
							fmt.Printf("\n%s\n", str)
						}
					case 3:
						fmt.Println("\nPress (x) to go back")
						fmt.Print("\nEnter new address\t: ")
						var newValue string
						newValue, err = helpers.ReadLine()
						checkError(err)

						if newValue == "x" {
							continue outer
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Address", "address", newValue)
						if err != nil {
							log.Printf("Error: %s\n", err.Error())
						} else {
							fmt.Printf("\n%s\n", str)
						}
					case 4:
						var newValue string
						mailUpdateLoop := true
						for mailUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new email\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue outer
							}

							emailIsValid, err := helpers.ValidateEmail(newValue)
							if emailIsValid {
								mailUpdateLoop = false
							} else {
								log.Println("Error:", err.Error())
							}
						}
						str, err := controllers.UpdateAccount(db, phoneNumber, "Email", "email", newValue)
						if err != nil {
							log.Printf("Error: %s\n", err.Error())
						} else {
							fmt.Printf("\n%s\n", str)
						}
					case 5:
						var newValue string
						phoneNumberUpdateLoop := true
						for phoneNumberUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new phone number\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue outer
							}

							isPhoneNumberValid, err := helpers.ValidatePhoneNumber(newValue)
							if isPhoneNumberValid {
								phoneNumberUpdateLoop = false
							} else {
								log.Println("Error:", err.Error())
							}
						}
						str, err := controllers.UpdateAccount(db, phoneNumber, "Phone number", "phone", newValue)
						if err != nil {
							log.Printf("Error: %s\n", err.Error())
						} else {
							fmt.Printf("\n%s\n", str)
						}
					case 6:
						var newValue string
						passHashing := ""
						passUpdateLoop := true
						for passUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new password\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue outer
							}

							isPassValid, err := helpers.ValidatePassword(newValue)
							if isPassValid {
								passHashing = helpers.HashPass(newValue)
								passUpdateLoop = false
							} else {
								log.Println("Error:", err.Error())
							}
						}
						str, err := controllers.UpdateAccount(db, phoneNumber, "Password", "password", passHashing)
						if err != nil {
							log.Printf("Error: %s\n", err.Error())
						} else {
							fmt.Printf("\n%s\n", str)
						}
					case 7:
						updateLoop = false
						fmt.Printf("\nUpdate complete.\n")
					}
				}
			}
		case 5:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				str, err := controllers.DeleteAccount(db, phoneNumber, password)
				if err != nil {
					fmt.Println("")
					log.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Println("")
					log.Printf("%s", str)
				}
			}

		case 7:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				var phoneNumberRecipient string
				var transferAmount float64
				fmt.Print("\nEnter recipient's phone number: ")
				fmt.Scanln(&phoneNumberRecipient)
				fmt.Print("\nEnter transfer amount: ")
				fmt.Scanln(&transferAmount)
				str, err := controllers.Transfer(db, phoneNumber, phoneNumberRecipient, transferAmount)
				if err != nil {
					fmt.Println("")
					log.Printf("Error: %s\n", err.Error())
				} else {
					fmt.Println("")
					log.Printf("%s\n", str)
				}
			}

		case 9:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				histories := controllers.DisplayTransferHistory(db, phoneNumber)
				fmt.Printf("\n")
				fmt.Println("-----------------------------------------")
				fmt.Printf("Your Transfer History: \n")
				fmt.Println("-----------------------------------------")
				transferCounter := 0
				for _, value := range histories {
					transferCounter++
					fmt.Printf("%+v\n", value)
				}
				fmt.Println("Count:", transferCounter)
			}

		case 10:
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\nYou have to login first!\n")
			} else {
				var otherPhoneNumber string
				fmt.Print("\nEnter other user's phone number\t: ")
				fmt.Scanln(&otherPhoneNumber)
				if otherPhoneNumber == phoneNumber {
					fmt.Println("")
					log.Printf("Choose option 3 to see your account information")
				} else {
					str, err := controllers.ReadOtherAccount(db, otherPhoneNumber)
					if err != nil {
						fmt.Println("")
						log.Printf("Error: %s\n", err.Error())
					} else {
						fmt.Printf("\n%s\n", str)
					}
				}
			}
		case 11:
			str := controllers.LogOutAccount(&phoneNumber, &password)
			fmt.Println(str)
		case 12:
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
