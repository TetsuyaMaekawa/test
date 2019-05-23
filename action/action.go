package action

import (
	"log"

	"github.com/heroku/test/dbaccess/mysql"

	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
)

// ResFollowEvent followEventに対して応答
func (i *InitLinebot) ResFollowEvent() {
	if _, err := i.Bot.ReplyMessage(i.Event.ReplyToken, linebot.NewTextMessage("友達追加ありがとうございます。")).Do(); err != nil {
		log.Print(err)
	}
}

// ResMessageEvent messageEventに対して応答
func (i *InitLinebot) ResMessageEvent() {
	switch message := i.Event.Message.(type) {
	case *linebot.TextMessage:
		i.resTextMessage(message)
	case *linebot.ImageMessage:
		i.resImageMessage()
	}
}

// ResPostBackEvent postBackEventに対して応答
func (i *InitLinebot) ResPostBackEvent() {
	if _, err := i.Bot.ReplyMessage(i.Event.ReplyToken, linebot.NewTemplateMessage(
		"this is a confilm template",
		linebot.NewConfirmTemplate(
			"key oa value",
			linebot.NewMessageAction(
				"key", "dbから抜き出したkey",
			),
			linebot.NewMessageAction(
				"key", "dbから抜き出したvalue",
			),
		),
	)).Do(); err != nil {
		log.Print(err)
	}
}

// resTextMessage textMessageの時に応答
func (i *InitLinebot) resTextMessage(message *linebot.TextMessage) {
	if message.Text == "情報" {
		userID := i.Event.Source.UserID
		profile, _ := i.Bot.GetProfile(userID).Do()
		if _, err := i.Bot.ReplyMessage(i.Event.ReplyToken, linebot.NewTextMessage("あなたの名前は「"+profile.DisplayName+"」\n"+"あなたのIDは「"+userID+"」です。")).Do(); err != nil {
			log.Print(err)
		}
	} else {
		if _, err := i.Bot.ReplyMessage(i.Event.ReplyToken, linebot.NewTextMessage("ユーザー情報を取得したい場合は「情報」と入力してください。")).Do(); err != nil {
			log.Print(err)
		}
	}
}

// resImageMessage imageMessageの時に応答
func (i *InitLinebot) resImageMessage() {
	m := mysql.Mytable{}
	i.DB.First(&m, "id=?", 1)
	if _, err := i.Bot.ReplyMessage(i.Event.ReplyToken, linebot.NewTemplateMessage(
		"this is a buttons template",
		linebot.NewButtonsTemplate(
			"",
			m.Name,
			"message text",
			linebot.NewMessageAction(
				"text",
				"text",
			),
			linebot.NewPostbackAction(
				"post back",
				"post back",
				"",
				""),
		),
	)).Do(); err != nil {
		log.Print(err)
	}
}

// InitLinebot ClientとEventを保持
type InitLinebot struct {
	Bot   *linebot.Client
	Event *linebot.Event
	DB    *gorm.DB
}
