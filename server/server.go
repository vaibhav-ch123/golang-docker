package server

import (
	"context"
	"fmt"
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
		r.Get("/server", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "hello server!")
		})
		r.Route("/", func(r chi.Router) {
			r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "hello user!")
			})
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
