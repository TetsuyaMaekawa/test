package handler

import (
	"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/heroku/test/action"
	"github.com/jinzhu/gorm"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

// Linebothandler LINEからのリクエストを受けて応答をハンドリング
func (in *InitDB) Linebothandler() {

	// client生成
	bot, err := linebot.New("0189c809a76170e6c965b62ac5c9f670",
		"hJ5OAGDvemzFZidHYjg1Ihr5SoHs9eqsgUuok/LoW4uXzKD3lEZpqyqDMKti8Q/bp0rb4aVW2zsjFroGMoi5xTZqdWVrGy/CQE/EbozdNI3+Fyvq7sd4O/5EHyFpZ9mMwA7snSk+JzX8WJjNyXUJJAdB04t89/1O/w1cDnyilFU=",
	)
	if err != nil {
		log.Print(err)
		return
	}

	i := action.InitLinebot{Bot: bot, DB: in.DB, RD: in.RD}

	// // Postのルーティング
	goji.Post("/callback", func(c web.C, w http.ResponseWriter, r *http.Request) {

		// requesut取得
		events, err := bot.ParseRequest(r)
		if err != nil {
			log.Print(err)
		}
		// event毎に処理分岐
		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeFollow:
				i.ResFollowEvent(event)
			case linebot.EventTypeMessage:
				i.ResMessageEvent(event)
			case linebot.EventTypePostback:
				i.ResPostBackEvent(event)
			default:
			}
		}
	})
	goji.Serve()
}

// InitDB ...
type InitDB struct {
	DB *gorm.DB
	RD *redis.Pool
}
