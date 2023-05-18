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
		id, full_name, identity_number, birth_date, address, email, phone, password, balance
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);
	`

	// Prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlStatement)
	checkErrorPrepare(err)
	defer stmt.Close()

	// Validating email format
	emailIsValid, err := helpers.ValidateEmail(user.Email)
	if !emailIsValid {
		log.Println("Error:", err.Error())
	}

	// Validating birth date format
	valid, birthDate, err := helpers.ValidateDate(user.BirthDate)
	if !valid {
		log.Println("Error:", err.Error())
	}

	// Validating password
	passIsValid, err := helpers.ValidatePassword(user.Password)
	passHashing := ""
	if passIsValid {
		// Hashing the password
		passHashing = helpers.HashPass(user.Password)
	} else {
		log.Println("Error:", err.Error())
	}

	// Generate user ID
	userID := helpers.GenerateID()

	// Insert new data to database
	_, err = stmt.Exec(
		userID, user.FullName, user.IdentityNumber, birthDate, user.Address,
		user.Email, user.PhoneNumber, passHashing, user.Balance,
	)
	if err != nil {
		return "", fmt.Errorf("add user: %v", err)
	}

	outputStr := "\n[SUCCESS] Account registered successfully.\n\n"
	return outputStr, nil
}

func DeleteAccount(db *sql.DB, phoneNumber, password string) (string, error) {
	sqlQuery1 := `DELETE FROM users WHERE phone=?`

	// Prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlQuery1)
	checkErrorPrepare(err)
	defer stmt.Close()

	result, err := stmt.Exec(phoneNumber)
	if err != nil {
		return "", fmt.Errorf("failed to delete account: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("failed to delete account: %v", err)
	}

	outputStr := ""

	if rowsAffected == 0 {
		outputStr = fmt.Sprintln("[FAIL] User not found. Cannot delete account")
	} else {
		outputStr = fmt.Sprintf("[SUCCESS] Account deleted successfully. Rows affected: %d\n", rowsAffected)
	}

	return outputStr, nil
}

func LoginAccount(db *sql.DB, phoneNumber, password string) (string, error) {
	var storedPassword string

	// Prepare the SQL statement
	sqlQuery := `SELECT password FROM users WHERE phone = ?`
	stmt, err := db.Prepare(sqlQuery)
	checkErrorPrepare(err)
	defer stmt.Close()

	// Query the password from the database
	err = stmt.QueryRow(phoneNumber).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("login failed, invalid phone number")
		}
		return "", fmt.Errorf("error querying password from database: %v", err)
	}

	// Compare the stored password with the provided password
	err = helpers.ComparePass([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("login failed, invalid password")
	}

	outputStr := "[SUCCESS] Login successful!\n"

	return outputStr, nil
}

func ReadOtherAccount(db *sql.DB, phoneNumber string) (string, error) {
	sqlStatement := `
		SELECT id, full_name, birth_date, address, email, phone
		FROM users
		WHERE phone=?
	`

	stmt, err := db.Prepare(sqlStatement)
	checkErrorPrepare(err)
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(phoneNumber).Scan(
		&user.ID, &user.FullName, &user.BirthDate, &user.Address, &user.Email, &user.PhoneNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user's account not found")
		}
		return "", fmt.Errorf("error querying user's account: %v", err)
	}

	outputStr := fmt.Sprintln("-----------------------------------------")
	outputStr += fmt.Sprintln("Account Information")
	outputStr += fmt.Sprintln("-----------------------------------------")
	outputStr += fmt.Sprintf("ID\t\t: %s\n", user.ID)
	outputStr += fmt.Sprintf("Full Name\t: %s\n", user.FullName)
	outputStr += fmt.Sprintf("Birth Date\t: %s\n", user.BirthDate)
	outputStr += fmt.Sprintf("Address\t\t: %s\n", user.Address)
	outputStr += fmt.Sprintf("Email\t\t: %s\n", user.Email)
	outputStr += fmt.Sprintf("Phone Number\t: %s\n", user.PhoneNumber)
	outputStr += fmt.Sprintln("-----------------------------------------")

	fmt.Println("")
	log.Printf("[SUCCESS] Account is found.\n")

	return outputStr, nil
}

func ReadAccount(db *sql.DB, phoneNumber, password string) (string, error) {
	sqlStatement :=
		`SELECT id, full_name, identity_number, birth_date, address, email, phone, balance
		FROM users WHERE phone=?`

	// prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlStatement)
	checkErrorPrepare(err)
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(phoneNumber).Scan(&user.ID, &user.FullName, &user.IdentityNumber, &user.BirthDate, &user.Address, &user.Email, &user.PhoneNumber, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("error reading account: %v", err)
	}

	outputStr := fmt.Sprintln("-----------------------------------------")
	outputStr += fmt.Sprintln("Your Account Information")
	outputStr += fmt.Sprintln("-----------------------------------------")
	outputStr += fmt.Sprintf("ID\t\t: %s\n", user.ID)
	outputStr += fmt.Sprintf("Full Name\t: %s\n", user.FullName)
	outputStr += fmt.Sprintf("Identity Number\t: %s\n", user.IdentityNumber)
	outputStr += fmt.Sprintf("Birth Date\t: %s\n", user.BirthDate)
	outputStr += fmt.Sprintf("Address\t\t: %s\n", user.Address)
	outputStr += fmt.Sprintf("Email\t\t: %s\n", user.Email)
	outputStr += fmt.Sprintf("Phone Number\t: %s\n", user.PhoneNumber)
	outputStr += fmt.Sprintf("Balance\t\t: %.2f\n", user.Balance)
	outputStr += fmt.Sprintln("-----------------------------------------")

	return outputStr, nil
}

func LogOutAccount(phoneNumber, password *string) string {
	*phoneNumber = ""
	*password = ""

	return "\n[SUCCESS] Log out success"
}

func UpdateAccount(db *sql.DB, phoneNumber, updateOption, column string, value interface{}) (string, error) {
	sqlQuery := "UPDATE users SET " + column + " = ? WHERE phone = ?"

	stmt, err := db.Prepare(sqlQuery)
	checkErrorPrepare(err)
	defer stmt.Close()

	_, err = stmt.Exec(value, phoneNumber)
	if err != nil {
		return "", fmt.Errorf("failed to update %s: %v", updateOption, err)
	}

	outputStr := fmt.Sprintf("[SUCCESS] %s updated successfully", updateOption)
	return outputStr, nil
}

func GetAccountInfo(db *sql.DB, phoneNumber string) (string, string) {
	sqlQuery := `SELECT id, full_name FROM users WHERE phone = ?`
	stmt, err := db.Prepare(sqlQuery)
	checkErrorPrepare(err)
	defer stmt.Close()

	var user models.User
	_ = stmt.QueryRow(phoneNumber).Scan(&user.ID, &user.FullName)

	return user.ID, user.FullName
}
