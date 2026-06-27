package handler

import (
	"httpserver/database"
	"httpserver/database/dbHelper"
	"httpserver/middlewares"
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
	}{}

	if err := utils.ParseBody(r.Body, &body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err, "failed parse request body!")
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

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-api-key")
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
        w.WriteHeader(http.StatusForbidden)
        return
	}

	err := dbHelper.DeleteSessionToken(userCtx.ID, token)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to logout user")
		return
	}
	utils.ResponseJSON(w, http.StatusAccepted, nil)
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
        w.WriteHeader(http.StatusForbidden)
        return
	}
	
	utils.ResponseJSON(w, http.StatusOK, userCtx)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("x-api-key")
	userCtx := middlewares.UserContext(r)
	if userCtx == nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
    
	if err := dbHelper.DeleteUserByID(userCtx.ID); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to delete user!")
		return
	}

	if err := dbHelper.DeleteSessionToken(userCtx.ID, token); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err, "failed to delete user!")
		return 
	}
    
	utils.ResponseJSON(w, http.StatusOK, nil)
}
