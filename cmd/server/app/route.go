package app

import "github.com/gin-gonic/gin"

const (
	statusCodeOK = 200
)

func (s *Server) addRoutes() {
	ping(s.Engine)
}

func ping(g *gin.Engine) {

	type pingResponse struct {
		Message string `json:"message"`
	}

	g.GET("/ping", func(c *gin.Context) {
		res := pingResponse{Message: "pong"}
		c.JSON(statusCodeOK, res)
	})
}
