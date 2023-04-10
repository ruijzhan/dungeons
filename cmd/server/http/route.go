package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ruijzhan/dungeons/pkg/resolve"
)

/*
AddPing registers a new route with the given Gin engine that responds to GET
requests to "/ping" with the pong function.

Parameters:
- g: A pointer to a Gin engine.

Return:
None.
*/
func AddPing(g *gin.Engine) {
	g.GET("/ping", pong)
}

// AddCheckHost registers a GET endpoint on the given gin engine that checks the
// availability of a given host using the provided DNS resolver.
func AddCheckHost(engine *gin.Engine, resolver resolve.Resolver) {
	// Define a new endpoint group for /host.
	hostGroup := engine.Group("/host")

	// Register the checkHost handler on the /host/:host endpoint.
	hostGroup.GET("/:host", checkHost(resolver))
}
