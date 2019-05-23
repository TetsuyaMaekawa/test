package main

import (
	"github.com/heroku/test/dbaccess/mysql"
	"github.com/heroku/test/handler"
	"github.com/zenazn/goji"
)

func main() {
	mysql.OpenMySQL()
	// Postのルーティング
	goji.Post("/callback", handler.LinebotHandler)
	goji.Serve()
}
