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

	outputStr := fmt.Sprintf("\n[SUCCESS] User with ID %d registered successfully.\n\n", id)
	return outputStr, nil
}

func DeleteAccount(db *sql.DB, phoneNumber, password string) (string, error) {
	sqlQuery1 := `DELETE FROM users WHERE phone=?`

	// prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlQuery1)
	if err != nil {
		log.Fatal(err)
	}

	var storedPassword string
	sqlQuery2 := `SELECT password FROM users WHERE phone = ?`
	err = db.QueryRow(sqlQuery2, phoneNumber).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found. Cannot delete account")
		}
		return "", fmt.Errorf("error querying password from database: %v", err)
	}

	err = helpers.ComparePass([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password. Cannot delete account")
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
		outputStr = fmt.Sprintln("[FAIL] User not found. Cannot delete account")
	} else {
		outputStr = fmt.Sprintf("\n[SUCCESS] Account deleted successfully. Rows affected: %d\n", rowsAffected)
	}

	return outputStr, nil
}

func LoginAccount(db *sql.DB, phoneNumber, password string) (string, error) {
	var storedPassword string
	sqlQuery := `SELECT password FROM users WHERE phone = ?`
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		log.Fatal(err)
	}

	err = stmt.QueryRow(phoneNumber).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("password not found")
		}
		return "", fmt.Errorf("error querying password from database: %v", err)
	}

	err = helpers.ComparePass([]byte(storedPassword), []byte(password))
	if err != nil {
		return "", fmt.Errorf("login failed: Invalid password")
	}

	outputStr := fmt.Sprintf("\n[SUCCESS] Login successful!\n")

	return outputStr, nil
}

func ReadOtherAccount(db *sql.DB, phoneNumber string) (models.User, error) {
	sqlStatement := `
		SELECT id, full_name, birth_date, address, email, phone, balance
		FROM users
		WHERE phone=?
	`

	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

	var user models.User
	err = stmt.QueryRow(phoneNumber).Scan(
		&user.ID, &user.FullName, &user.BirthDate, &user.Address, &user.Email, &user.PhoneNumber, &user.Balance,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, fmt.Errorf("user's account not found")
		}
		return models.User{}, fmt.Errorf("error querying user's account: %v", err)
	}

	fmt.Printf("\n[SUCCESS] Account is found.\n\n")
	return user, nil
}

func ReadAccount(db *sql.DB, phoneNumber, password string) (string, error) {
	sqlStatement := 
	`SELECT full_name, identity_number, birth_date, address, email, phone, balance, password 
		FROM users WHERE phone=?`

	// prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	//var fullName, identityNumber, birthDate, address, email, Password string
	//var Balance float64
	var user models.User
	err = stmt.QueryRow(phoneNumber).Scan(&user.FullName, &user.IdentityNumber, &user.BirthDate, &user.Address, &user.Email, &user.PhoneNumber, &user.Balance, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows{
			return "", fmt.Errorf("User not found")
		}
		return "", fmt.Errorf("Error reading account: %v", err)
	}

	err = helpers.ComparePass([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("login failed: Invalid password")
	}

	outputStr := fmt.Sprintf("Account Profil:\nFull Name: %s\nIndentity Number: %s\nBirth of Date: %s\nEmail: %s\nPhone Number: %s\nAddress: %s\nBalance: %f\n", user.FullName, user.IdentityNumber, user.BirthDate, user.Email, user.PhoneNumber, user.Address, user.Balance)

	return outputStr, nil
}