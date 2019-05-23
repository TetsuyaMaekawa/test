package handler

import (
	"log"
	"net/http"

	"github.com/heroku/test/action"
	"github.com/heroku/test/dbaccess/mysql"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zenazn/goji/web"
)

// LinebotHandler LINEからのリクエストを受けて応答をハンドリング
func LinebotHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	// client生成
	bot, err := linebot.New("0189c809a76170e6c965b62ac5c9f670",
		"hJ5OAGDvemzFZidHYjg1Ihr5SoHs9eqsgUuok/LoW4uXzKD3lEZpqyqDMKti8Q/bp0rb4aVW2zsjFroGMoi5xTZqdWVrGy/CQE/EbozdNI3+Fyvq7sd4O/5EHyFpZ9mMwA7snSk+JzX8WJjNyXUJJAdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Print(err)
	}

	// requesut取得
	events, err := bot.ParseRequest(r)
	if err != nil {
		log.Print(err)
	}

	db := mysql.OpenMySQL()

	// event毎に処理分岐
	for _, event := range events {
		i := action.InitLinebot{Bot: bot, Event: event, DB: db}
		switch event.Type {
		case linebot.EventTypeFollow:
			i.ResFollowEvent()
		case linebot.EventTypeMessage:
			i.ResMessageEvent()
		case linebot.EventTypePostback:
			i.ResPostBackEvent()
		default:
		}
	}
}
