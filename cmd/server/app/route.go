package app

import (
	"github.com/gin-gonic/gin"
	"github.com/ruijzhan/dungeons/pkg/resolve"
)

func (s *Server) addRoutes() {
	addPing(s.Engine)
	addCheckHost(s.Engine, s.dns)
}

func addPing(g *gin.Engine) {
	g.GET("/ping", pong)
}

func addCheckHost(g *gin.Engine, dns resolve.Resolver) {
	endpoint := g.Group("/host")
	endpoint.GET("/:host", checkHost(dns))
}
