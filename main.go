package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "lx",
		})
	})
	r.Run()
}
