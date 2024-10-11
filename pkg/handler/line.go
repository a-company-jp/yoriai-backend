package handler

import (
	"github.com/a-company/yoriai-backend/pkg/config"
	"github.com/a-company/yoriai-backend/pkg/service/line"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

type LINEWebhookHandler struct {
	lineBotSvc *line.LINEBotService
}

func NewLINEWebhookHandler(lineBotSvc *line.LINEBotService) *LINEWebhookHandler {
	return &LINEWebhookHandler{
		lineBotSvc: lineBotSvc,
	}
}

func (l *LINEWebhookHandler) Handle(c *gin.Context) {
	req, err := webhook.ParseRequest(config.Config.LineConfig.ChannelSecret, c.Request)
	if err != nil {
		return
	}

	for _, event := range req.Events {
		switch event.GetType() {
		case "join":
			// join: When your LINE Official Account joins a group chat or multi-person chat. You can reply to this event.
		case "accountLink":
			// accountLink: When a user has linked their LINE account with a provider's service. You can reply to this event.
		case "activated":
			// activated: When a user has linked their LINE account with a provider's service. You can reply to this event.
		case "beacon":
			// beacon: When a user enters the range of a LINE Beacon. You can reply to this event.
		case "botResumed":
			// botResumed: When a LINE Official Account that was suspended is resumed. You can reply to this event.
		case "botSuspended":
			// botSuspended: When a LINE Official Account is suspended. You can reply to this event.
		case "deactivated":
			// deactivated: When a user has unlinked their LINE account from a provider's service. You can reply to this event.
		case "delivery":
			// delivery: When a message is successfully delivered to a user. You can't reply to this event.
		case "follow":
			// follow: When a user adds your LINE Official Account as a friend, or unblocks your LINE Official Account. You can reply to this event.
			l.HandleFollowEvent(event.(webhook.FollowEvent))
		case "leave":
			// leave: When a user deletes your LINE Official Account or your LINE Official Account leaves, from a group chat or multi-person chat.
		case "memberJoined":
			// memberJoined: When a user joins a group chat or multi-person chat that your LINE Official Account is a member of. You can reply to this event.
		case "memberLeft":
			// memberLeft: When a user leaves a group chat or multi-person chat that your LINE Official Account is a member of.
		case "message":
			// message: When a user sends a message. You can reply to this event.
		case "module":
			// module: When a user interacts with a module. You can reply to this event.
		case "postback":
			// postback: When a user triggers a postback action. You can reply to this event.
		case "things":
			// things: When a user interacts with a LINE Things-compatible device. You can reply to this event.
		case "unfollow":
			// unfollow: When a user blocks your LINE Official Account.
		case "unsend":
			// unsend: When a user unsends a message. For more information on handling this event, see Processing on receipt of unsend event.
		case "videoPlayComplete":
			// videoPlayComplete: When a user finishes watching a video message that has a trackingId specified sent from the LINE Official Account. You can reply to this event.
		}
	}
}

func (l *LINEWebhookHandler) HandleFollowEvent(event webhook.FollowEvent) {
	l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "Thank you for adding me as a friend!")
}

func (l *LINEWebhookHandler) HandleMessageEvent(event webhook.MessageEvent) {
}

func (l *LINEWebhookHandler) HandleLeaveEvent(event webhook.UnfollowEvent) {
}
