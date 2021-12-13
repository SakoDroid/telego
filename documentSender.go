package telebot

import (
	"os"

	objs "github.com/SakoDroid/telebot/objects"
)

type DocumentSender struct {
	bot                                     *Bot
	chatIdInt                               int
	chatidString, caption, parseMode, thumb string
	replyTo                                 int
	captionEntities                         []objs.MessageEntity
	allowSendingWihoutReply                 bool
	replyMarkup                             objs.ReplyMarkup
	disableContentTypeDetection             bool
}

func (ds *DocumentSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return ds.bot.apiInterface.SendDocument(
		ds.chatIdInt, ds.chatidString, fileId, nil, ds.caption, ds.parseMode,
		ds.replyTo, ds.thumb, nil, silent, ds.allowSendingWihoutReply, ds.captionEntities,
		ds.disableContentTypeDetection, ds.replyMarkup,
	)
}

func (ds *DocumentSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return ds.bot.apiInterface.SendDocument(
		ds.chatIdInt, ds.chatidString, url, nil, ds.caption, ds.parseMode,
		ds.replyTo, ds.thumb, nil, silent, ds.allowSendingWihoutReply, ds.captionEntities,
		ds.disableContentTypeDetection, ds.replyMarkup,
	)
}

func (ds *DocumentSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return ds.bot.apiInterface.SendDocument(
		ds.chatIdInt, ds.chatidString, "attach://"+stat.Name(), file, ds.caption, ds.parseMode,
		ds.replyTo, ds.thumb, nil, silent, ds.allowSendingWihoutReply, ds.captionEntities,
		ds.disableContentTypeDetection, ds.replyMarkup,
	)
}
