package models

import (
	"fmt"
	"gooooo/db"
	"gooooo/utils"
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

type UserResponse struct {
	ID    int64
	Email string
	Alias string
	Admin bool
}

type UsersResponse struct {
	Users       []UserResponse `json:"users"`
	TotalCount  int            `json:"totalCount"`
	TotalPages  int            `json:"totalPages"`
	CurrentPage int            `json:"currentPage"`
	PageSize    int            `json:"pageSize"`
}

type SignIn struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := `INSERT INTO users (email, password, api_key, alias) VALUES ( ?, ?, ?, ?)`

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	hashedKey, err := utils.HashKey(u.APIKey)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, hashedPassword, hashedKey, u.Alias)
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

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with ID %d", id)
	}

	return nil
}

func GetAllUsers(page, pageSize int) ([]User, int, error) {
	offset := (page - 1) * pageSize

	var userCount int
	countQuery := `SELECT COUNT(*) FROM users`
	err := db.DB.QueryRow(countQuery).Scan(&userCount)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, email, alias, admin FROM users LIMIT ? OFFSET ?`
	rows, err := db.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Alias, &user.Admin); err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, userCount, nil
}
