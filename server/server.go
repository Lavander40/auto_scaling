package server

import (
	"io"
	"net/http"
	"github.com/gorilla/mux"
)

type Server struct {
	router  *mux.Router
}

func New() *Server {
	return &Server{
		router: mux.NewRouter(),
	}
}

func (s *Server) Start() error {
	s.configureRouter()
	return http.ListenAndServe(":4040", s.router)
}

func (s *Server) configureRouter() {
	s.router.HandleFunc("/", s.handleIndex())
}

func (s *Server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "server response")
	}
}
