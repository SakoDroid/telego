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
	thumbFile                               *os.File
}

/*Sends a file that already exists on telegram servers*/
func (ds *DocumentSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return ds.bot.apiInterface.SendDocument(
		ds.chatIdInt, ds.chatidString, fileId, nil, ds.caption, ds.parseMode,
		ds.replyTo, ds.thumb, ds.thumbFile, silent, ds.allowSendingWihoutReply, ds.captionEntities,
		ds.disableContentTypeDetection, ds.replyMarkup,
	)
}

/*Sends a file on the web. The file is downloaded on telegram server and then will be sent to the chat.*/
func (ds *DocumentSender) SendByURL(url string, silent bool) (*objs.SendMethodsResult, error) {
	return ds.bot.apiInterface.SendDocument(
		ds.chatIdInt, ds.chatidString, url, nil, ds.caption, ds.parseMode,
		ds.replyTo, ds.thumb, ds.thumbFile, silent, ds.allowSendingWihoutReply, ds.captionEntities,
		ds.disableContentTypeDetection, ds.replyMarkup,
	)
}

/*Sends a file that is located in this device.*/
func (ds *DocumentSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return ds.bot.apiInterface.SendDocument(
		ds.chatIdInt, ds.chatidString, "attach://"+stat.Name(), file, ds.caption, ds.parseMode,
		ds.replyTo, ds.thumb, ds.thumbFile, silent, ds.allowSendingWihoutReply, ds.captionEntities,
		ds.disableContentTypeDetection, ds.replyMarkup,
	)
}

/*This method sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (ds *DocumentSender) SetThumbnail(fileIdOrURL string) {
	ds.thumb = fileIdOrURL
}

/*This method sets the thumbnail of the file. It takes a file existing on the device*/
func (ds *DocumentSender) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	ds.thumbFile = file
	ds.thumb = "attach://" + stat.Name()
	return nil
}
