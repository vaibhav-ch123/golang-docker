package server

import (
	"context"
	"httpserver/handler"
	"httpserver/middlewares"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	writeTimeout      = 5 * time.Minute
	readHeaderTimeout = 2 * time.Minute
)

func SetupUpRoutes() *Server {

	router := chi.NewRouter()

	router.Route("/v1", func(r chi.Router) {
		r.Use(middlewares.CommonMiddlewares()...)
		r.Route("/", func(public chi.Router) {
			public.Post("/register", handler.RegisterUser)
			public.Post("/login", handler.LoginUser)
		})

		r.Route("/user", func(user chi.Router) {
			user.Use(middlewares.AuthMiddleware)
			user.Delete("/logout", handler.LogoutUser)
			user.Get("/", handler.GetUserInfo)
			user.Delete("/", handler.DeleteUser)
		})

		r.Route("/todo", func(todo chi.Router) {
			todo.Use(middlewares.AuthMiddleware)
			todo.Post("/", handler.AddTodo)
			todo.Get("/{id}", handler.GetTodo)
			todo.Get("/", handler.GetTodos)
			todo.Delete("/{id}", handler.DeleteTodo)
			todo.Patch("/{id}", handler.UpdateTodo)
			todo.Patch("/{id}/is-completed", handler.IsTodoCompleted)
		})

	})

	return &Server{
		Router: router,
	}
}

func (srv *Server) Run(Port string) error {
	srv.server = &http.Server{
		Addr:              Port,
		Handler:           srv.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}

	return srv.server.ListenAndServe()
}

func (srv *Server) ShutDown(timeout time.Duration) error {
	ctx, close := context.WithTimeout(context.Background(), timeout)
	defer close()
	return srv.server.Shutdown(ctx)
}
