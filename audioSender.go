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
	thumbFile                                                 *os.File
}

/*Sends a file that already exists on telegram servers*/
func (as *AudioSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return as.bot.apiInterface.SendAudio(
		as.chatIdInt, as.chatidString, fileId, nil, as.caption, as.parseMode,
		as.replyTo, as.thumb, as.thumbFile, silent, as.allowSendingWihoutReply,
		as.captionEntities, as.duration, as.performer, as.title, as.replyMarkup,
	)
}

/*Sends a file on the web. The file is downloaded on telegram server and then will be sent to the chat.*/
func (as *AudioSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return as.bot.apiInterface.SendAudio(
		as.chatIdInt, as.chatidString, url, nil, as.caption, as.parseMode,
		as.replyTo, as.thumb, as.thumbFile, silent, as.allowSendingWihoutReply,
		as.captionEntities, as.duration, as.performer, as.title, as.replyMarkup,
	)
}

/*Sends a file that is located in this device.*/
func (as *AudioSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return as.bot.apiInterface.SendAudio(
		as.chatIdInt, as.chatidString, "attach://"+stat.Name(), file, as.caption, as.parseMode,
		as.replyTo, as.thumb, as.thumbFile, silent, as.allowSendingWihoutReply,
		as.captionEntities, as.duration, as.performer, as.title, as.replyMarkup,
	)
}

/*This method sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (as *AudioSender) SetThumbnail(fileIdOrURL string) {
	as.thumb = fileIdOrURL
}

/*This method sets the thumbnail of the file. It takes a file existing on the device*/
func (as *AudioSender) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	as.thumbFile = file
	as.thumb = "attach://" + stat.Name()
	return nil
}
