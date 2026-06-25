package handler

import (
	"httpserver/database/dbHelper"
	"httpserver/utils"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func AddTodo(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		PendingAt   time.Time `json:"pendingAt,omitempty"`
	}{}

	if ParseErr := utils.ParseBody(r.Body, &body); ParseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, ParseErr, "failed to parse request body")
		return
	}

	if len(strings.TrimSpace(body.Name)) < 6 {
		utils.ResponseError(w, http.StatusBadRequest, nil, "todo name length must by greater 6")
		return
	}

	if len(strings.TrimSpace(body.Description)) < 10 {
		utils.ResponseError(w, http.StatusBadRequest, nil, "todo description length must by greater 10")
		return
	}

	UserId, err := dbHelper.CreateTodo(body.Name, body.Description, body.PendingAt)

	if err != nil {
		logrus.Errorf("failed to create todo in DB: %+v", err)
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to create todo in DB")
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, struct {
		UserId string `json:"userId"`
	}{
		UserId: UserId,
	})
}
