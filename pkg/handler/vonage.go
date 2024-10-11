package handler

import "github.com/gin-gonic/gin"

type VonageWebhook struct {
}

func NewVonageWebhook() *VonageWebhook {
	return &VonageWebhook{}
}

func (v *VonageWebhook) Handle(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}
