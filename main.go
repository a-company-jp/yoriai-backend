package main

import (
	"github.com/a-company/yoriai-backend/pkg/handler"
	"github.com/a-company/yoriai-backend/pkg/service/line"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func main() {
	e := gin.Default()

	e.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	lineBotSvc, err := line.NewLINEBotService()
	if err != nil {
		slog.Error("failed to initialize line bot service", err)
		return
	}

	lineWHandler := handler.NewLINEWebhookHandler(lineBotSvc)
	e.Any("/line/webhook", lineWHandler.Handle)
	e.Run(":8080")
}
