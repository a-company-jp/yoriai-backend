package handler

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/a-company/yoriai-backend/pkg/config"
	"github.com/a-company/yoriai-backend/pkg/model"
	"github.com/a-company/yoriai-backend/pkg/service/line"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
	"log/slog"
	"regexp"
	"time"
)

type LINEWebhookHandler struct {
	lineBotSvc *line.LINEBotService
	fs         *firestore.Client
}

func NewLINEWebhookHandler(lineBotSvc *line.LINEBotService,
	fs *firestore.Client,
) *LINEWebhookHandler {
	return &LINEWebhookHandler{
		lineBotSvc: lineBotSvc,
		fs:         fs,
	}
}

func (l *LINEWebhookHandler) Handle(c *gin.Context) {
	req, err := webhook.ParseRequest(config.Config.LineConfig.ChannelSecret, c.Request)
	if err != nil {
		return
	}

	for _, event := range req.Events {
		slog.Info("new line webhook event received", slog.String("type", event.GetType()))
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
			l.HandleMessageEvent(event.(webhook.MessageEvent))
		case "module":
			// module: When a user interacts with a module. You can reply to this event.
		case "postback":
			// postback: When a user triggers a postback action. You can reply to this event.
			l.HandlePostbackEvent(event.(webhook.PostbackEvent))
		case "things":
			// things: When a user interacts with a LINE Things-compatible device. You can reply to this event.
		case "unfollow":
			// unfollow: When a user blocks your LINE Official Account.
			l.HandleUnfollowEvent(event.(webhook.UnfollowEvent))
		case "unsend":
			// unsend: When a user unsends a message. For more information on handling this event, see Processing on receipt of unsend event.
		case "videoPlayComplete":
			// videoPlayComplete: When a user finishes watching a video message that has a trackingId specified sent from the LINE Official Account. You can reply to this event.
		}
	}
}

func (l *LINEWebhookHandler) HandleFollowEvent(event webhook.FollowEvent) {
	userID := event.Source.(webhook.UserSource).UserId
	// time from timestamp int64
	newUser := model.User{
		LINEID:  userID,
		AddDate: time.Unix(event.Timestamp/1000, 0),
	}
	ctx := context.Background()
	_, err := l.fs.Collection("users").Doc(userID).Create(ctx, newUser)
	if err != nil {
		slog.Error("failed to create user", err)
		l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "Failed to add you as a friend. Please try again later.")
		return
	}
	l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "友達登録ありがとうございます！\n\nまずはあなたのニックネームを入力してください")
}

func (l *LINEWebhookHandler) HandleUnfollowEvent(event webhook.UnfollowEvent) {
	userID := event.Source.(webhook.UserSource).UserId
	ctx := context.Background()
	_, err := l.fs.Collection("users").Doc(userID).Delete(ctx)
	if err != nil {
		slog.Error("failed to delete user", err)
		return
	}
}

func (l *LINEWebhookHandler) HandleMessageEvent(event webhook.MessageEvent) {
	if event.Message.GetType() != "text" {
		l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "テキストで回答してください🙏")
	}
	txtMsg := event.Message.(webhook.TextMessageContent)

	userdata := model.User{}
	ctx := context.Background()
	ref := l.fs.Collection("users").Doc(event.Source.(webhook.UserSource).UserId)
	doc, err := ref.Get(ctx)
	if err != nil {
		slog.Error("failed to get user", err)
		return
	}
	doc.DataTo(&userdata)

	if userdata.Target.Nickname == "" {
		userdata.Target.Nickname = txtMsg.Text
		_, err := ref.Set(ctx, userdata)
		if err != nil {
			slog.Error("failed to set user", err)
			return
		}
		l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "次に相手のニックネームを入力してください")
		return
	}
	if userdata.Target.RecipientNickname == "" {
		userdata.Target.RecipientNickname = txtMsg.Text
		_, err := ref.Set(ctx, userdata)
		if err != nil {
			slog.Error("failed to set user", err)
			return
		}
		l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "次に電話番号を入力してください (例: 09012345678)")
		return
	}
	if userdata.Target.Phone == "" {
		userdata.Target.Phone = model.PhoneNumber(txtMsg.Text)
		// 0から始まる10桁か11桁の数字
		phoneRgx := `^0\d{9,10}$`
		if !regexp.MustCompile(phoneRgx).MatchString(string(userdata.Target.Phone)) {
			l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "電話番号の形式が正しくありません。もう一度入力してください。(例: 09012345678)")
			return
		}
		// 先頭の0を81に変換
		userdata.Target.Phone = model.PhoneNumber("81" + string(userdata.Target.Phone)[1:])
		_, err := ref.Set(ctx, userdata)
		if err != nil {
			slog.Error("failed to set user", err)
			return
		}
		l.lineBotSvc.ReplyMessage(event.ReplyToken, []linebot.SendingMessage{line.CreateTimeSelectMessage()})
		return
	}
	if userdata.Target.CallTime == "" {
		l.lineBotSvc.ReplyMessage(event.ReplyToken,
			[]linebot.SendingMessage{
				linebot.NewTextMessage("時間の指定方法が正しくありません。もう一度選択してください"),
				line.CreateTimeSelectMessage()})
		return
	}
	if userdata.Target.RemindMessage == "" && !userdata.Target.Confirm {
		if txtMsg.Text == "なし" {
			userdata.Target.Confirm = true
		} else {
			userdata.Target.RemindMessage = txtMsg.Text
		}
		_, err := ref.Set(ctx, userdata)
		if err != nil {
			slog.Error("failed to set user", err)
			return
		}
		l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "登録完了しました！")
		return
	}
	l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "登録が完了しています！")
}

func (l *LINEWebhookHandler) HandlePostbackEvent(event webhook.PostbackEvent) {
	userID := event.Source.(webhook.UserSource).UserId
	ctx := context.Background()
	userdata := model.User{}
	ref := l.fs.Collection("users").Doc(userID)
	doc, err := ref.Get(ctx)
	if err != nil {
		slog.Error("failed to get user", err)
		return
	}
	if err := doc.DataTo(&userdata); err != nil {
		slog.Error("failed to load user", err)
		return
	}
	timeRgx := regexp.MustCompile(`call_time_picker_(\d{2}:00)`)
	if userdata.Target.CallTime == "" && timeRgx.MatchString(event.Postback.Data) {
		userdata.Target.CallTime = timeRgx.FindStringSubmatch(event.Postback.Data)[1]
		_, err := ref.Set(ctx, userdata)
		if err != nil {
			slog.Error("failed to set user", err)
			return
		}
		l.lineBotSvc.ReplyTextMessage(event.ReplyToken, "最後にリマインドメッセージがある場合は入力してください。ない場合は「なし」と入力してください")
		return
	}

}
