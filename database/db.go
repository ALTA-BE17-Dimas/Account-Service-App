package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func DBConnect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, fmt.Errorf("[FAIL] error connecting to database: %w", err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	fmt.Print("\n[OK] Database connected successfully.\n\n")

	return db, nil
}
