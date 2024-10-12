package handler

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"github.com/a-company/yoriai-backend/pkg/model"
	"github.com/a-company/yoriai-backend/pkg/service/line"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type VonageWebhookRequest struct {
	AgentID        string `json:"agent_id"`
	SessionID      string `json:"session_id"`
	ConversationID string `json:"conversation_id"`
	Feeling        string `json:"feeling"`
	PhoneNumber    string `json:"phone_number"`
	Message        string `json:"message"`
	TodayActivity  string `json:"today_activity"`
}

type VonageWebhook struct {
	line *line.LINEBotService
	fs   *firestore.Client
}

func NewVonageWebhook(
	svc *line.LINEBotService,
	fs *firestore.Client,
) *VonageWebhook {
	return &VonageWebhook{
		line: svc,
		fs:   fs,
	}
}

func (v *VonageWebhook) Handle(c *gin.Context) {
	var req VonageWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userdata := model.User{}
	doc, err := v.fs.Collection("users").
		Where("phone_number", "==", req.PhoneNumber).
		Documents(c).Next()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	doc.DataTo(&userdata)

	//気分: Feeling
	//今日の通話のサマリー: TodayActivity
	//伝言: Message

	notifyText := "本日の通話が終了しました\n\n"
	notifyText += fmt.Sprintf("気分: 元気\n")
	notifyText += fmt.Sprintf("今日の通話のサマリー: 今ハッカソンに出ています。\n")
	notifyText += fmt.Sprintf("伝言: 昨日、友達の太郎に久しぶりに会いました。\n")

	slog.Info("notifyText", slog.String("notifyText", notifyText))
	if err := v.line.PushTextMessage(userdata.LINEID, notifyText); err != nil {
		slog.Error("failed to push message", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
	})
}
