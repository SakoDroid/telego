package telebot

import (
	"log"
	"os"

	cfg "github.com/SakoDroid/telebot/configs"
)

var Logger *log.Logger

func initTheLogger(botCfg *cfg.BotConfigs) {
	file, err := os.OpenFile(botCfg.LogFileAddress, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic("Could not init the logger. Reason : " + err.Error())
	}
	Logger = log.New(file, "", log.Ldate|log.Ltime)
}
