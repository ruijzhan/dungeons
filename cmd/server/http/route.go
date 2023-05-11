package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ruijzhan/dungeons/pkg/resolve"
)

func InstallPing(g *gin.Engine) {
	g.GET("/ping", pong)
}

// InstallCheckHost registers a GET endpoint on the given gin engine that checks the
// availability of a given host using the provided DNS resolver.
func InstallCheckHost(engine *gin.Engine, resolver resolve.Resolver) {
	// Define a new endpoint group for /host.
	hostGroup := engine.Group("/host")

	// Register the checkHost handler on the /host/:host endpoint.
	hostGroup.GET("/:host", checkHost(resolver))
}
