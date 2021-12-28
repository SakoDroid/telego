package telebot

import (
	"log"
	"os"

	cfg "github.com/SakoDroid/telego/configs"
)

var Logger *log.Logger

func InitTheLogger(botCfg *cfg.BotConfigs) {
	if Logger == nil {
		file, err := os.OpenFile(botCfg.LogFileAddress, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic("Could not init the logger. Reason : " + err.Error())
		}
		Logger = log.New(file, "", log.Ldate|log.Ltime)
	}
}
