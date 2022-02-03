package telego

import (
	"os"

	objs "github.com/SakoDroid/telego/objects"
)

//MessageEditor is a tool for editing messsages.
type MessageEditor struct {
	bot          *Bot
	chatIdInt    int
	chatIdString string
}

//PhotoEditor is a tool for editing photos.
type PhotoEditor struct {
	mg                                  *MessageEditor
	messageId                           int
	inlineMessageId, caption, parseMode string
	captionEntities                     []objs.MessageEntity
	replyMarkup                         *objs.InlineKeyboardMarkup
}

/*EditByFileIdOrURL edits this photo by file id or url*/
func (pi *PhotoEditor) EditByFileIdOrURL(fileIdOrUrl string) (*objs.DefaultResult, error) {
	im := &objs.InputMediaPhoto{
		InputMediaDefault: fixTheDefault("photo", fileIdOrUrl, pi.caption, pi.parseMode, pi.captionEntities),
	}
	return pi.mg.editMedia(pi.messageId, pi.inlineMessageId, im, pi.replyMarkup, nil)
}

/*EditByFile edits this photo with an existing file in the device*/
func (pi *PhotoEditor) EditByFile(file *os.File) (*objs.DefaultResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	im := &objs.InputMediaPhoto{
		InputMediaDefault: fixTheDefault("photo", "attach://"+stat.Name(), pi.caption, pi.parseMode, pi.captionEntities),
	}
	return pi.mg.editMedia(pi.messageId, pi.inlineMessageId, im, pi.replyMarkup, file)
}

//VideoEditor is a tool for editing videos.
type VideoEditor struct {
	mg                                         *MessageEditor
	messageId                                  int
	inlineMessageId, caption, parseMode, thumb string
	captionEntities                            []objs.MessageEntity
	thumbFile                                  *os.File
	width, height, duration                    int
	supportsStreaming                          bool
	replyMarkup                                *objs.InlineKeyboardMarkup
}

/*EditByFileIdOrURL edits this video by file id or url*/
func (vi *VideoEditor) EditByFileIdOrURL(fileIdOrUrl string) (*objs.DefaultResult, error) {
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
	return vi.mg.editMedia(vi.messageId, vi.inlineMessageId, im, vi.replyMarkup, nil, vi.thumbFile)
}

/*EditByFile edits this video by file in the device*/
func (vi *VideoEditor) EditByFile(file *os.File) (*objs.DefaultResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
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
	return vi.mg.editMedia(vi.messageId, vi.inlineMessageId, im, vi.replyMarkup, file, vi.thumbFile)
}

/*EditThumbnail edits the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (vi *VideoEditor) EditThumbnail(fileIdOrURL string) {
	vi.thumb = fileIdOrURL
}

/*EditThumbnailFile edits the thumbnail of the file. It takes a file existing on the device*/
func (vi *VideoEditor) EditThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	vi.thumbFile = file
	vi.thumb = "attach://" + stat.Name()
	return nil
}

//AnimationEditor is a tool for editing animations.
type AnimationEditor struct {
	mg                                         *MessageEditor
	messageId                                  int
	inlineMessageId, caption, parseMode, thumb string
	captionEntities                            []objs.MessageEntity
	thumbFile                                  *os.File
	width, height, duration                    int
	replyMarkup                                *objs.InlineKeyboardMarkup
}

/*EditByFileIdOrURL edits this animation file by file id or url*/
func (ai *AnimationEditor) EditByFileIdOrURL(fileIdOrUrl string) (*objs.DefaultResult, error) {
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
	return ai.mg.editMedia(ai.messageId, ai.inlineMessageId, im, ai.replyMarkup, nil, ai.thumbFile)
}

/*EditByFile edits this animation by file in the device*/
func (ai *AnimationEditor) EditByFile(file *os.File) (*objs.DefaultResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
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
	return ai.mg.editMedia(ai.messageId, ai.inlineMessageId, im, ai.replyMarkup, file, ai.thumbFile)
}

/*EditThumbnail edits the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (ai *AnimationInserter) EditThumbnail(fileIdOrURL string) {
	ai.thumb = fileIdOrURL
}

/*EditThumbnailFile edits the thumbnail of the file. It takes a file existing on the device*/
func (ai *AnimationInserter) EditThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	ai.thumbFile = file
	ai.thumb = "attach://" + stat.Name()
	return nil
}

//AudioEditor is a tool for editing audios.
type AudioEditor struct {
	mg                                                           *MessageEditor
	messageId                                                    int
	inlineMessageId, caption, parseMode, thumb, performer, title string
	captionEntities                                              []objs.MessageEntity
	thumbFile                                                    *os.File
	duration                                                     int
	replyMarkup                                                  *objs.InlineKeyboardMarkup
}

/*EditByFileIdOrURL edits this file by file id or url*/
func (ai *AudioEditor) EditByFileIdOrURL(fileIdOrUrl string) (*objs.DefaultResult, error) {
	im := &objs.InputMediaAudio{
		InputMediaDefault: fixTheDefault("audio", fileIdOrUrl, ai.caption, ai.parseMode, ai.captionEntities),
		Thumb:             ai.thumb,
		Performer:         ai.performer,
		Title:             ai.title,
	}
	if ai.duration != 0 {
		im.Duration = ai.duration
	}
	return ai.mg.editMedia(ai.messageId, ai.inlineMessageId, im, ai.replyMarkup, nil, ai.thumbFile)
}

/*EditByFile edits this audio by file in the device*/
func (ai *AudioEditor) EditByFile(file *os.File) (*objs.DefaultResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
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
	return ai.mg.editMedia(ai.messageId, ai.inlineMessageId, im, ai.replyMarkup, file, ai.thumbFile)
}

/*EditThumbnail edits the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (ai *AudioEditor) EditThumbnail(fileIdOrURL string) {
	ai.thumb = fileIdOrURL
}

/*EditThumbnailFile edits the thumbnail of the file. It takes a file existing on the device*/
func (ai *AudioEditor) EditThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	ai.thumbFile = file
	ai.thumb = "attach://" + stat.Name()
	return nil
}

//DocumentEditor is a tool for editing documents.
type DocumentEditor struct {
	mg                                         *MessageEditor
	messageId                                  int
	inlineMessageId, caption, parseMode, thumb string
	captionEntities                            []objs.MessageEntity
	thumbFile                                  *os.File
	disableContentTypeDetection                bool
	replyMarkup                                *objs.InlineKeyboardMarkup
}

/*EditByFileIdOrURL edits this file by file id or url*/
func (di *DocumentEditor) EditByFileIdOrURL(fileIdOrUrl string) (*objs.DefaultResult, error) {
	im := &objs.InputMediaDocument{
		InputMediaDefault:           fixTheDefault("document", fileIdOrUrl, di.caption, di.parseMode, di.captionEntities),
		Thumb:                       di.thumb,
		DisableContentTypeDetection: di.disableContentTypeDetection,
	}
	return di.mg.editMedia(di.messageId, di.inlineMessageId, im, di.replyMarkup, nil, di.thumbFile)
}

/*EditByFile edits this document by file in the device*/
func (di *DocumentEditor) EditByFile(file *os.File) (*objs.DefaultResult, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	im := &objs.InputMediaDocument{
		InputMediaDefault:           fixTheDefault("document", "attach://"+stat.Name(), di.caption, di.parseMode, di.captionEntities),
		Thumb:                       di.thumb,
		DisableContentTypeDetection: di.disableContentTypeDetection,
	}
	return di.mg.editMedia(di.messageId, di.inlineMessageId, im, di.replyMarkup, file, di.thumbFile)
}

/*EditThumbnail edits the tumbnail of the file. It takes a fileId or a url. If you want to send a file use "setThumbnailFile" instead.*/
func (di *DocumentEditor) EditThumbnail(fileIdOrURL string) {
	di.thumb = fileIdOrURL
}

/*EditThumbnailFile edits the thumbnail of the file. It takes a file existing on the device*/
func (di *DocumentEditor) EditThumbnailFile(file *os.File) error {
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	di.thumbFile = file
	di.thumb = "attach://" + stat.Name()
	return nil
}

/*EditText can be used to edit text and game messages. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.*/
func (me *MessageEditor) EditText(messageId int, text, inlineMessageId, parseMode string, entities []objs.MessageEntity, disableWebPagePreview bool, keyboard *inlineKeyboard) (*objs.DefaultResult, error) {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return me.bot.apiInterface.EditMessageText(
		me.chatIdInt, me.chatIdString, messageId, inlineMessageId, text,
		parseMode, entities, disableWebPagePreview, &replyMarkup,
	)
}

/*EditCaption can be used to edit captions of messages. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.*/
func (me *MessageEditor) EditCaption(messageId int, caption, inlineMessageId, parseMode string, captionEntities []objs.MessageEntity, keyboard *inlineKeyboard) (*objs.DefaultResult, error) {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return me.bot.apiInterface.EditMessageCaption(
		me.chatIdInt, me.chatIdString, messageId, inlineMessageId, caption,
		parseMode, captionEntities, &replyMarkup,
	)
}

/*EditMediaPhoto returns a PhotoEditor to edit a photo*/
func (mg *MessageEditor) EditMediaPhoto(messageId int, caption, parseMode string, captionEntitie []objs.MessageEntity, keyboard *inlineKeyboard) *PhotoEditor {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return &PhotoEditor{mg: mg, messageId: messageId, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, replyMarkup: &replyMarkup}
}

/*EditMediaVideo returns a VideoEditor to edit a video*/
func (mg *MessageEditor) EditMediaVideo(messageId int, caption, parseMode string, width, height, duration int, supportsStreaming bool, captionEntitie []objs.MessageEntity, keyboard *inlineKeyboard) *VideoEditor {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return &VideoEditor{mg: mg, messageId: messageId, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, width: width, height: height, duration: duration, supportsStreaming: supportsStreaming, replyMarkup: &replyMarkup}
}

/*EditMediaAnimation returns an AnimationEditor to edit an animation*/
func (mg *MessageEditor) EditMediaAnimation(messageId int, caption, parseMode string, width, height, duration int, captionEntitie []objs.MessageEntity, keyboard *inlineKeyboard) *AnimationEditor {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return &AnimationEditor{mg: mg, messageId: messageId, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, width: width, height: height, duration: duration, replyMarkup: &replyMarkup}
}

/*EditMediaAudio returns an AudioEditor to edit an audio*/
func (mg *MessageEditor) EditMediaAudio(messageId int, caption, parseMode, performer, title string, duration int, captionEntitie []objs.MessageEntity, keyboard *inlineKeyboard) *AudioEditor {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return &AudioEditor{mg: mg, messageId: messageId, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, performer: performer, title: title, duration: duration, replyMarkup: &replyMarkup}
}

/*EditMediaDocument returns a DocumentEditor to edit a document*/
func (mg *MessageEditor) EditMediaDocument(messageId int, caption, parseMode string, disableContentTypeDetection bool, captionEntitie []objs.MessageEntity, keyboard *inlineKeyboard) *DocumentEditor {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return &DocumentEditor{mg: mg, messageId: messageId, caption: caption, parseMode: parseMode, captionEntities: captionEntitie, disableContentTypeDetection: disableContentTypeDetection, replyMarkup: &replyMarkup}
}

/*EditReplyMarkup can be used to edit only the reply markup of messages. On success, if the edited message is not an inline message, the edited Message is returned, otherwise True is returned.*/
func (me *MessageEditor) EditReplyMarkup(messageId int, inlineMessageId string, keyboard *inlineKeyboard) (*objs.DefaultResult, error) {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return me.bot.apiInterface.EditMessagereplyMarkup(
		me.chatIdInt, me.chatIdString, messageId, inlineMessageId, &replyMarkup,
	)
}

/*DeleteMessage can be used to delete a message, including service messages, with the following limitations:

- A message can only be deleted if it was sent less than 48 hours ago.

- A dice message in a private chat can only be deleted if it was sent more than 24 hours ago.

- Bots can delete outgoing messages in private chats, groups, and supergroups.

- Bots can delete incoming messages in private chats.

- Bots granted can_post_messages permissions can delete outgoing messages in channels.

- If the bot is an administrator of a group, it can delete any message there.

- If the bot has can_delete_messages permission in a supergroup or a channel, it can delete any message there.

Returns True on success.*/
func (me *MessageEditor) DeleteMessage(messageId int) (*objs.LogicalResult, error) {
	return me.bot.apiInterface.DeleteMessage(me.chatIdInt, me.chatIdString, messageId)
}

func (me *MessageEditor) editMedia(messageId int, inlineMessageId string, media objs.InputMedia, replyMarkup *objs.InlineKeyboardMarkup, file ...*os.File) (*objs.DefaultResult, error) {
	return me.bot.apiInterface.EditMessageMedia(
		me.chatIdInt, me.chatIdString, messageId, inlineMessageId, media,
		replyMarkup, file...,
	)
}
