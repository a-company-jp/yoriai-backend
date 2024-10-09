package main

import "github.com/gin-gonic/gin"

func main() {
	e := gin.Default()

	e.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})
	e.Run(":8080")
}
