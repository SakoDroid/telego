package telebot

import (
	"os"

	objs "github.com/SakoDroid/telebot/objects"
)

type VoiceSender struct {
	bot                              *Bot
	chatIdInt                        int
	chatidString, caption, parseMode string
	replyTo                          int
	captionEntities                  []objs.MessageEntity
	allowSendingWihoutReply          bool
	replyMarkup                      objs.ReplyMarkup
	duration                         int
}

/*Sends a file that already exists on telegram servers*/
func (vs *VoiceSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return vs.bot.apiInterface.SendVoice(
		vs.chatIdInt, vs.chatidString, fileId, nil, vs.caption, vs.parseMode,
		vs.duration, vs.replyTo, silent, vs.allowSendingWihoutReply, vs.captionEntities, vs.replyMarkup,
	)
}

/*Sends a file on the web. The file is downloaded on telegram server and then will be sent to the chat.*/
func (vs *VoiceSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return vs.bot.apiInterface.SendVoice(
		vs.chatIdInt, vs.chatidString, url, nil, vs.caption, vs.parseMode,
		vs.duration, vs.replyTo, silent, vs.allowSendingWihoutReply, vs.captionEntities, vs.replyMarkup,
	)
}

/*Sends a file that is located in this device.*/
func (vs *VoiceSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return vs.bot.apiInterface.SendVoice(
		vs.chatIdInt, vs.chatidString, "attach://"+stat.Name(), file, vs.caption, vs.parseMode,
		vs.duration, vs.replyTo, silent, vs.allowSendingWihoutReply, vs.captionEntities, vs.replyMarkup,
	)
}
