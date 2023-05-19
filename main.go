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
		log.Fatal("\033[91mError:\033[0m", err.Error())
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

firstOuter:
	for loop {
		fmt.Println(menu)

		if (phoneNumber == "") || (password == "") {
			fmt.Printf("\n\033[91mYou are not login yet!\033[0m\n")
		} else {
			id, name := controllers.GetAccountInfo(db, phoneNumber)
			fmt.Printf("\n\033[93mYou are login as (%s - %s)\033[0m\n", id, name)
		}

		fmt.Print("\nEnter menu option: ")
		var option int
		fmt.Scanln(&option)

		switch option {
		case 1: // Register new user account
			if (phoneNumber != "") || (password != "") {
				fmt.Printf("\n\033[91mYou have to log out first!\033[0m\n")
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
						log.Printf("\033[91mError: %s\033[0m\n", err.Error())
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
						log.Printf("\033[91mError: %s\033[0m\n", err.Error())
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
						log.Printf("\033[91mError: %s\033[0m\n", err.Error())
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
						log.Printf("\033[91mError: %s\033[0m\n", err.Error())
					}
				}

				// registering new user
				str, err := controllers.RegisterAccount(db, newUser)
				if err != nil {
					fmt.Printf("\n")
					log.Printf("\033[91mError: %s\033[0m\n", err.Error())
				} else {
					fmt.Printf("\n")
					log.Printf("\033[92m%s\033[0m\n", str)
				}
			}

		case 2: // Login user account
			if (phoneNumber != "") || (password != "") {
				fmt.Printf("\n\033[91mYou have to log out first!\033[0m\n")
			} else {
				fmt.Print("\nEnter phone number: ")
				fmt.Scanln(&phoneNumber)
				fmt.Print("\nEnter password: ")
				fmt.Scanln(&password)

				str, err := controllers.LoginAccount(db, phoneNumber, password)
				if err != nil {
					_ = controllers.LogOutAccount(&phoneNumber, &password)
					fmt.Println("")
					log.Printf("\033[91m[FAIL] %s\033[0m\n", err.Error())
				} else {
					fmt.Println("")
					log.Printf("\033[92m%s\033[0m\n", str)
				}
			}

		case 3: // Read user account
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
			} else {
				user, err := controllers.ReadAccount(db, phoneNumber, password)
				fmt.Print("\n")
				if err != nil {
					log.Printf("\033[91mError: %s\033[0m\n", err.Error())
				} else {
					log.Printf("\n\033[95m%s\033[0m\n", user)
				}
			}

		case 4: // Update user account
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
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
			secondOuter:
				for updateLoop {
					fmt.Println(updateMenu)
					fmt.Print("\nEnter update menu option: ")
					var option int
					fmt.Scanln(&option)

					switch option {
					case 1: // Update user account name
						fmt.Println("\nPress (x) to go back")
						fmt.Print("\nEnter new full name\t: ")
						var newValue string
						newValue, err = helpers.ReadLine()
						checkError(err)

						if newValue == "x" {
							continue secondOuter
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Full name", "full_name", newValue)
						if err != nil {
							log.Printf("\033[91mError: %s\033[0m\n", err.Error())
						} else {
							log.Printf("\n\033[92m%s\033[0m\n", str)
						}

					case 2: // Update user account birth date
						var newValue string
						var newBirthDate time.Time
						birthDateUpdateLoop := true
						for birthDateUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new birth date\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue secondOuter
							}

							birthDateIsValid, birthDate, err := helpers.ValidateDate(newValue)
							if birthDateIsValid {
								newBirthDate = birthDate
								birthDateUpdateLoop = false
							} else {
								log.Printf("\033[91mError: %s\033[0m\n", err.Error())
							}
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Birth date", "birth_date", newBirthDate)
						if err != nil {
							log.Printf("\033[91mError: %s\033[0m\n", err.Error())
						} else {
							log.Printf("\n\033[92m%s\033[0m\n", str)
						}

					case 3: // Update user account address
						fmt.Println("\nPress (x) to go back")
						fmt.Print("\nEnter new address\t: ")
						var newValue string
						newValue, err = helpers.ReadLine()
						checkError(err)

						if newValue == "x" {
							continue secondOuter
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Address", "address", newValue)
						if err != nil {
							log.Printf("\033[91mError: %s\033[0m\n", err.Error())
						} else {
							log.Printf("\n\033[92m%s\033[0m\n", str)
						}

					case 4: // Update user account email
						var newValue string
						mailUpdateLoop := true
						for mailUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new email\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue secondOuter
							}

							emailIsValid, err := helpers.ValidateEmail(newValue)
							if emailIsValid {
								mailUpdateLoop = false
							} else {
								log.Printf("\033[91mError: %s\033[0m\n", err.Error())
							}
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Email", "email", newValue)
						if err != nil {
							log.Printf("\033[91mError: %s\033[0m\n", err.Error())
						} else {
							log.Printf("\n\033[92m%s\033[0m\n", str)
						}

					case 5: // Update user account phone number
						var newValue string
						phoneNumberUpdateLoop := true
						for phoneNumberUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new phone number\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue secondOuter
							}

							isPhoneNumberValid, err := helpers.ValidatePhoneNumber(newValue)
							if isPhoneNumberValid {
								phoneNumberUpdateLoop = false
							} else {
								log.Printf("\033[91mError: %s\033[0m\n", err.Error())
							}
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Phone number", "phone", newValue)
						if err != nil {
							log.Printf("\033[91mError: %s\033[0m\n", err.Error())
						} else {
							_ = controllers.LogOutAccount(&phoneNumber, &password)
							log.Printf("\n\033[92m%s\033[0m\n", str)
							continue firstOuter
						}

					case 6: // Update user account password
						var newValue string
						var passHashing string
						passUpdateLoop := true
						for passUpdateLoop {
							fmt.Println("\nPress (x) to go back")
							fmt.Print("\nEnter new password\t: ")
							_, err = fmt.Scanln(&newValue)
							checkError(err)

							if newValue == "x" {
								continue secondOuter
							}

							isPassValid, err := helpers.ValidatePassword(newValue)
							if isPassValid {
								passHashing = helpers.HashPass(newValue)
								passUpdateLoop = false
							} else {
								log.Printf("\n\033[91mError: %s\033[0m\n", err.Error())
							}
						}

						str, err := controllers.UpdateAccount(db, phoneNumber, "Password", "password", passHashing)
						if err != nil {
							log.Printf("\033[91mError: %s\033[0m\n", err.Error())
						} else {
							_ = controllers.LogOutAccount(&phoneNumber, &password)
							log.Printf("\n\033[92m%s\033[0m\n", str)
							continue firstOuter
						}

					case 7: // Exit from update
						updateLoop = false
						fmt.Printf("\nUpdate complete.\n")
					}
				}
			}

		case 5: // Delete user account
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
			} else {
				str, err := controllers.DeleteAccount(db, phoneNumber, password)
				if err != nil {
					fmt.Println("")
					log.Printf("\033[91mError: %s\033[0m\n", err.Error())
				} else {
					_ = controllers.LogOutAccount(&phoneNumber, &password)
					fmt.Println("")
					log.Printf("\033[92m%s\033[0m\n", str)
					continue firstOuter
				}
			}

		case 6: // Top-up balance user account
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
			} else {
				var topupAmount float64
				fmt.Print("\nEnter top amount: ")
				fmt.Scanln(&topupAmount)
				str, err := controllers.Topup(db, phoneNumber, topupAmount)
				if err != nil {
					fmt.Printf("\n")
					log.Printf("\033[91mError: %s\033[0m\n", err.Error())
				} else {
					fmt.Printf("\n")
					log.Printf("\033[92m%s\033[0m\n", str)
				}
			}

		case 7: // Transfer balance to another account
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
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
					log.Printf("\033[91mError: %s\033[0m\n", err.Error())
				} else {
					fmt.Println("")
					log.Printf("\033[92m%s\033[0m\n", str)
				}
			}

		case 8: //Display top-up history
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
			} else {
				histories, err := controllers.DisplayTopupHistories(db, phoneNumber)
				checkError(err)
				fmt.Printf("\n")
				fmt.Println("-----------------------------------------")
				fmt.Printf("Your top-up history: \n")
				fmt.Println("-----------------------------------------")
				topupCounter := 0
				// Print top-up histories
				for _, history := range histories {
					topupCounter++
					fmt.Printf("User ID\t: %s\n", history.UserID)
					fmt.Printf("Amount\t: %.2f\n", history.Amount)
					fmt.Printf("Time\t: %s\n", history.CreatedAt.Format("2006-01-02 15:04:05"))
					fmt.Println("-----------------------------------------")
				}
				fmt.Println("Count:", topupCounter)
			}

		case 9: // Display transfer history user account
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
			} else {
				tfHistory := `
				Display transfer history as:

				[1].	Sender
				[2].	Recipient
				[3].	Exit
				`
				tfHistoryLoop := true
				for tfHistoryLoop {
					fmt.Println(tfHistory)
					fmt.Print("\nEnter transfer history option: ")
					var option int
					fmt.Scanln(&option)

					switch option {
					case 1:
						histories, err := controllers.DisplayTransferHistory(db, "sender", phoneNumber)
						checkError(err)
						fmt.Printf("\n")
						fmt.Println("-----------------------------------------")
						fmt.Printf("Your transfer history as sender: \n")
						fmt.Println("-----------------------------------------")
						transferCounter := 0
						for _, value := range histories {
							transferCounter++
							fmt.Printf(
								"transfer_id: %s, phone_recipient: %s, amount: %.2f, transaction_time: %s\n",
								value.ID, value.PhoneNumber, value.Amount, value.CreatedAt.Format("2006-01-02 15:04:05"),
							)
						}
						fmt.Println("Count:", transferCounter)

					case 2:
						histories, err := controllers.DisplayTransferHistory(db, "recipient", phoneNumber)
						checkError(err)
						fmt.Printf("\n")
						fmt.Println("-----------------------------------------")
						fmt.Printf("Your transfer history as recipient: \n")
						fmt.Println("-----------------------------------------")
						transferCounter := 0
						for _, value := range histories {
							transferCounter++
							fmt.Printf(
								"transfer_id: %s, phone_sender: %s, amount: %.2f, transaction_time: %s\n",
								value.ID, value.PhoneNumber, value.Amount, value.CreatedAt.Format("2006-01-02 15:04:05"),
							)
						}
						fmt.Println("Count:", transferCounter)

					case 3:
						tfHistoryLoop = false
					}
				}
			}

		case 10: // Read other user account
			if (phoneNumber == "") || (password == "") {
				fmt.Printf("\n\033[91mYou have to login first!\033[0m\n")
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
						log.Printf("\033[91mError: %s\033[0m\n", err.Error())
					} else {
						fmt.Printf("\033[95m\n%s\033[0m\n", str)
					}
				}
			}

		case 11: // Logout from user account
			str := controllers.LogOutAccount(&phoneNumber, &password)
			fmt.Printf("\033[92m%s\033[0m\n", str)

		case 12: // Exit from program
			loop = false
			fmt.Printf("\nExit program\n")
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Printf("\033[91mError: %s\033[0m\n", err.Error())
	}
}
