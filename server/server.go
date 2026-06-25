package server

import (
	"context"
	"httpserver/handler"
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

		r.Route("/", func(public chi.Router) {
			public.Post("/register", handler.RegisterUser)
			public.Post("/login", handler.LoginUser)
		})

		r.Route("/todo", func(r chi.Router) {

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
