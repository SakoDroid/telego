package telebot

import (
	"errors"
	"os"

	tba "github.com/SakoDroid/telebot/TBA"
	cfg "github.com/SakoDroid/telebot/configs"
	logger "github.com/SakoDroid/telebot/logger"
	objs "github.com/SakoDroid/telebot/objects"
)

type Bot struct {
	botCfg             *cfg.BotConfigs
	apiInterface       *tba.BotAPIInterface
	updateChannel      *chan *objs.Update
	pollUpdateChannel  *chan *objs.Update
	pollRoutineChannel *chan bool
}

/*Starts the bot. If the bot has already been started it returns an error.*/
func (bot *Bot) Run() error {
	logger.InitTheLogger(bot.botCfg)
	go bot.startPollUpdateRoutine()
	return bot.apiInterface.StartUpdateRoutine()
}

/*Returns the channel which new updates received from api server are pushed into.*/
func (bot *Bot) GetUpdateChannel() *chan *objs.Update {
	return bot.updateChannel
}

/*Returnes the received informations about the bot from api server.

---------------------

Official telegarm doc :

A simple method for testing your bot's authentication token. Requires no parameters. Returns basic information about the bot in form of a User object.*/
func (bot *Bot) GetMe() (*objs.UserResult, error) {
	return bot.apiInterface.GetMe()
}

/*Send a text message to a chat (not channel, use SendMessageToChannel method for sending messages to channles) and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.*/
func (bot *Bot) SendMessage(chatId int, text, parseMode string, replyTo int, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendMessage(chatId, "", text, parseMode, nil, false, silent, false, replyTo, nil)
}

/*Send a text message to a channel and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.*/
func (bot *Bot) SendMesssageToChannel(chatId, text, parseMode string, replyTo int, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendMessage(0, chatId, text, parseMode, nil, false, silent, false, replyTo, nil)
}

/*Returns a MessageForwarder which has several methods for forwarding a message*/
func (bot *Bot) ForwardMessage(messageId int, disableNotif bool) *MessageForwarder {
	return &MessageForwarder{bot: bot, messageId: messageId, disableNotif: disableNotif}
}

/*Returns a MessageCopier which has several methods for copying a message*/
func (bot *Bot) CopyMessage(messageId int, disableNotif bool) *MessageCopier {
	return &MessageCopier{bot: bot, messageId: messageId, disableNotif: disableNotif}
}

/*Returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to all types of chat except channels. To send a photo to a channel use "SendPhotoToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")*/
func (bot *Bot) SendPhoto(chatId, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: PHOTO, bot: bot, chatIdInt: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")
*/
func (bot *Bot) SendPhotoToChannel(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: PHOTO, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to all types of chat except channels. To send a video to a channel use "SendVideoToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendVideo(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEO, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendVideoToChannel(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEO, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to all types of chat except channels. To send a audio to a channel use "SendAudioToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.*/
func (bot *Bot) SendAudio(chatId, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: AUDIO, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.*/
func (bot *Bot) SendAudioToChannel(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: AUDIO, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to all types of chat except channels. To send a audio to a channel use "SendDocumentToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendDocument(chatId, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: DOCUMENT, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendDocumentToChannel(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: DOCUMENT, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to all types of chat except channels. To send a audio to a channel use "SendAnimationToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendAnimation(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: ANIMATION, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to channels
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendAnimationToChannel(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: ANIMATION, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a voice. This method is only used for sending a voice to all types of chat except channels. To send a voice to a channel use "SendVoiceToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendVoice(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VOICE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*Returns an MediaSender which has several methods for sending a voice. This method is only used for sending a voice to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *Bot) SendVoiceToChannel(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VOICE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*Returns a MediaSender which has several methods for sending a video note. This method is only used for sending a video note to all types of chat except channels. To send a video note to a channel use "SendVideoNoteToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.*/
func (bot *Bot) SendVideoNote(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*Returns an MediaSender which has several methods for sending a video note. This method is only used for sending a video note to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.*/
func (bot *Bot) SendVideoNoteToChannel(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*To ignore replyTo argument, pass 0.*/
func (bot *Bot) CreateAlbum(replyTo int) *MediaGroup {
	return &MediaGroup{replyTo: replyTo, bot: bot, media: make([]objs.InputMedia, 0), files: make([]*os.File, 0)}
}

/*Sends a venue to all types of chat but channels. To send it to channels use "SendVenueToChannel" method.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.*/
func (bot *Bot) SendVenue(chatId, replyTo int, latitude, longitude float32, title, address string, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendVenue(
		chatId, "", latitude, longitude, title, address, "", "", "", "", replyTo, silent, false, nil,
	)
}

/*Sends a venue to a channel.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.*/
func (bot *Bot) SendVenueTOChannel(chatId string, replyTo int, latitude, longitude float32, title, address string, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendVenue(
		0, chatId, latitude, longitude, title, address, "", "", "", "", replyTo, silent, false, nil,
	)
}

/*Sends a contact to all types of chat but channels. To send it to channels use "SendContactToChannel" method.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.*/
func (bot *Bot) SendContact(chatId, replyTo int, phoneNumber, firstName, lastName string, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendContact(
		chatId, "", phoneNumber, firstName, lastName, "", replyTo, silent, false, nil,
	)
}

/*Sends a contact to a channel.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.*/
func (bot *Bot) SendContactToChannel(chatId string, replyTo int, phoneNumber, firstName, lastName string, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendContact(
		0, chatId, phoneNumber, firstName, lastName, "", replyTo, silent, false, nil,
	)
}

/*Creates a poll for all types of chat but channels. To create a poll for channels use "CreatePollForChannel" method.

The poll type can be "regular" or "quiz"*/
func (bot *Bot) CreatePoll(chatId int, question, pollType string) (*Poll, error) {
	if pollType != "quiz" && pollType != "regular" {
		return nil, errors.New("poll type invalid : " + pollType)
	}
	return &Poll{bot: bot, pollType: pollType, chatIdInt: chatId, question: question, options: make([]string, 0)}, nil
}

/*Creates a poll for a channel.

The poll type can be "regular" or "quiz"*/
func (bot *Bot) CreatePollForChannel(chatId, question, pollType string) (*Poll, error) {
	if pollType != "quiz" && pollType != "regular" {
		return nil, errors.New("poll type invalid : " + pollType)
	}
	return &Poll{bot: bot, pollType: pollType, chatIdString: chatId, question: question, options: make([]string, 0)}, nil
}

/*Sends a dice message to all types of chat but channels. To send it to channels use "SendDiceToChannel" method.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned*/
func (bot *Bot) SendDice(chatId, replyTo int, emoji string, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendDice(
		chatId, "", emoji, replyTo, silent, false, nil,
	)
}

/*Sends a dice message to a channel.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned*/
func (bot *Bot) SendDiceToChannel(chatId string, replyTo int, emoji string, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendDice(
		0, chatId, emoji, replyTo, silent, false, nil,
	)
}

/*Sends a chat action message to all types of chat but channels. To send it to channels use "SendChatActionToChannel" method.

---------------------------------

Official telegram doc :

Use this method when you need to tell the user that something is happening on the bot's side. The status is set for 5 seconds or less (when a message arrives from your bot, Telegram clients clear its typing status). Returns True on success.

Example: The ImageBot needs some time to process a request and upload the image. Instead of sending a text message along the lines of “Retrieving image, please wait…”, the bot may use sendChatAction with action = upload_photo. The user will see a “sending photo” status for the bot.

We only recommend using this method when a response from the bot will take a noticeable amount of time to arrive.

action is the type of action to broadcast. Choose one, depending on what the user is about to receive: typing for text messages, upload_photo for photos, record_video or upload_video for videos, record_voice or upload_voice for voice notes, upload_document for general files, choose_sticker for stickers, find_location for location data, record_video_note or upload_video_note for video notes.*/
func (bot *Bot) SendChatAction(chatId int, action string) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendChatAction(chatId, "", action)
}

/*Sends a chat action message to a channel.

---------------------------------

Official telegram doc :

Use this method when you need to tell the user that something is happening on the bot's side. The status is set for 5 seconds or less (when a message arrives from your bot, Telegram clients clear its typing status). Returns True on success.

Example: The ImageBot needs some time to process a request and upload the image. Instead of sending a text message along the lines of “Retrieving image, please wait…”, the bot may use sendChatAction with action = upload_photo. The user will see a “sending photo” status for the bot.

We only recommend using this method when a response from the bot will take a noticeable amount of time to arrive.

action is the type of action to broadcast. Choose one, depending on what the user is about to receive: typing for text messages, upload_photo for photos, record_video or upload_video for videos, record_voice or upload_voice for voice notes, upload_document for general files, choose_sticker for stickers, find_location for location data, record_video_note or upload_video_note for video notes.*/
func (bot *Bot) SendChatActionToChannel(chatId, action string) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendChatAction(0, chatId, action)
}

/*Sends a location (not live) to all types of chats but channels. To send it to channel use "SendLocationToChannel" method.

You can not use this methods to send a live location. To send a live location use AdvancedBot.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (bot *Bot) SendLocation(chatId int, silent bool, latitude, longitude, accuracy float32, replyTo int) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendLocation(
		chatId, "", latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, false, nil,
	)
}

/*Sends a location (not live) to a channel.

You can not use this methods to send a live location. To send a live location use AdvancedBot.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (bot *Bot) SendLocationToChannel(chatId string, silent bool, latitude, longitude, accuracy float32, replyTo int) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendLocation(
		0, chatId, latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, false, nil,
	)
}

/*Gets the given user profile photos.

"userId" argument is required. Other arguments are optinoal and to ignore them pass 0.

---------------------------------

Official telegram doc :

Use this method to get a list of profile pictures for a user. Returns a UserProfilePhotos object.*/
func (bot *Bot) GetUserProfilePhotos(userId, offset, limit int) (*objs.ProfilePhototsResult, error) {
	return bot.apiInterface.GetUserProfilePhotos(userId, offset, limit)
}

/*Gets a file from telegram server. If it is successful the File object is returned.

If "download option is true, the file will be saved into the given file and if the given file is nil file will be saved in the same name as it has been saved in telegram servers.*/
func (bot *Bot) GetFile(fileId string, download bool, file *os.File) (*objs.File, error) {
	res, err := bot.apiInterface.GetFile(fileId)
	if err != nil {
		return nil, err
	}
	if download {
		err2 := bot.apiInterface.DownloadFile(res.Result, file)
		if err2 != nil {
			return res.Result, err2
		}
	}
	return res.Result, nil
}

/*Creates and returns a ChatManager for groups and other chats witch an integer id.

To manage supergroups and channels which have usernames use "GetChatManagerByUsername".*/
func (bot *Bot) GetChatManagerById(chatId int) *ChatManager {
	return &ChatManager{bot: bot, chatIdInt: chatId, chatIdString: ""}
}

/*Creates and returns a ChatManager for supergroups and channels which have usernames

To manage groups and other chats witch an integer id use "GetChatManagerById".*/
func (bot *Bot) GetChatManagerByUsrename(chatId int) *ChatManager {
	return &ChatManager{bot: bot, chatIdInt: chatId, chatIdString: ""}
}

/*Use this method to send answers to callback queries sent from inline keyboards. The answer will be displayed to the user as a notification at the top of the chat screen or as an alert. On success, True is returned.

Alternatively, the user can be redirected to the specified Game URL. For this option to work, you must first create a game for your bot via @Botfather and accept the terms. Otherwise, you may use links like t.me/your_bot?start=XXXX that open your bot with a parameter.*/
func (bot *Bot) AnswerCallbackQuery(callbackQueryId, text string, showAlert bool) (*objs.LogicalResult, error) {
	return bot.apiInterface.AnswerCallbackQuery(callbackQueryId, text, "", showAlert, 0)
}

/*Returns a command manager which has several method for manaing bot commands.*/
func (bot *Bot) GetCommandManager() *CommandsManager {
	return &CommandsManager{bot: bot}
}

/*Returns a MessageEditor for a chat with id which has several methods for editing messages.

To edit messages in a channel or a chat with username, use "GetMsgEditorWithUN"*/
func (bot *Bot) GetMsgEditor(chatId int) *MessageEditor {
	return &MessageEditor{bot: bot, chatIdInt: chatId}
}

/*Returns a MessageEditor for a chat with username which has several methods for editing messages.*/
func (bot *Bot) GetMsgEditorWithUN(chatId string) *MessageEditor {
	return &MessageEditor{bot: bot, chatIdInt: 0, chatIdString: chatId}
}

/*Returns a MediaSender which has several methods for sending an sticker to all types of chats but channels.
To send it to a channel use "SendStickerWithUN".

--------------------

Official telegram doc :


Use this method to send static .WEBP or animated .TGS stickers. On success, the sent Message is returned*/
func (bot *Bot) SendSticker(chatId, replyTo int) *MediaSender {
	return &MediaSender{mediaType: STICKER, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo}
}

/*Returns a MediaSender which has several methods for sending an sticker to channels.

--------------------

Official telegram doc :


Use this method to send static .WEBP or animated .TGS stickers. On success, the sent Message is returned*/
func (bot *Bot) SendStickerWithUn(chatId string, replyTo int) *MediaSender {
	return &MediaSender{mediaType: STICKER, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo}
}

/*Returns an sticker set with the given name*/
func (bot *Bot) GetStickerSet(name string) (*StickerSet, error) {
	res, err := bot.apiInterface.GetStickerSet(name)
	if err != nil {
		return nil, err
	}
	return &StickerSet{bot: bot, stickerSet: res.Result}, nil
}

/*
Use this method to upload a .PNG file with a sticker for later use in CreateNewStickerSet and AddStickerToSet methods (can be used multiple times). Returns the uploaded File on success.*/
func (bot *Bot) UploadStickerFile(userId int, stickerFile *os.File) (*objs.GetFileResult, error) {
	stat, err := stickerFile.Stat()
	if err != nil {
		return nil, err
	}
	return bot.apiInterface.UploadStickerFile(userId, "attach://"+stat.Name(), stickerFile)
}

/*
Use this method to create a new sticker set owned by a user. The bot will be able to edit the sticker set thus created. You must use exactly one of the fields pngSticker or tgsSticker. Returns the created sticker set on success.

png sticker can be passed as an file id or url (pngStickerFileIdOrUrl) or file(pngStickerFile).

"name" is the short name of sticker set, to be used in t.me/addstickers/ URLs (e.g., animals). Can contain only english letters, digits and underscores. Must begin with a letter, can't contain consecutive underscores and must end in “_by_<bot username>”. <bot_username> is case insensitive. 1-64 characters.*/
func (bot *Bot) CreateNewStickerSet(userId int, name, title, pngStickerFileIdOrUrl string, pngStickerFile *os.File, tgsSticker *os.File, emojies string, containsMask bool, maskPosition *objs.MaskPosition) (*StickerSet, error) {
	var res *objs.LogicalResult
	var err error
	if tgsSticker == nil {
		if pngStickerFile == nil {
			res, err = bot.apiInterface.CreateNewStickerSet(
				userId, name, title, pngStickerFileIdOrUrl, "", emojies, containsMask, maskPosition, nil,
			)
		} else {
			stat, er := pngStickerFile.Stat()
			if er != nil {
				return nil, er
			}
			res, err = bot.apiInterface.CreateNewStickerSet(
				userId, name, title, "attach://"+stat.Name(), "", emojies, containsMask, maskPosition, pngStickerFile,
			)
		}
	} else {
		stat, er := tgsSticker.Stat()
		if er != nil {
			return nil, er
		}
		res, err = bot.apiInterface.CreateNewStickerSet(
			userId, name, title, "", "attach://"+stat.Name(), emojies, containsMask, maskPosition, tgsSticker,
		)
	}
	if err != nil {
		return nil, err
	}
	if !res.Result {
		return nil, errors.New("false returned from server")
	}
	out := &StickerSet{bot: bot, stickerSet: &objs.StickerSet{
		Name: name, Title: title, ContainsMask: containsMask, Stickers: make([]objs.Sticker, 0),
	}}
	out.update()
	return out, nil
}

/*Returns an InlineQueryResponder which has several methods for answering an inline query.
To access more options use "AAsnwerInlineQuery" method in advanced bot.

--------------------------

Official telegram doc :

Use this method to send answers to an inline query. On success, True is returned.
No more than 50 results per query are allowed.*/
func (bot *Bot) AnswerInlineQuery(id string, cacheTime int) *InlineQueryResponder {
	return &InlineQueryResponder{bot: bot, id: id, cacheTime: cacheTime, results: make([]objs.InlineQueryResult, 0)}
}

/*Returnes an InvoiceSender which has several methods for creating and sending an invoice.

This method is suitable for sending this invoice to a chat that has an id, to send the invoice to channels use "CreateInvoiceUN" method.

To access more options, use "ACreateInvoice" method in advanced mode.*/
func (bot *Bot) CreateInvoice(chatId int, title, description, payload, providerToken, currency string) *InvoiceSender {
	return &InvoiceSender{
		bot: bot, chatIdInt: chatId, chatIdString: "", title: title, description: description, providerToken: providerToken, payload: payload, currency: currency, prices: make([]objs.LabeledPrice, 0),
	}
}

/*Returnes an InvoiceSender which has several methods for creating and sending an invoice.

To access more options, use "ACreateInvoiceUN" method in advanced mode.*/
func (bot *Bot) CreateInvoiceUN(chatId, title, description, payload, providerToken, currency string) *InvoiceSender {
	return &InvoiceSender{
		chatIdInt: 0, chatIdString: chatId, title: title, description: description, providerToken: providerToken, currency: currency, payload: payload, prices: make([]objs.LabeledPrice, 0),
	}
}

/*Answers an incoming shipping query.

-----------------------

Official telegram doc :

If you sent an invoice requesting a shipping address and the parameter is_flexible was specified, the Bot API will send an Update with a shipping_query field to the bot. Use this method to reply to shipping queries. On success, True is returned.

"ok" : Specify True if delivery to the specified address is possible and False if there are any problems (for example, if delivery to the specified address is not possible).

"shippingOptions" : Required if ok is True. A JSON-serialized array of available shipping options.

"errorMessage" : Required if ok is False. Error message in human readable form that explains why it is impossible to complete the order (e.g. "Sorry, delivery to your desired address is unavailable'). Telegram will display this message to the user.*/
func (bot *Bot) AnswerShippingQuery(shippingQueryId string, ok bool, shippingOptions []objs.ShippingOption, errorMessage string) (*objs.LogicalResult, error) {
	return bot.apiInterface.AnswerShippingQuery(shippingQueryId, ok, shippingOptions, errorMessage)
}

/*Answers a pre checkout query.

-----------------------

Official telegram doc :

Once the user has confirmed their payment and shipping details, the Bot API sends the final confirmation in the form of an Update with the field pre_checkout_query. Use this method to respond to such pre-checkout queries. On success, True is returned. Note: The Bot API must receive an answer within 10 seconds after the pre-checkout query was sent.

"ok" : Specify True if everything is alright (goods are available, etc.) and the bot is ready to proceed with the order. Use False if there are any problems.

"errorMessage" : Required if ok is False. Error message in human readable form that explains the reason for failure to proceed with the checkout (e.g. "Sorry, somebody just bought the last of our amazing black T-shirts while you were busy filling out your payment details. Please choose a different color or garment!"). Telegram will display this message to the user.*/
func (bot *Bot) AnswerPreCheckoutQuery(shippingQueryId string, ok bool, errorMessage string) (*objs.LogicalResult, error) {
	return bot.apiInterface.AnswerPreCheckoutQuery(shippingQueryId, ok, errorMessage)
}

/*Sends a game to the chat.

**To access more options use "ASendGame" method in advanced mode.

-----------------------

Official telegram doc :

Use this method to send a game. On success, the sent Message is returned.*/
func (bot *Bot) SendGame(chatId int, gameShortName string, silent bool, replyTo int) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendGame(
		chatId, gameShortName, silent, replyTo, false, nil,
	)
}

/*Sets the score of the given user.


**To access more option use "ASetGameScoe" in advanced mode.
-----------------------

Official telegram doc :

Use this method to set the score of the specified user in a game message. On success, if the message is not an inline message, the Message is returned, otherwise True is returned. Returns an error, if the new score is not greater than the user's current score in the chat and force is False.

"score" is new score, must be non-negative.
*/
func (bot *Bot) SetGameScore(userId, score, chatId, messageId int) (*objs.DefaultResult, error) {
	return bot.apiInterface.SetGameScore(
		userId, score, false, false, chatId, messageId, "",
	)
}

/*Returnes the high scores of the user.

-------------------------

Official telegram doc :


Use this method to get data for high score tables. Will return the score of the specified user and several of their neighbors in a game. On success, returns an Array of GameHighScore objects.

This method will currently return scores for the target user, plus two of their closest neighbors on each side. Will also return the top three users if the user and his neighbors are not among them. Please note that this behavior is subject to change.

"chatId" : Required if inline_message_id is not specified. Unique identifier for the target chat.

"messageId" : Required if inline_message_id is not specified. Identifier of the sent message.

"inlineMessageId" : Required if chat_id and message_id are not specified. Identifier of the inline message.*/
func (bot *Bot) GetGameHighScores(userId, chatId, messageId int, inlineMessageId string) (*objs.GameHighScoresResult, error) {
	return bot.apiInterface.GetGameHighScores(userId, chatId, messageId, inlineMessageId)
}

/*Stops the bot*/
func (bot *Bot) Stop() {
	bot.apiInterface.StopUpdateRoutine()
	*bot.pollRoutineChannel <- true
}

/*Returns and advanced version which gives more customized functions to iteract with the bot*/
func (bot *Bot) AdvancedMode() *AdvancedBot {
	return &AdvancedBot{bot: bot}
}

func (bot *Bot) startPollUpdateRoutine() {
loop:
	for {
		select {
		case <-*bot.pollRoutineChannel:
			break loop
		default:
			poll := <-*bot.pollUpdateChannel
			id := poll.Poll.Id
			pl := Polls[id]
			if pl == nil {
				logger.Logger.Println("Could not update poll `" + id + "`. Not found in the Polls map")
				*bot.updateChannel <- poll
				continue
			}
			err3 := pl.Update(poll.Poll)
			if err3 != nil {
				logger.Logger.Println("Could not update poll `" + id + "`." + err3.Error())
			}
		}
	}
}

/*Return a new bot instance with the specified configs*/
func NewBot(cfg *cfg.BotConfigs) (*Bot, error) {
	api, err := tba.CreateInterface(cfg)
	if err != nil {
		return nil, err
	}
	ch := make(chan bool)
	return &Bot{botCfg: cfg, apiInterface: api, updateChannel: api.GetUpdateChannel(), pollUpdateChannel: api.GetPollUpdateChannel(), pollRoutineChannel: &ch}, nil
}
