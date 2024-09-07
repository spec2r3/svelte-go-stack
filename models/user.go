package models

import (
	"fmt"
	"gooooo/db"
	_ "log"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
	Alias    string `binding:"required"`
	APIKey   string
	Admin    bool
}

func (u *User) Save() error {
	query := `INSERT INTO users (email, password, api_key, alias) VALUES ( ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Email, u.Password, u.APIKey, u.Alias)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = id

	return nil
}

func DeleteUserById(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	// Prepare the statement
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute the statement
	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	// Check the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", id)
	}

	return nil
}
