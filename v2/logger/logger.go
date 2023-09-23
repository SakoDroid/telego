package logger

import (
	"log"
	"os"

	cfg "github.com/SakoDroid/telego/v2/configs"
)

type BotLogger struct {
	// logger is the default logger of the bot.
	logger    *log.Logger
	colorized bool
}

// Logger is the default logger of the bot.
// var Logger *log.Logger

// colorized indicates if logs should be colored, default is true.
var colorized = true

const (
	HEADER    string = "\033[95m"
	OKBLUE    string = "\033[94m"
	OKCYAN    string = "\033[96m"
	OKGREEN   string = "\033[92m"
	WARNING   string = "\033[93m"
	FAIL      string = "\033[91m"
	ENDC      string = "\033[0m"
	BOLD      string = "\033[1m"
	UNDERLINE string = "\033[4m"
)

// InitiTheLogger initializes the default logger of the bot.
func InitTheLogger(botCfg *cfg.BotConfigs) *BotLogger {
	logger := &BotLogger{
		colorized: true,
	}
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
	logger.logger = log.New(file, botCfg.BotName+" | ", log.Ldate|log.Ltime)
	return logger
}

// Log logs the given paramteres based on the defined format.
func (l *BotLogger) Log(header, space, content, after, headerColor, contentColor, afterColor string) {
	if colorized {
		text := "| " + headerColor + header + ENDC + space + contentColor + content + ENDC + " |" + afterColor + after + ENDC
		l.logger.Println(text)
	} else {
		text := "| " + header + space + content + "|" + after
		l.logger.Println(text)
	}

}

// GetRaw returns the internal logger
func (l *BotLogger) GetRaw() *log.Logger {
	return l.logger
}

// Uncolor, clears the colors of the logs.
func (l *BotLogger) Uncolor() {
	colorized = false
}

// Color adds color to the logs.
func (l *BotLogger) Color() {
	colorized = true
}
