package dbHelper

import (
	"fmt"
	"httpserver/database"
	"httpserver/models"
	"strings"
	"time"
)

func CreateTodo(userID string, name string, description string, pendingAt time.Time) (string, error) {

	SQL := `INSERT INTO TODO (user_id, name, description, pending_at) VALUES ($1, $2, $3, $4) RETURNING id`
	var todoId string
	if err := database.Todo.QueryRowx(SQL, userID, name, description, pendingAt).Scan(&todoId); err != nil {
		return "", err
	}

	return todoId, nil

}

func GetTodoByID(todoID string, userID string) (models.Todo, error) {

	SQL := `SELECT id, user_id, name, description, pending_at, created_at
	        FROM TODO 
			WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	var todo models.Todo
	if err := database.Todo.Get(&todo, SQL, todoID, userID); err != nil {
		return todo, err
	}

	return todo, nil
}

func GetTodosByUserID(userID string) ([]models.Todo, error) {

	SQL := `SELECT id, user_id, name, description, pending_at, created_at
	FROM TODO
	WHERE user_id = $1 AND archived_at IS NULL`
	var todos []models.Todo
	if err := database.Todo.Select(&todos, SQL, userID); err != nil {
		return todos, err
	}
	return todos, nil
}

func DeleteTodoByID(todoID string, userID string) error {

	SQL := `UPDATE TODO
	        SET archived_at = NOW()
			WHERE id = $1 AND user_id = $2`
	
	_, err := database.Todo.Exec(SQL, todoID, userID)		

	return err
}

func UpdateTodoByID(todoID string, userID string, agr []string, val []any) error {

	SQL := `UPDATE TODO 
	        SET %s 
			WHERE id = ? AND user_id = ?`	

	SQL = fmt.Sprintf(SQL, strings.Join(agr, ", "))		

	SQL = database.ReplaceSQL(SQL, "?")
    
	val = append(val, todoID, userID)

	_, err := database.Todo.Exec(SQL, val...)

	return err
}