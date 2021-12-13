package telebot

import (
	"os"

	objs "github.com/SakoDroid/telebot/objects"
)

type AudioSender struct {
	bot                                                       *Bot
	chatIdInt                                                 int
	chatidString, caption, parseMode, thumb, performer, title string
	replyTo                                                   int
	captionEntities                                           []objs.MessageEntity
	allowSendingWihoutReply                                   bool
	replyMarkup                                               objs.ReplyMarkup
	duration                                                  int
}

func (as *AudioSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return as.bot.apiInterface.SendAudio(
		as.chatIdInt, as.chatidString, fileId, nil, as.caption, as.parseMode,
		as.replyTo, as.thumb, nil, silent, as.allowSendingWihoutReply,
		as.captionEntities, as.duration, as.performer, as.title, as.replyMarkup,
	)
}

func (as *AudioSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return as.bot.apiInterface.SendAudio(
		as.chatIdInt, as.chatidString, url, nil, as.caption, as.parseMode,
		as.replyTo, as.thumb, nil, silent, as.allowSendingWihoutReply,
		as.captionEntities, as.duration, as.performer, as.title, as.replyMarkup,
	)
}

func (as *AudioSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return as.bot.apiInterface.SendAudio(
		as.chatIdInt, as.chatidString, "attach://"+stat.Name(), file, as.caption, as.parseMode,
		as.replyTo, as.thumb, nil, silent, as.allowSendingWihoutReply,
		as.captionEntities, as.duration, as.performer, as.title, as.replyMarkup,
	)
}
