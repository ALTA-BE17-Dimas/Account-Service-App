package controllers

import (
	"alta/account-service-app/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func Transfer(db *sql.DB, phoneSender, phoneRecipient string, amount float64) (string, error) {
	// Get a Tx for making transaction requests.
	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Defer a rollback in case anything fails
	defer tx.Rollback() // statement is used to ensure that the prepared statement is closed after it has been executed or if an error occurs.

	// Query the sender's balance
	sqlQuery1 := `SELECT balance FROM users WHERE phone = ? AND deleted_at IS NULL`
	stmt, err := tx.Prepare(sqlQuery1)
	checkErrorPrepare(err)
	defer stmt.Close()

	// Update the sender's balance
	sqlQuery2 := `UPDATE users SET balance = balance + ? WHERE phone = ? AND deleted_at IS NULL`
	stmt, err = tx.Prepare(sqlQuery2)
	checkErrorPrepare(err)

	_, err = stmt.Exec(amount, phoneSender)
	if err != nil {
		return "", fmt.Errorf("failed to update sender's balance: %v", err)
	}

	// Select user id for sender and recipient
	sqlQuery3 := `SELECT id FROM users WHERE phone = ? AND deleted_at IS NULL`
	stmt, err = tx.Prepare(sqlQuery3)
	checkErrorPrepare(err)

	var senderID string
	err = stmt.QueryRow(phoneSender).Scan(&senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("sender's account not found")
		}
		return "", fmt.Errorf("error querying sender's account: %v", err)
	}

	var recipientID string
	err = stmt.QueryRow(phoneRecipient).Scan(&recipientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("recipient's account not found")
		}
		return "", fmt.Errorf("error querying recipient's account: %v", err)
	}

	// Update the recipient's balance
	sqlQuery4 := `UPDATE users SET balance = balance + ? WHERE phone = ? AND deleted_at IS NULL`
	stmt, err = tx.Prepare(sqlQuery4)
	checkErrorPrepare(err)

	_, err = stmt.Exec(amount, phoneRecipient)
	if err != nil {
		return "", fmt.Errorf("failed to update recipient's balance: %v", err)
	}

	// Insert a new row in the transfer_histories table
	sqlQuery5 := `INSERT INTO transfer_histories (user_id_sender, user_id_recipient, amount) VALUES (?, ?, ?)`
	stmt, err = tx.Prepare(sqlQuery5)
	checkErrorPrepare(err)

	_, err = stmt.Exec(senderID, recipientID, amount)
	if err != nil {
		return "", fmt.Errorf("failed to insert transfer history: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	outputStr := "[SUCCESS] Transfer was successful.\n"
	return outputStr, nil
}

func DisplayTransferHistory(db *sql.DB, role, phoneSender string) ([]models.History, error) {
	var sqlQuery string

	sqlQueryRecipient := `
		SELECT th.id, recipient.phone, th.amount, th.created_at
		FROM transfer_histories th
		INNER JOIN users sender ON th.user_id_sender = sender.id
		INNER JOIN users recipient ON th.user_id_recipient = recipient.id
		WHERE sender.phone = ? AND sender.deleted_at IS NULL
		ORDER BY th.id ASC
	`

	sqlQuerySender := `
		SELECT th.id, sender.phone, th.amount, th.created_at
		FROM transfer_histories th
		INNER JOIN users sender ON th.user_id_sender = sender.id
		INNER JOIN users recipient ON th.user_id_recipient = recipient.id
		WHERE recipient.phone = ? AND recipient.deleted_at IS NULL
		ORDER BY th.id ASC
	`

	if role == "sender" {
		sqlQuery = sqlQueryRecipient
	} else {
		sqlQuery = sqlQuerySender
	}

	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		return []models.History{}, fmt.Errorf("failed to prepare SQL statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(phoneSender)
	if err != nil {
		return []models.History{}, fmt.Errorf("failed to query transfer histories: %v", err)
	}
	defer rows.Close()

	var histories []models.History
	for rows.Next() {
		var history models.History
		var createdAt []uint8 // Use []byte to store the raw value

		err := rows.Scan(&history.ID, &history.PhoneNumber, &history.Amount, &createdAt)
		if err != nil {
			log.Printf("failed to scan transfer history: %v", err)
			continue
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

	if err := rows.Err(); err != nil {
		log.Printf("error iterating over transfer histories: %v", err)
	}

	return histories, nil
}

func checkErrorPrepare(err error) error {
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	return nil
}
