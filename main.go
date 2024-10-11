package main

import (
	"fmt"
	"github.com/a-company/yoriai-backend/pkg/config"
	"github.com/a-company/yoriai-backend/pkg/handler"
	"github.com/a-company/yoriai-backend/pkg/service/line"
	"github.com/a-company/yoriai-backend/pkg/util/firestore"
	"github.com/gin-gonic/gin"
	"log/slog"
	"strconv"
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

	fs := firestore.New()
	lineWHandler := handler.NewLINEWebhookHandler(lineBotSvc, fs)
	e.Any("/line/webhook", lineWHandler.Handle)

	vonageWHService := handler.NewVonageWebhook()
	e.Any("/vonage/webhook", vonageWHService.Handle)
	port := 8080
	if config.Config.General.Port != "" {
		port, err = strconv.Atoi(config.Config.General.Port)
		if err != nil {
			slog.Error("failed to load port, invalid format", err)
			return
		}
	}

	invokeHandler := handler.NewInvokeHandler(fs)
	e.POST("/invoke", invokeHandler.Handle)

	e.Run(fmt.Sprintf(":%d", port))
}
