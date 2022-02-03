package telego

import (
	"errors"
	"os"

	errs "github.com/SakoDroid/telego/errors"
	objs "github.com/SakoDroid/telego/objects"
)

//MediaGroup is a media group that can be sent
type MediaGroup struct {
	bot                     *Bot
	replyTo                 int
	allowSendingWihoutReply bool
	replyMarkup             objs.ReplyMarkup
	media                   []objs.InputMedia
	files                   []*os.File
}

//PhotoInserter is a tool for inserting photos into the MediaGroup.
type PhotoInserter struct {
	mg                 *MediaGroup
	caption, parseMode string
	captionEntities    []objs.MessageEntity
}

/*AddByFileIdOrURL adds this file by file id or url*/
func (pi *PhotoInserter) AddByFileIdOrURL(fileIdOrUrl string) {
	im := &objs.InputMediaPhoto{
		InputMediaDefault: fixTheDefault("photo", fileIdOrUrl, pi.caption, pi.parseMode, pi.captionEntities),
	}
	pi.mg.media = append(pi.mg.media, im)
}

/*AddByFile adds an existing file in the device*/
func (pi *PhotoInserter) AddByFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	im := &objs.InputMediaPhoto{
		InputMediaDefault: fixTheDefault("photo", "attach://"+stat.Name(), pi.caption, pi.parseMode, pi.captionEntities),
	}
	pi.mg.media = append(pi.mg.media, im)
	pi.mg.files = append(pi.mg.files, file)
	return nil
}

//VideoInserter is a tool for inserting videos into the MediaGroup.
type VideoInserter struct {
	mg                        *MediaGroup
	caption, parseMode, thumb string
	captionEntities           []objs.MessageEntity
	thumbFile                 *os.File
	width, height, duration   int
	supportsStreaming         bool
}

/*AddByFileIdOrURL adds this file by file id or url*/
func (vi *VideoInserter) AddByFileIdOrURL(fileIdOrUrl string) {
	im := &objs.InputMediaVideo{
		InputMediaDefault: fixTheDefault("video", fileIdOrUrl, vi.caption, vi.parseMode, vi.captionEntities),
		Thumb:             vi.thumb,
		SupportsStreaming: vi.supportsStreaming,
	}
	if vi.width != 0 {
		im.Width = vi.width
	}
	if vi.height != 0 {
		im.Height = vi.height
	}
	if vi.duration != 0 {
		im.Duration = vi.duration
	}
	vi.mg.media = append(vi.mg.media, im)
	if vi.thumbFile != nil {
		vi.mg.files = append(vi.mg.files, vi.thumbFile)
	}
}

/*AddByFile adds an existing file in the device*/
func (vi *VideoInserter) AddByFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	im := &objs.InputMediaVideo{
		InputMediaDefault: fixTheDefault("video", "attach://"+stat.Name(), vi.caption, vi.parseMode, vi.captionEntities),
		Thumb:             vi.thumb,
		SupportsStreaming: vi.supportsStreaming,
	}
	if vi.width != 0 {
		im.Width = vi.width
	}
	if vi.height != 0 {
		im.Height = vi.height
	}
	if vi.duration != 0 {
		im.Duration = vi.duration
	}
	vi.mg.media = append(vi.mg.media, im)
	vi.mg.files = append(vi.mg.files, file)
	if vi.thumbFile != nil {
		vi.mg.files = append(vi.mg.files, vi.thumbFile)
	}
	return nil
}

/*SetThumbnail sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (vi *VideoInserter) SetThumbnail(fileIdOrURL string) {
	vi.thumb = fileIdOrURL
}

/*SetThumbnailFile sets the tumbnail of the file. It takes a file existing on the device*/
func (vi *VideoInserter) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	vi.thumbFile = file
	vi.thumb = "attach://" + stat.Name()
	return nil
}

//AnimationInserter is a tool for inserting animations into the MediaGroup.
type AnimationInserter struct {
	mg                        *MediaGroup
	caption, parseMode, thumb string
	captionEntities           []objs.MessageEntity
	thumbFile                 *os.File
	width, height, duration   int
}

/*AddByFileIdOrURL adds this file by file id or url*/
func (ai *AnimationInserter) AddByFileIdOrURL(fileIdOrUrl string) {
	im := &objs.InputMediaAnimation{
		InputMediaDefault: fixTheDefault("animation", fileIdOrUrl, ai.caption, ai.parseMode, ai.captionEntities),
		Thumb:             ai.thumb,
	}
	if ai.width != 0 {
		im.Width = ai.width
	}
	if ai.height != 0 {
		im.Height = ai.height
	}
	if ai.duration != 0 {
		im.Duration = ai.duration
	}
	ai.mg.media = append(ai.mg.media, im)
	if ai.thumbFile != nil {
		ai.mg.files = append(ai.mg.files, ai.thumbFile)
	}
}

/*AddByFile adds an existing file in the device*/
func (ai *AnimationInserter) AddByFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	im := &objs.InputMediaAnimation{
		InputMediaDefault: fixTheDefault("animation", "attach://"+stat.Name(), ai.caption, ai.parseMode, ai.captionEntities),
		Thumb:             ai.thumb,
	}
	if ai.width != 0 {
		im.Width = ai.width
	}
	if ai.height != 0 {
		im.Height = ai.height
	}
	if ai.duration != 0 {
		im.Duration = ai.duration
	}
	ai.mg.media = append(ai.mg.media, im)
	ai.mg.files = append(ai.mg.files, file)
	if ai.thumbFile != nil {
		ai.mg.files = append(ai.mg.files, ai.thumbFile)
	}
	return nil
}

/*SetThumbnail sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (ai *AnimationInserter) SetThumbnail(fileIdOrURL string) {
	ai.thumb = fileIdOrURL
}

/*SetThumbnailFile sets the tumbnail of the file. It takes a file existing on the device*/
func (ai *AnimationInserter) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	ai.thumbFile = file
	ai.thumb = "attach://" + stat.Name()
	return nil
}

//AudioInserter is a tool for inserting audios into the MediaGroup.
type AudioInserter struct {
	mg                                          *MediaGroup
	caption, parseMode, thumb, performer, title string
	captionEntities                             []objs.MessageEntity
	thumbFile                                   *os.File
	duration                                    int
}

/*AddByFileIdOrURL adds this file by file id or url*/
func (ai *AudioInserter) AddByFileIdOrURL(fileIdOrUrl string) {
	im := &objs.InputMediaAudio{
		InputMediaDefault: fixTheDefault("audio", fileIdOrUrl, ai.caption, ai.parseMode, ai.captionEntities),
		Thumb:             ai.thumb,
		Performer:         ai.performer,
		Title:             ai.title,
	}
	if ai.duration != 0 {
		im.Duration = ai.duration
	}
	ai.mg.media = append(ai.mg.media, im)
	if ai.thumbFile != nil {
		ai.mg.files = append(ai.mg.files, ai.thumbFile)
	}
}

/*AddByFile adds an existing file in the device*/
func (ai *AudioInserter) AddByFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	im := &objs.InputMediaAudio{
		InputMediaDefault: fixTheDefault("audio", "attach://"+stat.Name(), ai.caption, ai.parseMode, ai.captionEntities),
		Thumb:             ai.thumb,
		Performer:         ai.performer,
		Title:             ai.title,
	}
	if ai.duration != 0 {
		im.Duration = ai.duration
	}
	ai.mg.media = append(ai.mg.media, im)
	ai.mg.files = append(ai.mg.files, file)
	if ai.thumbFile != nil {
		ai.mg.files = append(ai.mg.files, ai.thumbFile)
	}
	return nil
}

/*SetThumbnail sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (ai *AudioInserter) SetThumbnail(fileIdOrURL string) {
	ai.thumb = fileIdOrURL
}

/*SetThumbnailFile sets the tumbnail of the file. It takes a file existing on the device*/
func (ai *AudioInserter) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	ai.thumbFile = file
	ai.thumb = "attach://" + stat.Name()
	return nil
}

//DocumentInserter is a tool for inserting documents into the MediaGroup.
type DocumentInserter struct {
	mg                          *MediaGroup
	caption, parseMode, thumb   string
	captionEntities             []objs.MessageEntity
	thumbFile                   *os.File
	disableContentTypeDetection bool
}

/*AddByFileIdOrURL adds this file by file id or url*/
func (di *DocumentInserter) AddByFileIdOrURL(fileIdOrUrl string) {
	im := &objs.InputMediaDocument{
		InputMediaDefault:           fixTheDefault("document", fileIdOrUrl, di.caption, di.parseMode, di.captionEntities),
		Thumb:                       di.thumb,
		DisableContentTypeDetection: di.disableContentTypeDetection,
	}
	di.mg.media = append(di.mg.media, im)
	if di.thumbFile != nil {
		di.mg.files = append(di.mg.files, di.thumbFile)
	}
}

/*AddByFile adds an existing file in the device*/
func (di *DocumentInserter) AddByFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	im := &objs.InputMediaDocument{
		InputMediaDefault:           fixTheDefault("document", "attach://"+stat.Name(), di.caption, di.parseMode, di.captionEntities),
		Thumb:                       di.thumb,
		DisableContentTypeDetection: di.disableContentTypeDetection,
	}
	di.mg.media = append(di.mg.media, im)
	di.mg.files = append(di.mg.files, file)
	if di.thumbFile != nil {
		di.mg.files = append(di.mg.files, di.thumbFile)
	}
	return nil
}

/*SetThumbnail sets the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (di *DocumentInserter) SetThumbnail(fileIdOrURL string) {
	di.thumb = fileIdOrURL
}

/*SetThumbnailFile sets the tumbnail of the file. It takes a file existing on the device*/
func (di *DocumentInserter) SetThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	di.thumbFile = file
	di.thumb = "attach://" + stat.Name()
	return nil
}

/*Send sends this album (to all types of chat but channels, to send to channels use "SendToChannel" method)

--------------------

Official telegram doc :

Use this method to send a group of photos, videos, documents or audios as an album. Documents and audio files can be only grouped in an album with messages of the same type. On success, an array of Messages that were sent is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (mg *MediaGroup) Send(chatId int, silent, protectContent bool) (*objs.SendMediaGroupMethodResult, error) {
	if len(mg.media) < 2 {
		return nil, errors.New("the number os medias should be greater than 1")
	}
	return mg.bot.apiInterface.SendMediaGroup(
		chatId, "", mg.replyTo, mg.media, silent, mg.allowSendingWihoutReply, protectContent,
		mg.replyMarkup, mg.files...,
	)
}

/*Send sends this album to a channel.


--------------------

Official telegram doc :

Use this method to send a group of photos, videos, documents or audios as an album. Documents and audio files can be only grouped in an album with messages of the same type. On success, an array of Messages that were sent is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (mg *MediaGroup) SendToChannel(chatId string, silent, protectContent bool) (*objs.SendMediaGroupMethodResult, error) {
	if len(mg.media) < 2 {
		return nil, errors.New("the number os medias should be greater than 1")
	}
	return mg.bot.apiInterface.SendMediaGroup(
		0, chatId, mg.replyTo, mg.media, silent, mg.allowSendingWihoutReply, protectContent,
		mg.replyMarkup, mg.files...,
	)
}

/*AddPhoto returns a PhotoInserter to add a photo to the album*/
func (mg *MediaGroup) AddPhoto(caption, parseMode string, captionEntitie []objs.MessageEntity) (*PhotoInserter, error) {
	if len(mg.media) == 10 {
		return nil, &errs.MediaGroupFullError{}
	}
	return &PhotoInserter{mg: mg, caption: caption, parseMode: parseMode, captionEntities: captionEntitie}, nil
}

/*AddVideo returns a VideoInserter to add a video to the album*/
func (mg *MediaGroup) AddVideo(caption, parseMode string, width, height, duration int, supportsStreaming bool, captionEntitie []objs.MessageEntity) (*VideoInserter, error) {
	if len(mg.media) == 10 {
		return nil, &errs.MediaGroupFullError{}
	}
	return &VideoInserter{mg: mg, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, width: width, height: height, duration: duration, supportsStreaming: supportsStreaming}, nil
}

/*AddAnimation returns an AnimationInserter to add an animation to the album*/
func (mg *MediaGroup) AddAnimation(caption, parseMode string, width, height, duration int, captionEntitie []objs.MessageEntity) (*AnimationInserter, error) {
	if len(mg.media) == 10 {
		return nil, &errs.MediaGroupFullError{}
	}
	return &AnimationInserter{mg: mg, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, width: width, height: height, duration: duration}, nil
}

/*AddAudio returns an AudioInserter to add an audio to the album*/
func (mg *MediaGroup) AddAudio(caption, parseMode, performer, title string, duration int, captionEntitie []objs.MessageEntity) (*AudioInserter, error) {
	if len(mg.media) == 10 {
		return nil, &errs.MediaGroupFullError{}
	}
	return &AudioInserter{mg: mg, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, performer: performer, title: title, duration: duration}, nil
}

/*AddDocument returns a DocumentInserter to add a document to the album*/
func (mg *MediaGroup) AddDocument(caption, parseMode string, disableContentTypeDetection bool, captionEntitie []objs.MessageEntity) (*DocumentInserter, error) {
	if len(mg.media) == 10 {
		return nil, &errs.MediaGroupFullError{}
	}
	return &DocumentInserter{mg: mg, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, disableContentTypeDetection: disableContentTypeDetection}, nil
}

func fixTheDefault(tp, media, caption, parseMode string, captionEnt []objs.MessageEntity) objs.InputMediaDefault {
	return objs.InputMediaDefault{Type: tp, Media: media, Caption: caption, ParseMode: parseMode, CaptionEntities: captionEnt}
}
