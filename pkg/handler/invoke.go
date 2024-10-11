package handler

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type CallInvoke struct {
	fs *firestore.Client
}

func NewInvokeHandler(
	fs *firestore.Client,
) *CallInvoke {
	return &CallInvoke{
		fs: fs,
	}
}

func (h *CallInvoke) Handle(c *gin.Context) {
	// get time
	timeVal := fmt.Sprintf("%2d:%2d", time.Now().Hour(), 0)
	res := h.fs.Collection("users").Where("call_time", "==", timeVal).Documents(c)
	if res == nil {
		c.JSON(400, gin.H{"error": "no data"})
		return
	}
	for {
		doc, err := res.Next()
		if err != nil {
			break
		}
		userdata := doc.Data()
		slog.Info("invoke call on user", slog.Any("data", userdata))

		// vonage callを発火
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}
