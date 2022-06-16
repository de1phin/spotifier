package main

import (
	"database/sql"
	"../dbmaster"
	"../easybot"
	"log"

	tgbotapi "github.com/telegram-bot-api-bot-api-5.0"
)

var bot *tgbotapi.BotAPI
var db	*sql.DB

const tgtoken string = "1489112563:AAH99DLrk9DWYmIHUitiCKFT1pnmW-KddMo"

func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	db, err = dbmaster.DBConnect("Spotifier", "root", "root")
	if err != nil {
		panic(err)
	}
	InitImgGen()
}

func main() {
	easybot.SEPARATOR = " "

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	InitAllHandlers()
	InitAllInstances()

	easybot.Run(&updates)
}