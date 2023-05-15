package controllers

import (
	"alta/account-service-app/helpers"
	"alta/account-service-app/models"
	"database/sql"
	"fmt"
)

func UserRegister(db *sql.DB, user models.User) (int64, error) {
	sqlStatement := `
	INSERT INTO users (
		full_name, single_identity_number, birth_date, address, email, phone, password, balance
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	passHashing := helpers.HashPass(user.Password)

	result, err := db.Exec(
		sqlStatement, user.FullName, user.SingleIdentityNumber, user.BirthDate, user.Address,
		user.Email, user.PhoneNumber, passHashing, user.Balance,
	)
	if err != nil {
		return 0, fmt.Errorf("add user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add user: %v", err)
	}

	return id, nil
}
