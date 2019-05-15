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

			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// if getResMessage(message.Text) == "情報" {
				if message.Text == "情報" {
					userID := event.Source.UserID
					// resp, _ := http.Get("https://api.line.me/v2/bot/profile/" + userID)
					// defer resp.Body.Close()
					// var profile profile
					// if err := json.NewDecoder(resp.Body).D\ecode(&profile); err != nil {
					// log.Print(err)
					// }
					res, _ := bot.GetProfile(userID).Do()
					userName := res.DisplayName

					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("あなたのユーザーIDは："+userID+"\n"+"あなたのユーザ名は："+userName)).Do(); err != nil {
						log.Print(err)
					}
				}
			}

			// // メッセージが送信されてきた時の振る舞い
			// if event.Type == linebot.EventTypeMessage {
			// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("そうですね。")).Do(); err != nil {
			// 		log.Print(err)
			// 	}
			// }
		}
	})

	router.Run(":" + port)
}

// JSONデコード用構造体
// type profile struct {
// 	displayName   string `json:"displayname"`
// 	userId        string `json:"userid"`
// 	pictureUrl    string `json:"pictureUrl"`
// 	statusMessage string `json:"statusMessage"`
// }
