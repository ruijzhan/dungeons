package app

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruijzhan/dungeons/pkg/resolve"
	"github.com/spf13/cast"
)

func (s *Server) addRoutes() {
	addPing(s.Engine)
	addCheckHost(s.Engine, s.dns)
}

func addPing(g *gin.Engine) {

	type pingResponse struct {
		Message string `json:"message"`
	}

	g.GET("/ping", func(c *gin.Context) {
		res := pingResponse{Message: "pong"}
		c.JSON(http.StatusOK, res)
	})
}

func addCheckHost(g *gin.Engine, dns resolve.Resolver) {

	endpoint := g.Group("/host")
	endpoint.GET("/:host", func(c *gin.Context) {
		host := c.Param("host")
		ctx, cancel := context.WithTimeout(context.Background(), cast.ToDuration("5s"))
		defer cancel()
		ips, err := dns.Resolve(ctx, host)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"host": host,
			"ips":  ips,
		})
	})
}
