package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zlingqu/nvidia-gpu-mem-monitor/handlers"
	"net/http"
)

func main() {

	r := gin.Default()

	r.GET("/metrics", func(c *gin.Context) {
		r := handlers.Metrics()
		c.String(http.StatusOK, r)
	})
	r.Run(":80")

}
