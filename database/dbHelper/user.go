package dbHelper

import (
	"database/sql"
	"httpserver/database"
	"httpserver/models"
	"httpserver/utils"
	"strings"

	"github.com/jmoiron/sqlx"
)

func IsUserExists(email string) (bool, error) {

	SQL := `SELECT id FROM USERS WHERE email = TRIM(LOWER($1)) AND archived_at IS NULL`
	var id string
	err := database.Todo.Get(&id, SQL, strings.ToLower(email))

	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func CreateUser(db sqlx.Ext, name, email, password string) (string, error) {

	SQL := `INSERT INTO users(name, email, password) VALUES ($1, TRIM(LOWER($2)), $3) RETURNING id`
	var userID string
	if err := db.QueryRowx(SQL, name, email, password).Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func CreateUserRole(db sqlx.Ext, userID string, role models.Role) error {

	SQL := `INSERT INTO user_roles(user_id, role) VALUES ($1, $2)`
	_, err := db.Exec(SQL, userID, role)
	return err
}

func CreateUserSession(db sqlx.Ext, userID, sessionToken string) error {

	SQL := `INSERT INTO user_session(user_id, session_token) VALUES ($1, $2)`
	_, err := db.Exec(SQL, userID, sessionToken)
	return err
}

func GetUserIDByPassword(email, password string) (string, error) {

	SQL := `SELECT 
	            id,
				password
			FROM 
			    users
			WHERE
			    archived_at IS NULL 
				AND email = TRIM(LOWER($1))`
    var userID, passwordHash string
	err := database.Todo.QueryRowx(SQL, email).Scan(&userID, &passwordHash)
    if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if err == sql.ErrNoRows {
		return "", nil
	}

	if passwordErr := utils.CheckPassword(password, passwordHash); passwordErr != nil {
		return "", passwordErr
	}

	return userID, nil
}
