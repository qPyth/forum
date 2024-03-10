package server

import (
	"net/http"

	"forum/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
