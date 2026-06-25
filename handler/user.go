package handler

import (
	"httpserver/database"
	"httpserver/database/dbHelper"
	"httpserver/models"
	"httpserver/utils"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Name     string      `json:"name"`
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Role     models.Role `json:"role"`
	}{}

	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err, "failed parse request body!")
		return
	}

	if !body.Role.IsValid() {
		utils.ResponseError(w, http.StatusBadRequest, nil, "invalid role type provided!")
		return
	}

	if len(body.Password) < 8 {
		utils.ResponseError(w, http.StatusBadRequest, nil, "password must be 8 chars long!")
		return
	}

	exists, existsErr := dbHelper.IsUserExists(body.Email)

	if existsErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, existsErr, "failed to check user existence!")
		return
	}

	if exists {
		utils.ResponseError(w, http.StatusBadRequest, nil, "user already exists!")
		return
	}

	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, hasErr, "failed to secure password!")
		return
	}

	sessionToken := utils.HashString(body.Email + time.Now().String())
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
		if saveErr != nil {
			return saveErr
		}

		roleErr := dbHelper.CreateUserRole(tx, userID, body.Role)
		if roleErr != nil {
			return roleErr
		}

		sessionErr := dbHelper.CreateUserSession(tx, userID, sessionToken)
		if sessionErr != nil {
			return sessionErr
		}

		return nil
	})

	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})

}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	body := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, parseErr, "failed to parse request body!")
		return
	}

	userID, userErr := dbHelper.GetUserIDByPassword(body.Email, body.Password)
	if userErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, userErr, "failed to find user!")
		return
	}

	sessionToken := utils.HashString(body.Email + time.Now().String())
	sessionErr := dbHelper.CreateUserSession(database.Todo, userID, sessionToken)
	if sessionErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, sessionErr, "failed to create user session")
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
	}{
		Token: sessionToken,
	})
}
