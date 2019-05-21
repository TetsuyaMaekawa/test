package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/heroku/go-getting-started/dbaccess/mysql"
	"github.com/heroku/go-getting-started/dbaccess/redis"
	_ "github.com/heroku/x/hmetrics/onload"

	// SDK追加
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	port := os.Getenv("PORT")

	port = "80"
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// LINE bot instanceの作成
	bot, err := linebot.New(
		"0189c809a76170e6c965b62ac5c9f670",
		"hJ5OAGDvemzFZidHYjg1Ihr5SoHs9eqsgUuok/LoW4uXzKD3lEZpqyqDMKti8Q/bp0rb4aVW2zsjFroGMoi5xTZqdWVrGy/CQE/EbozdNI3+Fyvq7sd4O/5EHyFpZ9mMwA7snSk+JzX8WJjNyXUJJAdB04t89/1O/w1cDnyilFU=",
		// os.Getenv("CHANNEL_SECRET"),
		// os.Getenv("CHANNEL_TOKEN"),
	)

	if err != nil {
		log.Fatal(err)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	// Line Messaging API用　Routing設定
	router.POST("/callback", func(c *gin.Context) {

		events, err := bot.ParseRequest(c.Request)

		if err != nil {
			if err == linebot.ErrInvalidSignature {
				log.Print(err)
			}
			return
		}
		for _, event := range events {
			// フォローイベントの場合
			if event.Type == linebot.EventTypeFollow {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("初めまして。よろしくお願いします。")).Do(); err != nil {
					log.Print(err)
				}
			}
			// ポストバックイベントの場合
			if event.Type == linebot.EventTypePostback {
				// redisに接続
				redisConn := redis.RedisConnection()
				redis.RedisSet("key1", "value1", 30, redisConn)
			}
			// メッセージイベントの場合
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				// テキストメッセージの場合
				case *linebot.TextMessage:
					if message.Text == "情報" {
						userID := event.Source.UserID
						// ユーザーIDからプロフィールを取得
						profile, _ := bot.GetProfile(userID).Do()
						if err != nil {
							log.Print(err)
						}
						// 情報と入力された場合に自己情報を返す
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("あなたのユーザ名は："+profile.DisplayName+"\n"+"あなたのユーザーIDは："+profile.UserID)).Do(); err != nil {
							log.Print(err)
						}
					} else {
						// その他のメッセージを受けた場合はhelpを返す
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("プロフィール情報が見たい場合は「情報」と入力してください。")).Do(); err != nil {
							log.Print(err)
						}
					}
				// 画像メッセージの場合
				case *linebot.ImageMessage:
					// dbに接続しIDから情報を抽出
					db := mysql.GormConnect()
					defer db.Close()
					r := mysql.RcvData{}
					r.ID = 1
					db.First(&r, "id=?", "1")
					// ボタンテンプレート
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage(
						"this is a buttons template",
						linebot.NewButtonsTemplate(
							"",
							r.Name,
							"messsage text",
							linebot.NewMessageAction("text", "text"),
							linebot.NewPostbackAction("postback", "postback", "", ""),
						),
					)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})
	router.Run(":" + port)
}
