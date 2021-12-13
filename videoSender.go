package telebot

import (
	"os"

	objs "github.com/SakoDroid/telebot/objects"
)

type VideoSender struct {
	bot                                     *Bot
	chatIdInt                               int
	chatidString, caption, parseMode, thumb string
	replyTo                                 int
	captionEntities                         []objs.MessageEntity
	allowSendingWihoutReply                 bool
	replyMarkup                             objs.ReplyMarkup
	duration                                int
	supportsStreaming                       bool
}

func (vs *VideoSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return vs.bot.apiInterface.SendVideo(
		vs.chatIdInt, vs.chatidString, fileId,
		nil, vs.caption, vs.parseMode, vs.replyTo, vs.thumb, nil, silent, vs.allowSendingWihoutReply,
		vs.captionEntities, vs.duration, vs.supportsStreaming, vs.replyMarkup,
	)
}

func (vs *VideoSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return vs.bot.apiInterface.SendVideo(
		vs.chatIdInt, vs.chatidString, url,
		nil, vs.caption, vs.parseMode, vs.replyTo, vs.thumb, nil, silent, vs.allowSendingWihoutReply,
		vs.captionEntities, vs.duration, vs.supportsStreaming, vs.replyMarkup,
	)
}

func (vs *VideoSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return vs.bot.apiInterface.SendVideo(
		vs.chatIdInt, vs.chatidString, "attach://"+stat.Name(),
		file, vs.caption, vs.parseMode, vs.replyTo, vs.thumb, nil, silent, vs.allowSendingWihoutReply,
		vs.captionEntities, vs.duration, vs.supportsStreaming, vs.replyMarkup,
	)
}
