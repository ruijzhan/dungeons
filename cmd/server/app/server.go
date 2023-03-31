package app

import (
	"github.com/gin-gonic/gin"
)

type ServerOption struct {
	Addr string
}

type Server struct {
	opt ServerOption
	*gin.Engine
}

// NewServer creates a new HTTP server with custom options.
func NewServer(opt ServerOption) *Server {
	srv := &Server{
		opt:    opt,
		Engine: gin.New(),
	}
	srv.addRoutes()
	return srv
}

func (s *Server) Run() error {
	return s.Engine.Run(s.opt.Addr)
}
