package dbHelper

import (
	"fmt"
	"httpserver/database"
	"httpserver/models"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
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

	SQL := `SELECT id, user_id, name, description, is_completed, pending_at, created_at
	        FROM TODO 
			WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	var todo models.Todo
	if err := database.Todo.Get(&todo, SQL, todoID, userID); err != nil {
		return todo, err
	}

	return todo, nil
}

func GetTodosByUserID(userID string) ([]models.Todo, error) {

	SQL := `SELECT id, user_id, name, description, is_completed, pending_at, created_at
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
			WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	_, err := database.Todo.Exec(SQL, todoID, userID)

	return err
}

func DeleteTodosByUserID(db sqlx.Ext, userID string) error {

	SQL := `UPDATE TODO
	        SET archived_at = NOW()
			WHERE user_id = $1 AND archived_at IS NULL`

	_, err := db.Exec(SQL, userID)

	return err
}

func UpdateTodoByID(todoID string, userID string, agr []string, val []any) error {

	SQL := `UPDATE TODO 
	        SET %s 
			WHERE id = ? AND user_id = ? AND archived_at IS NULL`

	SQL = fmt.Sprintf(SQL, strings.Join(agr, ", "))

	SQL = database.ReplaceSQL(SQL, "?")

	val = append(val, todoID, userID)

	_, err := database.Todo.Exec(SQL, val...)

	return err
}

func IsCompletedTodoByID(todoID string, userID string) (bool, error) {

	SQL := `SELECT is_completed from todo WHERE id = $1 AND user_id = $2 AND archived_at IS NULL`

	var iscompleted bool
	if err := database.Todo.Get(&iscompleted, SQL, todoID, userID); err != nil {
		return false, err
	}

	return iscompleted, nil
}

func MarkCompletedTodoByID(todoID string, userID string, is_completed bool) error {

	SQL := `UPDATE todo 
	        SET is_completed = $1 
			WHERE id = $2 AND user_id = $3 AND archived_at IS NULL`
	logrus.Println(is_completed)
	_, err := database.Todo.Exec(SQL, !is_completed, todoID, userID)

	return err
}
