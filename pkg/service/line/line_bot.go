package line

import (
	"fmt"
	"github.com/a-company/yoriai-backend/pkg/config"
	"github.com/line/line-bot-sdk-go/v8/linebot"
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

func (l *LINEBotService) ReplyMessage(replyToken string, flexMessage []linebot.SendingMessage) {
	if _, err := l.client.ReplyMessage(replyToken, flexMessage...).Do(); err != nil {
		slog.Error("failed to reply flex message", err)
	}
}

func (l *LINEBotService) PushTextMessage(to, message string) error {
	if _, err := l.client.PushMessage(to, linebot.NewTextMessage(message)).Do(); err != nil {
		slog.Error("failed to push message", err)
		return err
	}
	return nil
}

func CreateTimeSelectMessage() *linebot.FlexMessage {
	bubble := &linebot.BubbleContainer{
		Type: linebot.FlexContainerTypeBubble,
		Body: &linebot.BoxComponent{
			Type:   linebot.FlexComponentTypeBox,
			Layout: linebot.FlexBoxLayoutTypeVertical,
			Contents: []linebot.FlexComponent{
				&linebot.TextComponent{
					Type:   linebot.FlexComponentTypeText,
					Text:   "電話をかける時間を\n選択してください",
					Weight: linebot.FlexTextWeightTypeBold,
					Size:   linebot.FlexTextSizeTypeLg,
					Margin: linebot.FlexComponentMarginTypeMd,
					Align:  linebot.FlexComponentAlignTypeCenter,
					Wrap:   true,
				},
				createHourRow("09:00", "10:00", "11:00"),
				createHourRow("12:00", "13:00", "14:00"),
				createHourRow("15:00", "16:00", "17:00"),
				createHourRow("18:00", "19:00", "20:00"),
			},
		},
	}

	flexMessage := linebot.NewFlexMessage("TimePicker", bubble)
	return flexMessage
}

func createHourRow(hour1, hour2, hour3 string) *linebot.BoxComponent {
	return &linebot.BoxComponent{
		Type:   linebot.FlexComponentTypeBox,
		Layout: linebot.FlexBoxLayoutTypeHorizontal,
		Contents: []linebot.FlexComponent{
			&linebot.ButtonComponent{
				Type: linebot.FlexComponentTypeButton,
				Action: &linebot.PostbackAction{
					Label:       hour1,
					Data:        "call_time_picker_" + hour1,
					DisplayText: hour1,
				},
				Flex:   linebot.IntPtr(1),
				Margin: linebot.FlexComponentMarginTypeSm,
			},
			&linebot.ButtonComponent{
				Type: linebot.FlexComponentTypeButton,
				Action: &linebot.PostbackAction{
					Label:       hour2,
					Data:        "call_time_picker_" + hour2,
					DisplayText: hour2,
				},
				Flex:   linebot.IntPtr(1),
				Margin: linebot.FlexComponentMarginTypeSm,
			},
			&linebot.ButtonComponent{
				Type: linebot.FlexComponentTypeButton,
				Action: &linebot.PostbackAction{
					Label:       hour3,
					Data:        "call_time_picker_" + hour3,
					DisplayText: hour3,
				},
				Flex:   linebot.IntPtr(1),
				Margin: linebot.FlexComponentMarginTypeSm,
			},
		},
		Margin: linebot.FlexComponentMarginTypeSm,
	}
}
