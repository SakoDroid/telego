package telebot

import (
	"os"

	objs "github.com/SakoDroid/telebot/objects"
)

type PhotoSender struct {
	bot                              *Bot
	chatIdInt                        int
	chatidString, caption, parseMode string
	replyTo                          int
	captionEntities                  []objs.MessageEntity
	allowSendingWihoutReply          bool
	replyMarkup                      objs.ReplyMarkup
}

func (ps *PhotoSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return ps.bot.apiInterface.SendPhoto(
		ps.chatIdInt, ps.chatidString, fileId, nil, ps.caption, ps.parseMode,
		ps.replyTo, silent, ps.allowSendingWihoutReply, ps.replyMarkup, ps.captionEntities,
	)
}

func (ps *PhotoSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return ps.bot.apiInterface.SendPhoto(
		ps.chatIdInt, ps.chatidString, url, nil, ps.caption, ps.parseMode,
		ps.replyTo, silent, ps.allowSendingWihoutReply, ps.replyMarkup, ps.captionEntities,
	)
}

func (ps *PhotoSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return ps.bot.apiInterface.SendPhoto(
		ps.chatIdInt, ps.chatidString, "attach://"+stats.Name(), file, ps.caption, ps.parseMode,
		ps.replyTo, silent, ps.allowSendingWihoutReply, ps.replyMarkup, ps.captionEntities,
	)
}
