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

func GetUserBySession(sessionToken string) (*models.User, error) {
    
	SQL := `SELECT
	            u.id,
				u.name,
				u.email,
				u.created_at
	        FROM users u
			JOIN user_session us on u.id = us.user_id 
			WHERE u.archived_at IS NULL AND us.session_token = $1`

	var user models.User
	err := database.Todo.Get(&user, SQL, sessionToken)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	
	return &user, nil
}

func DeleteSessionToken(userID, token string) error {
	SQL := `DELETE FROM user_session WHERE user_id = $1 AND session_token = $2`
	_, err := database.Todo.Exec(SQL, userID, token)
	return err
}

func DeleteUserByID(userID string) error {
	SQL := `UPDATE users
	        SET archived_at = NOW()
			WHERE id = $1`
	_, err := database.Todo.Exec(SQL, userID)
	return err		
}