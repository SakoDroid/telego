package logger

import (
	"log"
	"os"

	cfg "github.com/SakoDroid/telego/configs"
)

//Logger is the default logger of the bot.
var Logger *log.Logger

//InitiTheLogger initializes the default logger of the bot.
func InitTheLogger(botCfg *cfg.BotConfigs) {
	if Logger == nil {
		var file *os.File
		if botCfg.LogFileAddress == "STDOUT" {
			file = os.Stdout
		} else {
			var err error
			file, err = os.OpenFile(botCfg.LogFileAddress, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
			if err != nil {
				panic("Could not init the logger. Reason : " + err.Error())
			}
		}
		Logger = log.New(file, "", log.Ldate|log.Ltime)
	}
}
