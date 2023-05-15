package controllers

import (
	"alta/account-service-app/helpers"
	"alta/account-service-app/models"
	"database/sql"
	"fmt"
	"log"
)

func UserRegister(db *sql.DB, user models.User) (string, error) {
	sqlStatement := `
	INSERT INTO users (
		full_name, single_identity_number, birth_date, address, email, phone, password, balance
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	// validating email format
	if emailIsValid := helpers.ValidateEmail(user.Email); !emailIsValid {
		log.Fatal("email format is invalid")
	}

	// validating birth date format
	valid, birthDate, _ := helpers.ValidateDate(user.BirthDate)
	if !valid {
		log.Fatal("date is invalid. date format expected for input is (DD-MM-YYYY)")
	}

	// validating password
	passIsValid := helpers.ValidatePassword(user.Password)
	passHashing := ""
	if passIsValid {
		// hashing the password
		passHashing = helpers.HashPass(user.Password)
	} else {
		log.Fatal("password should contain lowercase, uppercase, special character, and the length is more than 7")
	}

	// insert new data to database
	result, err := db.Exec(
		sqlStatement, user.FullName, user.SingleIdentityNumber, birthDate, user.Address,
		user.Email, user.PhoneNumber, passHashing, user.Balance,
	)
	if err != nil {
		return "", fmt.Errorf("add user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("add user: %v", err)
	}

	outputStr := fmt.Sprintf("User with ID %d registered successfully.", id)
	return outputStr, nil
}
