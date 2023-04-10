package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ruijzhan/dungeons/pkg/resolve"
)

const (
	defaultTimeout = 5 * time.Second
)

func pong(c *gin.Context) {
	res := struct{ Message string }{Message: "pong"}
	c.JSON(http.StatusOK, res)
}

func checkHost(dns resolve.Resolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		host := c.Param("host")
		ips, err := dns.Resolve(ctx, host)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, resolve.Response{
			Host: host,
			IPs:  ips,
		})
	}
}
