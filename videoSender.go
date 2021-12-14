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
	thumbFile                               *os.File
}

/*Sends a file that already exists on telegram servers*/
func (vs *VideoSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return vs.bot.apiInterface.SendVideo(
		vs.chatIdInt, vs.chatidString, fileId,
		nil, vs.caption, vs.parseMode, vs.replyTo, vs.thumb, vs.thumbFile, silent, vs.allowSendingWihoutReply,
		vs.captionEntities, vs.duration, vs.supportsStreaming, vs.replyMarkup,
	)
}

/*Sends a file on the web. The file is downloaded on telegram server and then will be sent to the chat.*/
func (vs *VideoSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return vs.bot.apiInterface.SendVideo(
		vs.chatIdInt, vs.chatidString, url,
		nil, vs.caption, vs.parseMode, vs.replyTo, vs.thumb, vs.thumbFile, silent, vs.allowSendingWihoutReply,
		vs.captionEntities, vs.duration, vs.supportsStreaming, vs.replyMarkup,
	)
}

/*Sends a file that is located in this device.*/
func (vs *VideoSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return vs.bot.apiInterface.SendVideo(
		vs.chatIdInt, vs.chatidString, "attach://"+stat.Name(),
		file, vs.caption, vs.parseMode, vs.replyTo, vs.thumb, vs.thumbFile, silent, vs.allowSendingWihoutReply,
		vs.captionEntities, vs.duration, vs.supportsStreaming, vs.replyMarkup,
	)
}

/*This method sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (vs *VideoSender) SetThumbnail(fileIdOrURL string) {
	vs.thumb = fileIdOrURL
}

/*This method sets the thumbnail of the file. It takes a file existing on the device*/
func (vs *VideoSender) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	vs.thumbFile = file
	vs.thumb = "attach://" + stat.Name()
	return nil
}
