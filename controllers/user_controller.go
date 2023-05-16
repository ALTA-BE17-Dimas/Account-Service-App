package controllers

import (
	"alta/account-service-app/helpers"
	"alta/account-service-app/models"
	"database/sql"
	"fmt"
	"log"
)

func RegisterAccount(db *sql.DB, user models.User) (string, error) {
	sqlStatement := `
	INSERT INTO users (
		full_name, identity_number, birth_date, address, email, phone, password, balance
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	// prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	// validating email format
	emailIsValid, err := helpers.ValidateEmail(user.Email)
	if !emailIsValid {
		log.Println("Error:", err.Error())
	}

	// validating birth date format
	valid, birthDate, err := helpers.ValidateDate(user.BirthDate)
	if !valid {
		log.Println("Error:", err.Error())
	}

	// validating password
	passIsValid, err := helpers.ValidatePassword(user.Password)
	passHashing := ""
	if passIsValid {
		// hashing the password
		passHashing = helpers.HashPass(user.Password)
	} else {
		log.Println("Error:", err.Error())
	}

	// insert new data to database
	result, err := stmt.Exec(
		user.FullName, user.IdentityNumber, birthDate, user.Address,
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

func DeleteAccount(db *sql.DB, phoneNumber string) (string, error) {
	sqlStatement := `DELETE FROM users WHERE phone=?`

	// prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	result, err := stmt.Exec(phoneNumber)
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}

	outputStr := ""

	if rowsAffected == 0 {
		outputStr = fmt.Sprintln("User not found. Cannot delete account")
	} else {
		outputStr = fmt.Sprintf("User deleted successfully. Rows affected: %d\n", rowsAffected)
	}

	return outputStr, nil
}
