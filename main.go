package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"

	// SDK追加
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// LINE bot instanceの作成
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
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
			// 友達追加された時の振る舞い
			if event.Type == linebot.EventTypeFollow {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("初めまして。よろしくお願いします。")).Do(); err != nil {
					log.Print(err)
				}
			}
			// メッセージを受けた場合のふるまい
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "情報" {
					userID := event.Source.UserID
					// ユーザーIDからプロフィールを取得
					profile, _ := bot.GetProfile(userID).Do()
					if err != nil {
						log.Print(err)
					}
					// 構造体に値をセット
					profileStruct := personalInfo{profile.DisplayName, profile.UserID}
					// 情報と入力された場合に自己情報を返す
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("あなたのユーザ名は："+profileStruct.name+"\n"+"あなたのユーザーIDは："+profileStruct.id)).Do(); err != nil {
						log.Print(err)
					}
				} else {
					// その他のメッセージを受けた場合はhelpを返す
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("プロフィール情報が見たい場合は「情報」と入力してください。")).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	router.Run(":" + port)
}

// 自己情報を保持する構造体
type personalInfo struct {
	name string
	id   string
}
