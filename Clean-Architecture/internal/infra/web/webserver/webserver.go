package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

// AddHandler cadastra as rotas
func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() error {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	log.Printf("Server listening on :%s", s.WebServerPort)
	return http.ListenAndServe(s.WebServerPort, s.Router)
}
