package telego

import (
	"errors"
	"os"

	objs "github.com/SakoDroid/telego/objects"
)

type MediaType int

const (
	PHOTO     MediaType = 1
	VIDEO     MediaType = 2
	AUDIO     MediaType = 3
	ANIMATION MediaType = 4
	DOCUMENT  MediaType = 5
	VIDEONOTE MediaType = 6
	VOICE     MediaType = 7
	STICKER   MediaType = 8
)

//MediaSender is a tool for sending media messages.
type MediaSender struct {
	bot                                                       *Bot
	chatIdInt                                                 int
	mediaType                                                 MediaType
	chatidString, caption, parseMode, thumb, performer, title string
	replyTo                                                   int
	captionEntities                                           []objs.MessageEntity
	allowSendingWihoutReply                                   bool
	replyMarkup                                               objs.ReplyMarkup
	duration, length, width, height                           int
	supportsStreaming, disableContentTypeDetection            bool
	thumbFile                                                 *os.File
}

/*SendByFileIdOrUrl sends a file that already exists on telegram servers (file id) or a url on the web.*/
func (ms *MediaSender) SendByFileIdOrUrl(fileIdOrUrl string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	switch ms.mediaType {
	case PHOTO:
		return ms.bot.apiInterface.SendPhoto(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl, nil, ms.caption, ms.parseMode,
			ms.replyTo, silent, ms.allowSendingWihoutReply, protectContent, ms.replyMarkup, ms.captionEntities,
		)
	case VIDEO:
		return ms.bot.apiInterface.SendVideo(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl,
			nil, ms.caption, ms.parseMode, ms.replyTo, ms.thumb, ms.thumbFile, silent, ms.allowSendingWihoutReply, protectContent,
			ms.captionEntities, ms.duration, ms.supportsStreaming, ms.replyMarkup,
		)
	case AUDIO:
		return ms.bot.apiInterface.SendAudio(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl, nil, ms.caption, ms.parseMode,
			ms.replyTo, ms.thumb, ms.thumbFile, silent, ms.allowSendingWihoutReply, protectContent,
			ms.captionEntities, ms.duration, ms.performer, ms.title, ms.replyMarkup,
		)
	case ANIMATION:
		return ms.bot.apiInterface.SendAnimation(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl, nil, ms.caption, ms.parseMode,
			ms.width, ms.height, ms.duration, ms.replyTo, ms.thumb, ms.thumbFile,
			silent, ms.allowSendingWihoutReply, protectContent, ms.captionEntities, ms.replyMarkup,
		)
	case DOCUMENT:
		return ms.bot.apiInterface.SendDocument(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl, nil, ms.caption, ms.parseMode,
			ms.replyTo, ms.thumb, ms.thumbFile, silent, ms.allowSendingWihoutReply, protectContent, ms.captionEntities,
			ms.disableContentTypeDetection, ms.replyMarkup,
		)
	case VIDEONOTE:
		return ms.bot.apiInterface.SendVideoNote(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl, nil, ms.caption, ms.parseMode,
			ms.length, ms.duration, ms.replyTo, ms.thumb, ms.thumbFile, silent,
			ms.allowSendingWihoutReply, protectContent, ms.captionEntities, ms.replyMarkup,
		)
	case VOICE:
		return ms.bot.apiInterface.SendVoice(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl, nil, ms.caption, ms.parseMode,
			ms.duration, ms.replyTo, silent, ms.allowSendingWihoutReply, protectContent, ms.captionEntities, ms.replyMarkup,
		)
	case STICKER:
		return ms.bot.apiInterface.SendSticker(
			ms.chatIdInt, ms.chatidString, fileIdOrUrl, silent, ms.allowSendingWihoutReply, protectContent,
			ms.replyTo, ms.replyMarkup, nil,
		)
	default:
		return nil, errors.New("wrong media type")
	}

}

/*SendByFile sends a file that is located in this device.*/
func (ms *MediaSender) SendByFile(file *os.File, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	switch ms.mediaType {
	case PHOTO:
		return ms.bot.apiInterface.SendPhoto(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(), file, ms.caption, ms.parseMode,
			ms.replyTo, silent, ms.allowSendingWihoutReply, protectContent, ms.replyMarkup, ms.captionEntities,
		)
	case VIDEO:
		return ms.bot.apiInterface.SendVideo(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(),
			file, ms.caption, ms.parseMode, ms.replyTo, ms.thumb, ms.thumbFile, silent, ms.allowSendingWihoutReply, protectContent,
			ms.captionEntities, ms.duration, ms.supportsStreaming, ms.replyMarkup,
		)
	case AUDIO:
		return ms.bot.apiInterface.SendAudio(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(), file, ms.caption, ms.parseMode,
			ms.replyTo, ms.thumb, ms.thumbFile, silent, ms.allowSendingWihoutReply, protectContent,
			ms.captionEntities, ms.duration, ms.performer, ms.title, ms.replyMarkup,
		)
	case ANIMATION:
		return ms.bot.apiInterface.SendAnimation(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(), file, ms.caption, ms.parseMode,
			ms.width, ms.height, ms.duration, ms.replyTo, ms.thumb, ms.thumbFile,
			silent, ms.allowSendingWihoutReply, protectContent, ms.captionEntities, ms.replyMarkup,
		)
	case DOCUMENT:
		return ms.bot.apiInterface.SendDocument(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(), file, ms.caption, ms.parseMode,
			ms.replyTo, ms.thumb, ms.thumbFile, silent, ms.allowSendingWihoutReply, protectContent, ms.captionEntities,
			ms.disableContentTypeDetection, ms.replyMarkup,
		)
	case VIDEONOTE:
		return ms.bot.apiInterface.SendVideoNote(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(), file, ms.caption, ms.parseMode,
			ms.length, ms.duration, ms.replyTo, ms.thumb, ms.thumbFile, silent,
			ms.allowSendingWihoutReply, protectContent, ms.captionEntities, ms.replyMarkup,
		)
	case VOICE:
		return ms.bot.apiInterface.SendVoice(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(), file, ms.caption, ms.parseMode,
			ms.duration, ms.replyTo, silent, ms.allowSendingWihoutReply, protectContent, ms.captionEntities, ms.replyMarkup,
		)
	case STICKER:
		return ms.bot.apiInterface.SendSticker(
			ms.chatIdInt, ms.chatidString, "attach://"+stat.Name(), silent, ms.allowSendingWihoutReply, protectContent,
			ms.replyTo, ms.replyMarkup, file,
		)
	default:
		return nil, errors.New("wrong media type")
	}
}

/*SetThumbnail sets the tumbnail of the file. It takes a file id or a url. If you want to send a file use "setThumbnailFile" instead.
If this media does not support thumbnail, the thumbnail will be ignored.*/
func (ms *MediaSender) SetThumbnail(fileIdOrURL string) {
	ms.thumb = fileIdOrURL
}

/*SetThumbnailFile sets the thumbnail of the file. It takes a file existing on the device.
If this media does not support thumbnail, the thumbnail will be ignored.*/
func (ms *MediaSender) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	ms.thumbFile = file
	ms.thumb = "attach://" + stat.Name()
	return nil
}
