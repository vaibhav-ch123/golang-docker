package dbHelper

import (
	"httpserver/database"
	"time"
)

func CreateTodo(name string, description string, pendingAt time.Time) (string, error) {

	SQL := `INSERT INTO TODO (name, description, pending_at) VALUES ($1, $2, $3) RETURNING id`
	var userId string
	if err := database.Todo.QueryRowx(SQL, name, description, pendingAt).Scan(&userId); err != nil {
		return "", err
	}

	return userId, nil

}

// func GetTodo(id string){

// 	SQL := `SELECT * FROM TODO T WHERE T.id = $1`

// 	var todo models.Todo

// }
