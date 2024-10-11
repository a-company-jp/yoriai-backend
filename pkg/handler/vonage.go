package handler

import "github.com/gin-gonic/gin"

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
}

func NewVonageWebhook() *VonageWebhook {
	return &VonageWebhook{}
}

func (v *VonageWebhook) Handle(c *gin.Context) {
	var req VonageWebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// TODO: forward the request to LINE
	c.JSON(200, gin.H{
		"message": "success",
	})
}
