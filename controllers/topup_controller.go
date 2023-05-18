package controllers

import (
	"alta/account-service-app/models"
	"database/sql"
	"fmt"
	"time"
	"log"
)

func Topup(db *sql.DB, phoneNumber string, amount float64) (string, error) {
	// making transaction requests.
	transaction, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Defer a rollback in case anything fails
	defer transaction.Rollback() // statement is used to ensure that the prepared statement is closed after it has been executed or if an error occurs.
	
	// Query the user's balance
	sqlQuery1 := `SELECT balance FROM users WHERE phone = ?`
	stmt, err := transaction.Prepare(sqlQuery1)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()
	
	var balance float64
	err = stmt.QueryRow(phoneNumber).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return"", fmt.Errorf("user's account not found")
		}
		return"", fmt.Errorf("failed to query user's account: %v", err)
	}

	//update user balance
	sqlQuery2 := `UPDATE users SET balance = balance + ? WHERE phone = ?`
	stmt, err = transaction.Prepare(sqlQuery2)
	checkErrorPrepare(err)
		
	_, err = stmt.Exec(amount, phoneNumber)
	if err != nil {
		return "", fmt.Errorf("failed to update user balance: %v", err)
	}
	
	// Get user ID
	sqlQuery3 := `SELECT id FROM users WHERE phone = ?`
	stmt, err = transaction.Prepare(sqlQuery3)
	checkErrorPrepare(err)
	
	var userID string
	err = stmt.QueryRow(phoneNumber).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return"", fmt.Errorf("user's account not found")
		}
		return"", fmt.Errorf("failed to query user's account: %v", err)
	}

	// Insert a new row in the topup_histories table
	sqlQuery4 := `INSERT INTO top_up_histories (user_id, amount, created_at) VALUES (?, ?, NOW())`
	stmt, err = transaction.Prepare(sqlQuery4)
	checkErrorPrepare(err)

	_, err = stmt.Exec(userID, amount)
	if err != nil {
		return "", fmt.Errorf("failed to insert transfer history: %v", err)
	}

	// commit transaction
	if err = transaction.Commit(); err != nil {
		return"", fmt.Errorf("failed to commit transaction: %v", err)
	}

	outputStr := "\n[SUCCESS] Top Up was successfull.\n"
	return outputStr, nil
}

func DisplayTopupHistories(db *sql.DB, phoneNumber string) ([]models.TopUpHistory, error) {
	// Query top-up histories for a specific user
	sqlQuery := `
		SELECT th.user_id, th.amount, th.created_at
		FROM top_up_histories AS th
		JOIN users AS u ON th.user_id = u.id
		WHERE u.phone = ?
		ORDER BY th.created_at DESC
	`
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare SQL query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	histories := []models.TopUpHistory{}
	for rows.Next() {
		var history models.TopUpHistory
		var createdAt []uint8 // Use []byte to store the raw value
		
		err := rows.Scan(&history.UserID, &history.Amount, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan top up history: %v", err)
		}
		// Parse the createdAt value into a time.Time variable
		createdAtStr := string(createdAt)
		history.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			log.Printf("failed to parse created_at value: %v", err)
			continue
		}
		histories = append(histories, history)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to prepare insert statement: %v", err)
	}

	return histories, nil
}

	