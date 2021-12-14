package telebot

import (
	"os"

	objs "github.com/SakoDroid/telebot/objects"
)

type VideoNoteSender struct {
	bot                                     *Bot
	chatIdInt                               int
	chatidString, caption, parseMode, thumb string
	replyTo                                 int
	captionEntities                         []objs.MessageEntity
	allowSendingWihoutReply                 bool
	replyMarkup                             objs.ReplyMarkup
	duration, length                        int
	thumbFile                               *os.File
}

/*Sends a file that already exists on telegram servers*/
func (vns *VideoNoteSender) SendByFileId(fileId string, silent bool) (*objs.SendMethodsResult, error) {
	return vns.bot.apiInterface.SendVideoNote(
		vns.chatIdInt, vns.chatidString, fileId, nil, vns.caption, vns.parseMode,
		vns.length, vns.duration, vns.replyTo, vns.thumb, vns.thumbFile, silent,
		vns.allowSendingWihoutReply, vns.captionEntities, vns.replyMarkup,
	)
}

/*Sends a file that is located in this device.*/
func (vns *VideoNoteSender) SendByFile(file *os.File, silent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	return vns.bot.apiInterface.SendVideoNote(
		vns.chatIdInt, vns.chatidString, "attach://"+stat.Name(), file, vns.caption, vns.parseMode,
		vns.length, vns.duration, vns.replyTo, vns.thumb, vns.thumbFile, silent,
		vns.allowSendingWihoutReply, vns.captionEntities, vns.replyMarkup,
	)
}

/*This method sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (vns *VideoNoteSender) SetThumbnail(fileIdOrURL string) {
	vns.thumb = fileIdOrURL
}

/*This method sets the thumbnail of the file. It takes a file existing on the device*/
func (vns *VideoNoteSender) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	vns.thumbFile = file
	vns.thumb = "attach://" + stat.Name()
	return nil
}
