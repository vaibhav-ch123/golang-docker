package handler

import (
	"database/sql"
	"httpserver/database/dbHelper"
	"httpserver/middlewares"
	"httpserver/utils"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func AddTodo(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PendingAt   time.Time `json:"pendingAt"`
	}{}

	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
        w.WriteHeader(http.StatusForbidden)
        return
	}

	if ParseErr := utils.ParseBody(r.Body, &body); ParseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, ParseErr, "failed to parse request body")
		return
	}

	if len(strings.TrimSpace(body.Name)) < 6 {
		utils.ResponseError(w, http.StatusBadRequest, nil, "todo name length must greater then 6")
		return
	}

	if len(strings.TrimSpace(body.Description)) < 10 {
		utils.ResponseError(w, http.StatusBadRequest, nil, "todo description length must by greater 10")
		return
	}

	todoId, err := dbHelper.CreateTodo(userCtx.ID, body.Name, body.Description, body.PendingAt)

	if err != nil {
		logrus.Errorf("failed to create todo in DB: %+v", err)
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to create todo in DB")
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, struct {
		TodoId string `json:"todoId"`
	}{
		TodoId: todoId,
	})
}

func GetTodo(w http.ResponseWriter, r *http.Request) {

	ID := chi.URLParam(r, "id")

	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
        w.WriteHeader(http.StatusForbidden)
        return
	}

	todo, err := dbHelper.GetTodoByID(ID, userCtx.ID)

	if err != nil && err == sql.ErrNoRows {
		utils.ResponseError(w, http.StatusNotFound, err, "todo not found!")
		return
	}

	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to get todo!")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, todo)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {

	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
        w.WriteHeader(http.StatusForbidden)
        return
	}

	todos, err := dbHelper.GetTodosByUserID(userCtx.ID)

	if err != nil && err == sql.ErrNoRows {
		utils.ResponseError(w, http.StatusNotFound, err, "todos not found!")
		return
	}

	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to get todos!")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, todos)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {

	todoID := chi.URLParam(r, "id")
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
        w.WriteHeader(http.StatusForbidden)
        return
	}

	if err := dbHelper.DeleteTodoByID(todoID, userCtx.ID); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to delete todo!")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, nil)

}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {

    body := struct{
		Name *string              `json:"name"`
		Description *string       `json:"description"`
		PendingAt *time.Time      `json:"pendingAt"`
	}{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, parseErr, "failed to parse body!")
		return 
	}

	if body.Name != nil && len(strings.TrimSpace(*body.Name)) < 6 {
		utils.ResponseError(w, http.StatusBadRequest, nil, "todo name length must greater then 6")
		return
	}

	if body.Description != nil && len(strings.TrimSpace(*body.Description)) < 10 {
		utils.ResponseError(w, http.StatusBadRequest, nil, "todo description length must by greater 10")
		return
	}

	todoID := chi.URLParam(r, "id")
	userCtx := middlewares.UserContext(r)

	if userCtx == nil {
		w.WriteHeader(http.StatusForbidden)
        return
	}
    agr := make([]string, 0)
	val := make([]any, 0)
	
	if body.Name != nil {
		agr = append(agr, "name = ?")
		val = append(val, *body.Name)
	}

	if body.Description != nil {
		agr = append(agr, "description = ?")
		val = append(val, *body.Description)
	}

	if body.PendingAt != nil {
        agr = append(agr, "pending_at = ?")
		val = append(val, *body.PendingAt)
	}
	
	if err := dbHelper.UpdateTodoByID(todoID, userCtx.ID, agr, val); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to update todo!")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, nil)
}
