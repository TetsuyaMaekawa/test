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

	// dbpoolの生成
	redisPool := redis.NewPool()
	db := mysql.GormConnect()
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
				// redisPool := redis.NewPool()
				redisKey := "key1"
				redisValue := "value1"
				// redisに接続
				redisConn := redisPool.Get()
				redis.RedisSet(redisKey, redisValue, 30, redisConn)
				// // イメージマップ
				// ibs := ImagemapBaseSize{1040, 1040}
				// ia := ImagemapArea{520, 0, 520, 1040}
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImagemapMessage("https://rootship.co.jp/images/logo.png",
				// 	"this is a imagemap",
				// 	linebot.ImagemapBaseSize(ibs),
				// 	linebot.NewMessageImagemapAction("test", linebot.ImagemapArea((ia))),
				// )).Do(); err != nil {
				// 	log.Print(err)
				// }
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage(
					"this is a buttons template",
					linebot.NewConfirmTemplate(
						"key oa value",
						linebot.NewMessageAction("key", redis.RedisGetKey(redisKey, redisConn)),
						linebot.NewMessageAction("value", redis.RedisGetValue(redisKey, redisConn)),
					),
				)).Do(); err != nil {
					log.Print(err)
				}
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
					// db := mysql.GormConnect()
					// dbからデータ取得
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

// // ImagemapBaseSize ...
// type ImagemapBaseSize struct {
// 	Width  int "json:\"width\""
// 	Height int "json:\"height\""
// }

// // ImagemapArea ...
// type ImagemapArea struct {
// 	X      int "json:\"x\""
// 	Y      int "json:\"y\""
// 	Width  int "json:\"width\""
// 	Height int "json:\"height\""
// }
