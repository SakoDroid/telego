package telego

import (
	"errors"
	"os"

	cfg "github.com/SakoDroid/telego/configs"
	logger "github.com/SakoDroid/telego/logger"
	objs "github.com/SakoDroid/telego/objects"
	upp "github.com/SakoDroid/telego/parser"
	tba "github.com/SakoDroid/telego/tba"
)

type Bot struct {
	botCfg                 *cfg.BotConfigs
	apiInterface           *tba.BotAPIInterface
	channelsMap            map[string]map[string]*chan *objs.Update
	interfaceUpdateChannel *chan *objs.Update
	chatUpdateChannel      *chan *objs.ChatUpdate
	prcRoutineChannel      *chan bool
	ab                     *AdvancedBot
}

/*Run starts the bot. If the bot has already been started it returns an error.*/
func (bot *Bot) Run(pause bool) error {
	logger.InitTheLogger(bot.botCfg)
	if !bot.checkWebHook() {
		logger.Logger.Fatalln("Webhook check failed. See the logs for more info.")
	}
	go bot.startChatUpdateRoutine()
	go bot.startUpdateProcessing()
	cfg.Dump(bot.botCfg)
	go bot.botCfg.StartCfgUpdateRoutine()
	var err error
	if bot.botCfg.Webhook {
		err = tba.StartWebHook(bot.botCfg, bot.interfaceUpdateChannel, bot.chatUpdateChannel)
	} else {
		err = bot.apiInterface.StartUpdateRoutine()
	}
	if err != nil {
		return err
	}
	ch := make(chan bool)
	<-ch
	return nil
}

func (bot *Bot) checkWebHook() bool {
	wi, err := bot.apiInterface.GetWebhookInfo()
	if err != nil {
		logger.Logger.Println(err)
		return false
	}
	if wi.Result.URL == "" {
		if bot.botCfg.Webhook {
			err2 := bot.setWebhook()
			if err2 != nil {
				logger.Logger.Println("Unable to set a new webhook.", err2)
				return false
			}
		}
	} else {
		if bot.botCfg.Webhook {
			if wi.Result.URL != bot.botCfg.WebHookConfigs.URL {
				logger.Logger.Println("A webhook is already set in the API server to this url :", wi.Result.URL, ". Deleting the webhook ...")
				err2 := bot.deleteWebhook()
				if err2 != nil {
					logger.Logger.Println("Unable to delete webhook.", err2)
					return false
				}
				err3 := bot.setWebhook()
				if err3 != nil {
					logger.Logger.Println("Unable to set webhook.", err3)
					return false
				}
			}
		} else {
			logger.Logger.Println("A webhook has been set.")
			err2 := bot.deleteWebhook()
			if err2 != nil {
				logger.Logger.Println("Unable to delete webhook.", err2)
				return false
			}
		}
	}
	return true
}

/*Sets a new webhook*/
func (bot *Bot) setWebhook() error {
	logger.Logger.Println("Setting webhook ...")
	whcfg := bot.botCfg.WebHookConfigs
	var fl *os.File
	if whcfg.SelfSigned {
		var err2 error
		fl, err2 = os.Open(whcfg.CertFile)
		if err2 != nil {
			return err2
		}
	}
	res, err3 := bot.apiInterface.SetWebhook(whcfg.URL, whcfg.IP, whcfg.MaxConnections, whcfg.AllowedUpdates, whcfg.DropPendingUpdates, fl)
	if err3 != nil {
		return err3
	}
	if !res.Result {
		return errors.New("unable to set webhook. API server returned false")
	}
	return nil
}

/*Deletes a webhook*/
func (bot *Bot) deleteWebhook() error {
	logger.Logger.Println("Deleting the webhook ...")
	res, err2 := bot.apiInterface.DeleteWebhook(false)
	if err2 != nil {
		return err2
	}
	if !res.Result {
		return errors.New("failed to delete the webhook. API server returned false")
	}
	return nil
}

// BlockUser blocks a user based on their ID and username.
func (bot *Bot) BlockUser(user *objs.User) {
	for _, us := range bot.botCfg.BlockedUsers {
		if us.UserID == user.Id {
			return
		}
	}
	us := cfg.BlockedUser{UserID: user.Id, UserName: user.Username}
	bot.botCfg.BlockedUsers = append(bot.botCfg.BlockedUsers, us)
}

/*GetUpdateChannel returns the channel which new updates received from api server are pushed into.*/
func (bot *Bot) GetUpdateChannel() *chan *objs.Update {
	return bot.channelsMap["global"]["all"]
}

/*
AddHandler adds a handler for a text message that matches the given regex pattern and chatType.

"pattern" is a regex pattern.

"chatType" must be "private","group","supergroup","channel" or "all". Any other value will cause the function to return an error.
*/
func (bot *Bot) AddHandler(pattern string, handler func(*objs.Update), chatTypes ...string) error {
	for _, val := range chatTypes {
		if val != "private" && val != "group" && val != "supergroup" && val != "channel" && val != "all" {
			return errors.New("unknown chat type : " + val)
		}
	}
	return upp.AddHandler(pattern, handler, chatTypes...)

}

/*
GetMe returns the received informations about the bot from api server.

---------------------

Official telegarm doc :

A simple method for testing your bot's authentication token. Requires no parameters. Returns basic information about the bot in form of a User object.
*/
func (bot *Bot) GetMe() (*objs.UserResult, error) {
	return bot.apiInterface.GetMe()
}

/*
SendMessage sens a text message to a chat (not channel, use SendMessageUN method for sending messages to channles) and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendMessage(chatId int, text, parseMode string, replyTo int, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendMessage(chatId, "", text, parseMode, nil, false, silent, false, protectContent, replyTo, nil)
}

/*
SendMesssageUN sens a text message to a channel and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendMessageUN(chatId, text, parseMode string, replyTo int, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendMessage(0, chatId, text, parseMode, nil, false, silent, false, protectContent, replyTo, nil)
}

func (bot *Bot) PinChatMessage(chatIdInt int, chatIdString string, messageId int, disableNotification bool) (*objs.LogicalResult, error) {
	return bot.apiInterface.PinChatMessage(chatIdInt, chatIdString, messageId, disableNotification)
}

func (bot *Bot) UnpinChatMessage(chatIdInt int, chatIdString string, messageId int) (*objs.LogicalResult, error) {
	return bot.apiInterface.UnpinChatMessage(chatIdInt, chatIdString, messageId)
}

func (bot *Bot) UnpinAllChatMessages(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	return bot.apiInterface.UnpinAllChatMessages(chatIdInt, chatIdString)
}

func (bot *Bot) CreateChatInviteLink(chatIdInt int, chatIdString, name string, expireDate, memberLimit int, createsJoinRequest bool) (*objs.ChatInviteLinkResult, error) {
	return bot.apiInterface.CreateChatInviteLink(chatIdInt, chatIdString, name, expireDate, memberLimit, createsJoinRequest)
}

func (bot *Bot) GetChatMember(chatIdInt int, chatIdString string, userId int) (*objs.DefaultResult, error) {
	return bot.apiInterface.GetChatMember(chatIdInt, chatIdString, userId)
}

func (bot *Bot) BanChatMember(chatIdInt int, chatIdString string, userId, untilDate int, revokeMessages bool) (*objs.LogicalResult, error) {
	return bot.apiInterface.BanChatMember(chatIdInt, chatIdString, userId, untilDate, revokeMessages)
}

func (bot *Bot) UnbanChatMember(chatIdInt int, chatIdString string, userId int, onlyIfBanned bool) (*objs.LogicalResult, error) {
	return bot.apiInterface.UnbanChatMember(chatIdInt, chatIdString, userId, onlyIfBanned)
}

func (bot *Bot) SetMyCommands(commands []objs.BotCommand, scope objs.BotCommandScope, languageCode string) (*objs.LogicalResult, error) {
	return bot.apiInterface.SetMyCommands(commands, scope, languageCode)
}

/*
ForwardMessage returns a MessageForwarder which has several methods for forwarding a message

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) ForwardMessage(messageId int, disableNotif, protectContent bool) *MessageForwarder {
	return &MessageForwarder{bot: bot, messageId: messageId, disableNotif: disableNotif, protectContent: protectContent}
}

/*
CopyMessage returns a MessageCopier which has several methods for copying a message

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) CopyMessage(messageId int, disableNotif, protectContent bool) *MessageCopier {
	return &MessageCopier{bot: bot, messageId: messageId, disableNotif: disableNotif, protectContent: protectContent}
}

/*
SendPhoto returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to all types of chat except channels. To send a photo to a channel use "SendPhotoUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")
*/
func (bot *Bot) SendPhoto(chatId, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: PHOTO, bot: bot, chatIdInt: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendPhotoUN returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")
*/
func (bot *Bot) SendPhotoUN(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: PHOTO, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendVideo returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to all types of chat except channels. To send a video to a channel use "SendVideoUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendVideo(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEO, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendVideoUN returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendVideoUN(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEO, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendAudio returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to all types of chat except channels. To send a audio to a channel use "SendAudioUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.
*/
func (bot *Bot) SendAudio(chatId, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: AUDIO, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendAudioUN returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.
*/
func (bot *Bot) SendAudioUN(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: AUDIO, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendDocument returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to all types of chat except channels. To send a audio to a channel use "SendDocumentUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendDocument(chatId, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: DOCUMENT, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendDocumentUN returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendDocumentUN(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: DOCUMENT, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

/*
SendAnimation returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to all types of chat except channels. To send a audio to a channel use "SendAnimationUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendAnimation(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: ANIMATION, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*
SendAnimationUN returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to channels
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendAnimationUN(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: ANIMATION, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*
SendVoice returns a MediaSender which has several methods for sending a voice. This method is only used for sending a voice to all types of chat except channels. To send a voice to a channel use "SendVoiceUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendVoice(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VOICE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*
SendVoiceUN returns an MediaSender which has several methods for sending a voice. This method is only used for sending a voice to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.
*/
func (bot *Bot) SendVoiceUN(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VOICE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*
SendVideoNote returns a MediaSender which has several methods for sending a video note. This method is only used for sending a video note to all types of chat except channels. To send a video note to a channel use "SendVideoNoteUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.
*/
func (bot *Bot) SendVideoNote(chatId int, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*
SendVideoNoteUN returns an MediaSender which has several methods for sending a video note. This method is only used for sending a video note to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.
*/
func (bot *Bot) SendVideoNoteUN(chatId string, replyTo int, caption, parseMode string) *MediaSender {
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot, caption: caption, parseMode: parseMode}
}

/*
CreateAlbum returns a MediaGroup for grouping media messages.
To ignore replyTo argument, pass 0.
*/
func (bot *Bot) CreateAlbum(replyTo int) *MediaGroup {
	return &MediaGroup{replyTo: replyTo, bot: bot, media: make([]objs.InputMedia, 0), files: make([]*os.File, 0)}
}

/*
SendVenue sends a venue to all types of chat but channels. To send it to channels use "SendVenueUN" method.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendVenue(chatId, replyTo int, latitude, longitude float32, title, address string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendVenue(
		chatId, "", latitude, longitude, title, address, "", "", "", "", replyTo, silent, false, protectContent, nil,
	)
}

/*
SendVenueUN sends a venue to a channel.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendVenueUN(chatId string, replyTo int, latitude, longitude float32, title, address string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendVenue(
		0, chatId, latitude, longitude, title, address, "", "", "", "", replyTo, silent, false, protectContent, nil,
	)
}

/*
SendContact sends a contact to all types of chat but channels. To send it to channels use "SendContactUN" method.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendContact(chatId, replyTo int, phoneNumber, firstName, lastName string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendContact(
		chatId, "", phoneNumber, firstName, lastName, "", replyTo, silent, false, protectContent, nil,
	)
}

/*
SendContactUN sends a contact to a channel.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendContactUN(chatId string, replyTo int, phoneNumber, firstName, lastName string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendContact(
		0, chatId, phoneNumber, firstName, lastName, "", replyTo, silent, false, protectContent, nil,
	)
}

/*
CreatePoll creates a poll for all types of chat but channels. To create a poll for channels use "CreatePollForChannel" method.

The poll type can be "regular" or "quiz"
*/
func (bot *Bot) CreatePoll(chatId int, question, pollType string) (*Poll, error) {
	if pollType != "quiz" && pollType != "regular" {
		return nil, errors.New("poll type invalid : " + pollType)
	}
	return &Poll{bot: bot, pollType: pollType, chatIdInt: chatId, question: question, options: make([]string, 0)}, nil
}

/*
CreatePollForChannel creates a poll for a channel.

The poll type can be "regular" or "quiz"
*/
func (bot *Bot) CreatePollForChannel(chatId, question, pollType string) (*Poll, error) {
	if pollType != "quiz" && pollType != "regular" {
		return nil, errors.New("poll type invalid : " + pollType)
	}
	return &Poll{bot: bot, pollType: pollType, chatIdString: chatId, question: question, options: make([]string, 0)}, nil
}

/*
SendDice sends a dice message to all types of chat but channels. To send it to channels use "SendDiceUN" method.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendDice(chatId, replyTo int, emoji string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendDice(
		chatId, "", emoji, replyTo, silent, false, protectContent, nil,
	)
}

/*
SendDiceUN sends a dice message to a channel.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendDiceUN(chatId string, replyTo int, emoji string, silent, protectContent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendDice(
		0, chatId, emoji, replyTo, silent, false, protectContent, nil,
	)
}

/*
SendChatAction sends a chat action message to all types of chat but channels. To send it to channels use "SendChatActionUN" method.

---------------------------------

Official telegram doc :

Use this method when you need to tell the user that something is happening on the bot's side. The status is set for 5 seconds or less (when a message arrives from your bot, Telegram clients clear its typing status). Returns True on success.

Example: The ImageBot needs some time to process a request and upload the image. Instead of sending a text message along the lines of “Retrieving image, please wait…”, the bot may use sendChatAction with action = upload_photo. The user will see a “sending photo” status for the bot.

We only recommend using this method when a response from the bot will take a noticeable amount of time to arrive.

action is the type of action to broadcast. Choose one, depending on what the user is about to receive: typing for text messages, upload_photo for photos, record_video or upload_video for videos, record_voice or upload_voice for voice notes, upload_document for general files, choose_sticker for stickers, find_location for location data, record_video_note or upload_video_note for video notes.
*/
func (bot *Bot) SendChatAction(chatId int, action string) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendChatAction(chatId, "", action)
}

/*
SendChatActionUN sends a chat action message to a channel.

---------------------------------

Official telegram doc :

Use this method when you need to tell the user that something is happening on the bot's side. The status is set for 5 seconds or less (when a message arrives from your bot, Telegram clients clear its typing status). Returns True on success.

Example: The ImageBot needs some time to process a request and upload the image. Instead of sending a text message along the lines of “Retrieving image, please wait…”, the bot may use sendChatAction with action = upload_photo. The user will see a “sending photo” status for the bot.

We only recommend using this method when a response from the bot will take a noticeable amount of time to arrive.

action is the type of action to broadcast. Choose one, depending on what the user is about to receive: typing for text messages, upload_photo for photos, record_video or upload_video for videos, record_voice or upload_voice for voice notes, upload_document for general files, choose_sticker for stickers, find_location for location data, record_video_note or upload_video_note for video notes.
*/
func (bot *Bot) SendChatActionUN(chatId, action string) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendChatAction(0, chatId, action)
}

/*
SendLocation sends a location (not live) to all types of chats but channels. To send it to channel use "SendLocationUN" method.

You can not use this methods to send a live location. To send a live location use AdvancedBot.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendLocation(chatId int, silent, protectContent bool, latitude, longitude, accuracy float32, replyTo int) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendLocation(
		chatId, "", latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, false, protectContent, nil,
	)
}

/*
SendLocationUN sends a location (not live) to a channel.

You can not use this methods to send a live location. To send a live location use AdvancedBot.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.
*/
func (bot *Bot) SendLocationUN(chatId string, silent, protectContent bool, latitude, longitude, accuracy float32, replyTo int) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendLocation(
		0, chatId, latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, false, protectContent, nil,
	)
}

/*
GetUserProfilePhotos gets the given user profile photos.

"userId" argument is required. Other arguments are optinoal and to ignore them pass 0.

---------------------------------

Official telegram doc :

Use this method to get a list of profile pictures for a user. Returns a UserProfilePhotos object.
*/
func (bot *Bot) GetUserProfilePhotos(userId, offset, limit int) (*objs.ProfilePhototsResult, error) {
	return bot.apiInterface.GetUserProfilePhotos(userId, offset, limit)
}

/*
GetFile gets a file from telegram server. If it is successful the File object is returned.

If "download option is true, the file will be saved into the given file and if the given file is nil file will be saved in the same name as it has been saved in telegram servers.
*/
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

/*
GetChatManagerById creates and returns a ChatManager for groups and other chats witch an integer id.

To manage supergroups and channels which have usernames use "GetChatManagerByUsername".
*/
func (bot *Bot) GetChatManagerById(chatId int) *ChatManager {
	return &ChatManager{bot: bot, chatIdInt: chatId, chatIdString: ""}
}

/*
GetChatManagerByUsrename creates and returns a ChatManager for supergroups and channels which have usernames

To manage groups and other chats witch an integer id use "GetChatManagerById".
*/
func (bot *Bot) GetChatManagerByUsrename(userName string) *ChatManager {
	return &ChatManager{bot: bot, chatIdInt: 0, chatIdString: userName}
}

/*
AnswerCallbackQuery can be used to send answers to callback queries sent from inline keyboards. The answer will be displayed to the user as a notification at the top of the chat screen or as an alert. On success, True is returned.

Alternatively, the user can be redirected to the specified Game URL. For this option to work, you must first create a game for your bot via @Botfather and accept the terms. Otherwise, you may use links like t.me/your_bot?start=XXXX that open your bot with a parameter.
*/
func (bot *Bot) AnswerCallbackQuery(callbackQueryId, text string, showAlert bool) (*objs.LogicalResult, error) {
	return bot.apiInterface.AnswerCallbackQuery(callbackQueryId, text, "", showAlert, 0)
}

/*GetCommandManager returns a command manager which has several method for manaing bot commands.*/
func (bot *Bot) GetCommandManager() *CommandsManager {
	return &CommandsManager{bot: bot}
}

/*
GetMsgEditor returns a MessageEditor for a chat with id which has several methods for editing messages.

To edit messages in a channel or a chat with username, use "GetMsgEditorWithUN"
*/
func (bot *Bot) GetMsgEditor(chatId int) *MessageEditor {
	return &MessageEditor{bot: bot, chatIdInt: chatId}
}

/*GetMsgEditorWithUN returns a MessageEditor for a chat with username which has several methods for editing messages.*/
func (bot *Bot) GetMsgEditorWithUN(chatId string) *MessageEditor {
	return &MessageEditor{bot: bot, chatIdInt: 0, chatIdString: chatId}
}

/*
SendSticker returns a MediaSender which has several methods for sending an sticker to all types of chats but channels.
To send it to a channel use "SendStickerWithUN".

--------------------

Official telegram doc :

Use this method to send static .WEBP or animated .TGS stickers. On success, the sent Message is returned
*/
func (bot *Bot) SendSticker(chatId, replyTo int) *MediaSender {
	return &MediaSender{mediaType: STICKER, bot: bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo}
}

/*
SendStickerWithUn returns a MediaSender which has several methods for sending an sticker to channels.

--------------------

Official telegram doc :

Use this method to send static .WEBP or animated .TGS stickers. On success, the sent Message is returned
*/
func (bot *Bot) SendStickerWithUn(chatId string, replyTo int) *MediaSender {
	return &MediaSender{mediaType: STICKER, bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo}
}

/*GetStickerSet returns an sticker set with the given name*/
func (bot *Bot) GetStickerSet(name string) (*StickerSet, error) {
	res, err := bot.apiInterface.GetStickerSet(name)
	if err != nil {
		return nil, err
	}
	return &StickerSet{bot: bot, stickerSet: res.Result}, nil
}

/*UploadStickerFile can be used to upload a .PNG file with a sticker for later use in CreateNewStickerSet and AddStickerToSet methods (can be used multiple times). Returns the uploaded File on success.*/
func (bot *Bot) UploadStickerFile(userId int, stickerFile *os.File) (*objs.GetFileResult, error) {
	stat, err := stickerFile.Stat()
	if err != nil {
		return nil, err
	}
	return bot.apiInterface.UploadStickerFile(userId, "attach://"+stat.Name(), stickerFile)
}

/*
CreateNewStickerSet can be used to create a new sticker set owned by a user. The bot will be able to edit the sticker set thus created. You must use exactly one of the fields pngSticker or tgsSticker or webmSticker. Returns the created sticker set on success.

png sticker can be passed as an file id or url (pngStickerFileIdOrUrl) or file(pngStickerFile).

"name" is the short name of sticker set, to be used in t.me/addstickers/ URLs (e.g., animals). Can contain only english letters, digits and underscores. Must begin with a letter, can't contain consecutive underscores and must end in “_by_<bot username>”. <bot_username> is case insensitive. 1-64 characters.
*/
func (bot *Bot) CreateNewStickerSet(userId int, name, title, pngStickerFileIdOrUrl string, pngStickerFile *os.File, tgsSticker *os.File, webmSticker *os.File, emojies string, containsMask bool, maskPosition *objs.MaskPosition) (*StickerSet, error) {
	var res *objs.LogicalResult
	var err error
	if tgsSticker == nil {
		if pngStickerFile == nil {
			if pngStickerFileIdOrUrl == "" {
				if webmSticker == nil {
					return nil, errors.New("wrong file id or url")
				} else {
					stat, er := webmSticker.Stat()
					if er != nil {
						return nil, er
					}
					res, err = bot.apiInterface.CreateNewStickerSet(
						userId, name, title, "", "", "attach://"+stat.Name(), emojies, containsMask, maskPosition, pngStickerFile,
					)
				}
			}
			res, err = bot.apiInterface.CreateNewStickerSet(
				userId, name, title, pngStickerFileIdOrUrl, "", "", emojies, containsMask, maskPosition, nil,
			)
		} else {
			stat, er := pngStickerFile.Stat()
			if er != nil {
				return nil, er
			}
			res, err = bot.apiInterface.CreateNewStickerSet(
				userId, name, title, "attach://"+stat.Name(), "", "", emojies, containsMask, maskPosition, pngStickerFile,
			)
		}
	} else {
		stat, er := tgsSticker.Stat()
		if er != nil {
			return nil, er
		}
		res, err = bot.apiInterface.CreateNewStickerSet(
			userId, name, title, "", "attach://"+stat.Name(), "", emojies, containsMask, maskPosition, tgsSticker,
		)
	}
	if err != nil {
		return nil, err
	}
	if !res.Result {
		return nil, errors.New("false returned from server")
	}
	out := &StickerSet{bot: bot, userId: userId, stickerSet: &objs.StickerSet{
		Name: name, Title: title, ContainsMask: containsMask, Stickers: make([]objs.Sticker, 0),
	}}
	out.update()
	return out, nil
}

/*
AnswerInlineQuery returns an InlineQueryResponder which has several methods for answering an inline query or web app query.
To access more options use "AAsnwerInlineQuery" method in advanced bot.

--------------------------

Official telegram doc :

Use this method to send answers to an inline query. On success, True is returned.
No more than 50 results per query are allowed.
*/
func (bot *Bot) AnswerInlineQuery(id string, cacheTime int) *InlineQueryResponder {
	return &InlineQueryResponder{bot: bot, id: id, cacheTime: cacheTime, results: make([]objs.InlineQueryResult, 0)}
}

/*
AnswerWebAppQuery returns an InlineQueryResponder which has several methods for answering an inline query or web app query.

--------------------------

Official telegram doc :

Use this method to set the result of an interaction with a Web App and send a corresponding message on behalf of the user to the chat from which the query originated. On success, a SentWebAppMessage object is returned.
*/
func (bot *Bot) AnswerWebAppQuery(webAppQueryId string) *InlineQueryResponder {
	return &InlineQueryResponder{id: webAppQueryId, isWebApp: true, results: make([]objs.InlineQueryResult, 0)}
}

/*
CreateInvoice returns an InvoiceSender which has several methods for creating and sending an invoice.

This method is suitable for sending this invoice to a chat that has an id, to send the invoice to channels use "CreateInvoiceUN" method.

To access more options, use "ACreateInvoice" method in advanced mode.
*/
func (bot *Bot) CreateInvoice(chatId int, title, description, payload, providerToken, currency string) *Invoice {
	return &Invoice{
		bot: bot, chatIdInt: chatId, chatIdString: "", title: title, description: description, providerToken: providerToken, payload: payload, currency: currency, prices: make([]objs.LabeledPrice, 0),
	}
}

/*
CreateInvoiceUN returns an InvoiceSender which has several methods for creating and sending an invoice.

To access more options, use "ACreateInvoiceUN" method in advanced mode.
*/
func (bot *Bot) CreateInvoiceUN(chatId, title, description, payload, providerToken, currency string) *Invoice {
	return &Invoice{
		chatIdInt: 0, chatIdString: chatId, title: title, description: description, providerToken: providerToken, currency: currency, payload: payload, prices: make([]objs.LabeledPrice, 0),
	}
}

/*
AnswerShippingQuery answers an incoming shipping query.

-----------------------

Official telegram doc :

If you sent an invoice requesting a shipping address and the parameter is_flexible was specified, the Bot API will send an Update with a shipping_query field to the bot. Use this method to reply to shipping queries. On success, True is returned.

"ok" : Specify True if delivery to the specified address is possible and False if there are any problems (for example, if delivery to the specified address is not possible).

"shippingOptions" : Required if ok is True. A JSON-serialized array of available shipping options.

"errorMessage" : Required if ok is False. Error message in human readable form that explains why it is impossible to complete the order (e.g. "Sorry, delivery to your desired address is unavailable'). Telegram will display this message to the user.
*/
func (bot *Bot) AnswerShippingQuery(shippingQueryId string, ok bool, shippingOptions []objs.ShippingOption, errorMessage string) (*objs.LogicalResult, error) {
	return bot.apiInterface.AnswerShippingQuery(shippingQueryId, ok, shippingOptions, errorMessage)
}

/*
AnswerPreCheckoutQuery answers a pre checkout query.

-----------------------

Official telegram doc :

Once the user has confirmed their payment and shipping details, the Bot API sends the final confirmation in the form of an Update with the field pre_checkout_query. Use this method to respond to such pre-checkout queries. On success, True is returned. Note: The Bot API must receive an answer within 10 seconds after the pre-checkout query was sent.

"ok" : Specify True if everything is alright (goods are available, etc.) and the bot is ready to proceed with the order. Use False if there are any problems.

"errorMessage" : Required if ok is False. Error message in human readable form that explains the reason for failure to proceed with the checkout (e.g. "Sorry, somebody just bought the last of our amazing black T-shirts while you were busy filling out your payment details. Please choose a different color or garment!"). Telegram will display this message to the user.
*/
func (bot *Bot) AnswerPreCheckoutQuery(shippingQueryId string, ok bool, errorMessage string) (*objs.LogicalResult, error) {
	return bot.apiInterface.AnswerPreCheckoutQuery(shippingQueryId, ok, errorMessage)
}

/*
SendGame sends a game to the chat.

**To access more options use "ASendGame" method in advanced mode.

-----------------------

Official telegram doc :

Use this method to send a game. On success, the sent Message is returned.
*/
func (bot *Bot) SendGame(chatId int, gameShortName string, silent bool, replyTo int) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendGame(
		chatId, gameShortName, silent, replyTo, false, nil,
	)
}

/*
SetGameScore sets the score of the given user.

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

/*
GetGameHighScores returns the high scores of the user.

-------------------------

Official telegram doc :

Use this method to get data for high score tables. Will return the score of the specified user and several of their neighbors in a game. On success, returns an Array of GameHighScore objects.

This method will currently return scores for the target user, plus two of their closest neighbors on each side. Will also return the top three users if the user and his neighbors are not among them. Please note that this behavior is subject to change.

"chatId" : Required if inline_message_id is not specified. Unique identifier for the target chat.

"messageId" : Required if inline_message_id is not specified. Identifier of the sent message.

"inlineMessageId" : Required if chat_id and message_id are not specified. Identifier of the inline message.
*/
func (bot *Bot) GetGameHighScores(userId, chatId, messageId int, inlineMessageId string) (*objs.GameHighScoresResult, error) {
	return bot.apiInterface.GetGameHighScores(userId, chatId, messageId, inlineMessageId)
}

/*
GetChatMenuButton gets the current menu button of given chat.

-------------------------

Official telegram doc :

Use this method to get the current value of the bot's menu button in a private chat, or the default menu button. Returns MenuButton on success.
*/
func (bot *Bot) GetChatMenuButton(chatId int64) (*objs.MenuButtonResult, error) {
	return bot.apiInterface.GetChatMenuButton(chatId)
}

/*
SetCommandChatMenuButton sets the current menu button of given chat to command meaning that it opens the bot's list of commands.

-------------------------

Official telegram doc :

Use this method to change the bot's menu button in a private chat, or the default menu button. Returns True on success.
*/
func (bot *Bot) SetCommandChatMenuButton(chatId int64) (*objs.LogicalResult, error) {
	return bot.apiInterface.SetChatMenuButton(chatId, &objs.MenuButton{Type: "commands"})
}

/*
SetDefaultChatMenuButton sets the current menu button of given chat to command meaning that it describes that no specific value for the menu button was set.

-------------------------

Official telegram doc :

Use this method to change the bot's menu button in a private chat, or the default menu button. Returns True on success.
*/
func (bot *Bot) SetDefaultChatMenuButton(chatId int64) (*objs.LogicalResult, error) {
	return bot.apiInterface.SetChatMenuButton(chatId, &objs.MenuButton{Type: "default"})
}

/*
SetWebAppChatMenuButton sets the current menu button of given chat to web_app meaning that it launches a Web App.

-------------------------

Official telegram doc :

Use this method to change the bot's menu button in a private chat, or the default menu button. Returns True on success.
*/
func (bot *Bot) SetWebAppChatMenuButton(chatId int64, text, url string) (*objs.LogicalResult, error) {
	return bot.apiInterface.SetChatMenuButton(chatId, &objs.MenuButton{Type: "web_app", Text: text, WebApp: &objs.WebAppInfo{URL: url}})
}

/*
CreateKeyboard creates a keyboard an returns it. The created keyboard has some methods for adding buttons to it.

You can send the keyboard along with messages by passing the keyboard as the "keyboard" argument of the method. The methods that supoort keyboard are mostly located in the advanced mode.

Arguments (as described in telegram bot api):

1. resizeKeyboard : Requests clients to resize the keyboard vertically for optimal fit (e.g., make the keyboard smaller if there are just two rows of buttons). Defaults to false, in which case the custom keyboard is always of the same height as the app's standard keyboard.

2. oneTimeKeyboard : Requests clients to hide the keyboard as soon as it's been used. The keyboard will still be available, but clients will automatically display the usual letter-keyboard in the chat – the user can press a special button in the input field to see the custom keyboard again. Defaults to false

3. inputFieldPlaceholder : The placeholder to be shown in the input field when the keyboard is active; 1-64 characters.

4. selective : Use this parameter if you want to show the keyboard to specific users only. Targets: 1) users that are @mentioned in the text of the Message object; 2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.

Example: A user requests to change the bot's language, bot replies to the request with a keyboard to select the new language. Other users in the group don't see the keyboard.
*/
func (bot *Bot) CreateKeyboard(resizeKeyboard, oneTimeKeyboard, selective bool, inputFieldPlaceholder string) *keyboard {
	return &keyboard{
		keys:                  make([][]*objs.KeyboardButton, 0),
		resizeKeyBoard:        resizeKeyboard,
		oneTimeKeyboard:       oneTimeKeyboard,
		selective:             selective,
		inputFieldPlaceHolder: inputFieldPlaceholder,
	}
}

/*
CreateInlineKeyboard creates a keyboard an returns it. The created keyboard has some methods for adding buttons to it.

You can send the keyboard along with messages by passing the keyboard as the "keyboard" argument of a method. The methods that supoort keyboard are mostly located in the advanced mode.
*/
func (bot *Bot) CreateInlineKeyboard() *inlineKeyboard {
	return &inlineKeyboard{}
}

/*GetTextFormatter returns a MessageFormatter that can be used for formatting a text message. You can add bold,italic,underline,spoiler,mention,url,link and some other texts with this tool.*/
func (bot *Bot) GetTextFormatter() *TextFormatter {
	return &TextFormatter{entites: make([]objs.MessageEntity, 0)}
}

/*VerifyJoin verifies if the user has joined the given channel or supergroup. Returns true if the user is present in the given chat, returns false if not or an error has occured.*/
func (bot *Bot) VerifyJoin(userID int, UserName string) bool {
	_, err := bot.apiInterface.GetChatMember(0, UserName, userID)
	return err == nil
}

/*Stop stops the bot*/
func (bot *Bot) Stop() {
	bot.apiInterface.StopUpdateRoutine()
	*bot.prcRoutineChannel <- true
}

/*AdvancedMode returns and advanced version of the bot which gives more customized functions to iteract with the bot*/
func (bot *Bot) AdvancedMode() *AdvancedBot {
	return bot.ab
}

func (bot *Bot) processUpdate(update *objs.Update, mapKey string) bool {
	out := true
	upType := update.GetType()
	if upType == "poll" {
		bot.processPoll(update)
	} else {
		ch := bot.channelsMap[mapKey][upType]
		if ch != nil {
			*ch <- update
		} else {
			out = false
		}
	}
	return out
}

func (bot *Bot) startUpdateProcessing() {
loop:
	for {
		select {
		case <-*bot.prcRoutineChannel:
			break loop
		case up := <-*bot.interfaceUpdateChannel:
			if !bot.processUpdate(up, "global") {
				*bot.channelsMap["global"]["all"] <- up
			}
		}
	}
}

func (bot *Bot) processPoll(update *objs.Update) {
	id := update.Poll.Id
	pl := Polls[id]
	if pl == nil {
		logger.Log("Error", "\t\t\t", "Could not update poll `"+id+"`. Not found in the Polls map", "917", logger.BOLD+logger.FAIL, logger.WARNING, "")
		*bot.channelsMap["global"]["all"] <- update
	} else {
		err3 := pl.Update(update.Poll)
		if err3 != nil {
			logger.Log("Error", "\t\t\t", "Could not update poll `"+id+"`."+err3.Error(), "922", logger.BOLD+logger.FAIL, logger.WARNING, "")
		}
	}
}

func (bot *Bot) startChatUpdateRoutine() {
loop:
	for {
		select {
		case <-*bot.prcRoutineChannel:
			break loop
		case up := <-*bot.chatUpdateChannel:
			if !bot.processUpdate(up.Update, up.ChatId) {
				chatChannel := bot.channelsMap[up.ChatId]["all"]
				if chatChannel != nil {
					*chatChannel <- up.Update
				} else {
					*bot.interfaceUpdateChannel <- up.Update
				}
			}
		}
	}
}

/*NewBot returns a new bot instance with the specified configs*/
func NewBot(cfg *cfg.BotConfigs) (*Bot, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}
	if !cfg.Check() {
		return nil, errors.New("config check failed. Please check the configs")
	}
	api, err := tba.CreateInterface(cfg)
	if err != nil {
		return nil, err
	}
	ch := make(chan bool)
	uc := make(chan *objs.Update)
	bt := &Bot{botCfg: cfg, apiInterface: api, interfaceUpdateChannel: api.GetUpdateChannel(), chatUpdateChannel: api.GetChatUpdateChannel(), prcRoutineChannel: &ch, channelsMap: make(map[string]map[string]*chan *objs.Update)}
	bt.channelsMap["global"] = make(map[string]*chan *objs.Update)
	bt.channelsMap["global"]["all"] = &uc
	bt.ab = &AdvancedBot{bot: bt}
	return bt, nil
}
