package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ruijzhan/dungeons/pkg/resolve"
)

type ServerOption struct {
	Addr string
}

type Server struct {
	opt ServerOption
	*gin.Engine
	dns resolve.Resolver
}

// NewServer creates a new HTTP server with custom options.
func NewServer(opt ServerOption) *Server {
	srv := &Server{
		opt:    opt,
		Engine: gin.New(),
		dns:    resolve.New(),
	}
	srv.addRoutes()
	return srv
}

func (s *Server) Run() error {
	return s.Engine.Run(s.opt.Addr)
}
