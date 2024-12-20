package server

import (
	"YandexPracticum-go-final-TODO/internal/config"
	"log"
	"net/http"
)

type Server struct {
	httpServer *http.Server
	Handler    http.Handler
}

var port = config.Port()

func (s *Server) Run(router http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    port,
		Handler: router,
	}

	log.Printf("Server started on %s", s.httpServer.Addr)

	return s.httpServer.ListenAndServe()
}
