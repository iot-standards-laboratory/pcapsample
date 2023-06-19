package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const addr = "localhost:4242"
const message = "Computer Networks Packet Capture"

func newRouter() *gin.Engine {
	r := gin.New()

	r.Any("/*any", func(c *gin.Context) {
		c.String(http.StatusOK, "text/plain", message)
	})
	return r
}
