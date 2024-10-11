package line

import (
	"fmt"
	"github.com/a-company/yoriai-backend/pkg/config"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log/slog"
)

type LINEBotService struct {
	client *linebot.Client
}

func NewLINEBotService() (*LINEBotService, error) {
	secret := config.Config.LineConfig.ChannelSecret
	tkn := config.Config.LineConfig.AccessToken
	bot, err := linebot.New(secret, tkn)
	if err != nil {
		return nil, fmt.Errorf("failed to create line bot: %w", err)
	}
	return &LINEBotService{
		client: bot,
	}, nil
}

func (l *LINEBotService) ReplyTextMessage(replyToken, message string) {
	if _, err := l.client.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		slog.Error("failed to reply message", err)
	}
}
