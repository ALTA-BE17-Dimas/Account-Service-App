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
		full_name, identity_number, birth_date, address, email, phone, password, balance
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
	`

	// prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}

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

//proses login
//mendeklarasikan variabel phone dan pass sebagai parameter fungsi loginUser
func LoginUsers(db *sql.DB, phonenumber, password string) (string, error) {
	
	//query untuk memeriksa kecocokkan username dan password
	//mendefinisikan query
	//query:= "SELECT id, phone, password FROM users WHERE phone = ? LIMIT 1"
	query := "SELECT id, phone, password FROM users WHERE phone = ? LIMIT 1"

	// prepared statement from the SQL statement before executed
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	var user models.User 
	//eksekusi pemanggilan query kedatabase
	err = stmt.QueryRow(phonenumber).Scan(&user.ID, &user.PhoneNumber, &user.Password)
	if err!= nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("login failed: Invalid phone")
		} else {
			return "", err
		}
	}
	//mengembalikan objek user yang berhasil ditemukan
	
	//compare password dengan hash password
	err = helpers.ComparePass([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("login failed: Invalid password")
	}

	outputStr := fmt.Sprint("Login successful!")
	return outputStr, nil
	
}