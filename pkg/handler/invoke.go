package handler

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/a-company/yoriai-backend/pkg/model"
	"github.com/a-company/yoriai-backend/pkg/service/vonage"
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
	// specify tokyo and get hour
	hour := time.Now().In(time.FixedZone("Asia/Tokyo", 9*60*60)).Hour()
	timeVal := fmt.Sprintf("%02d:%02d", hour, 0)
	slog.Info("invoke call", slog.String("time", timeVal))
	res := h.fs.Collection("users").Where("call_time", "==", timeVal).Documents(c)
	if res == nil {
		c.JSON(400, gin.H{"error": "no data"})
		return
	}

	v := vonage.NewVonage()
	for {
		doc, err := res.Next()
		if err != nil {
			break
		}
		userdata := model.User{}
		doc.DataTo(&userdata)
		slog.Info("invoke call on user", slog.Any("data", userdata))

		// vonage callを発火
		v.CallPhoneAPI(
			vonage.PhoneAPIInput{
				PhoneNumber:   string(userdata.Phone),
				ReceiverName:  userdata.Nickname,
				CallerName:    "Yoriai",
				RemindMessage: "今日の通話の時間です",
			})
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}
