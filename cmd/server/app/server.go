package app

import (
	"github.com/gin-gonic/gin"
)

type Option struct {
	Addr string
}

func (o *Option) newServer() *Server {
	s := &Server{
		addr:   o.Addr,
		Engine: gin.New(),
	}
	return s
}

type Server struct {
	addr string
	*gin.Engine
}

func (s *Server) run() error {
	s.route()
	return s.Run(s.addr)
}

func (s *Server) route() {
	s.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
