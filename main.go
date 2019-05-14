package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

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
			if event.Type == linebot.EventTypeMessage {
				respMessage = array(
					image,
					"https://tenshoku.mynavi.jp/sites/all/knowhow/heroes_file/img/top167_19.jpg",
					"https://tenshoku.mynavi.jp/sites/all/knowhow/heroes_file/img/top167_19.jpg",
				)
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(respMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	router.Run(":" + port)
}

// func getMessage(message *linebot.TextMessage) (rtnMessage string) {
// 	resMessage := [3]string{"ありがとう", "どういたしまして", "おやすみなさい"}

// 	// 乱数生成
// 	rand.Seed(time.Now().UnixNano())
// 	for {
// 		if math := rand.Intn(3); math != 3 {
// 			rtnMessage = resMessage[math]
// 			break
// 		}
// 	}
// 	return
}
