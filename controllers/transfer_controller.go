package controllers

import (
	"database/sql"
	"fmt"
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
	sqlQuery1 := `SELECT balance FROM users WHERE phone = ?`
	stmt, err := tx.Prepare(sqlQuery1)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	// Update the sender's balance
	sqlQuery2 := `UPDATE users SET balance = balance + ? WHERE phone = ?`
	stmt, err = tx.Prepare(sqlQuery2)
	if err != nil {
		return "", fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(amount, phoneSender)
	if err != nil {
		return "", fmt.Errorf("failed to update sender's balance: %v", err)
	}

	// Select user id for sender and recipient
	sqlQuery3 := `SELECT id FROM users WHERE phone = ?`
	stmt, err = tx.Prepare(sqlQuery3)
	if err != nil {
		return "", fmt.Errorf("failed to prepare query: %v", err)
	}
	defer stmt.Close()

	var senderID string
	var recipientID string
	err = stmt.QueryRow(phoneSender).Scan(&senderID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("sender's account not found")
		}
		return "", fmt.Errorf("error querying sender's account: %v", err)
	}
	err = stmt.QueryRow(phoneRecipient).Scan(&recipientID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("recipient's account not found")
		}
		return "", fmt.Errorf("error querying recipient's account: %v", err)
	}

	// Update the recipient's balance
	sqlQuery4 := `UPDATE users SET balance = balance + ? WHERE phone = ?`
	stmt, err = tx.Prepare(sqlQuery4)
	if err != nil {
		return "", fmt.Errorf("failed to prepare update statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(amount, phoneRecipient)
	if err != nil {
		return "", fmt.Errorf("failed to update recipient's balance: %v", err)
	}

	// Insert a new row in the transfer_histories table
	sqlQuery5 := `INSERT INTO transfer_histories (user_id_sender, user_id_recipient, amount) VALUES (?, ?, ?)`
	stmt, err = tx.Prepare(sqlQuery5)
	if err != nil {
		return "", fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(senderID, recipientID, amount)
	if err != nil {
		return "", fmt.Errorf("failed to insert transfer history: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	outputStr := "\n[SUCCESS] Transfer was successful.\n"
	return outputStr, nil
}
