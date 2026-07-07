package middlewares

import (
	"context"
	
	"httpserver/database/dbHelper"
	"httpserver/models"
	"httpserver/utils"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type ContextKeys string

const (
	userContext ContextKeys = "__userContext"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		auth := r.Header.Get("Authorization")

		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtToken := strings.TrimPrefix(auth, "Bearer ")
		userInfo, jwtErr := utils.VerifyJwtToken(jwtToken)
		// fmt.Println(userInfo)
		if jwtErr != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		jwtUser := &models.User{
			ID:    userInfo["UserID"].(string),
			Name:  userInfo["UserName"].(string),
			Email: userInfo["UserEmail"].(string),
		}

		apiKey := r.Header.Get("x-api-key")
		sessionUser, err := dbHelper.GetUserBySession(apiKey)

		if err != nil || sessionUser == nil {
			logrus.WithError(err).Errorf("failed to get user with token: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContext, jwtUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserContext(r *http.Request) *models.User {
	if user, ok := r.Context().Value(userContext).(*models.User); ok && user != nil {
		return user
	}
	return nil
}
