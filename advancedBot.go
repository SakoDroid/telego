package telego

import (
	"errors"
	"os"

	objs "github.com/SakoDroid/telego/objects"
)

/*AdvancedBot is an advanced type of bot which will give you alot more customization for the bot.
Methods which are uniquely for advanced bot start with 'A' .*/
type AdvancedBot struct {
	bot *Bot
}

/*ASendMessage sends a text message to a chat (not channel, use SendMessageUN method for sending messages to channles) and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.*/
func (bot *AdvancedBot) ASendMessage(chatId int, text, parseMode string, replyTo int, silent, protectContent bool, entites []objs.MessageEntity, disabelWebPagePreview, allowSendingWithoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendMessage(chatId, "", text, parseMode, entites, disabelWebPagePreview,
		silent, allowSendingWithoutReply, protectContent, replyTo, replyMarkup)
}

/*ASendMesssageUN sends a text message to a channel and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.*/
func (bot *AdvancedBot) ASendMesssageUN(chatId, text, parseMode string, replyTo int, silent, protectContent bool, entites []objs.MessageEntity, disabelWebPagePreview, allowSendingWithoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendMessage(0, chatId, text, parseMode, entites, disabelWebPagePreview,
		silent, allowSendingWithoutReply, protectContent, replyTo, replyMarkup)
}

/*ACopyMessage returns a MessageCopier which has several methods for copying a message*/
func (bot *AdvancedBot) ACopyMessage(messageId int, disableNotif bool, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, allowSendingWithoutReply bool, keyboard MarkUps) *MessageCopier {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MessageCopier{bot: bot.bot, messageId: messageId, disableNotif: disableNotif, caption: caption, parseMode: parseMode, captionEntities: captionEntites, allowSendingWihtouReply: allowSendingWithoutReply, replyTo: replyTo, replyMarkup: replyMarkup}
}

/*ASendPhoto returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to all types of chat except channels. To send a photo to a channel use "SendPhotoUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")*/
func (bot *AdvancedBot) ASendPhoto(chatId, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: PHOTO, bot: bot.bot, chatIdInt: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendPhotoUN returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")
*/
func (bot *AdvancedBot) ASendPhotoUN(chatId string, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: PHOTO, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendVideo returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to all types of chat except channels. To send a video to a channel use "SendVideoUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVideo(chatId int, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, duration int, supportsStreaming, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: VIDEO, bot: bot.bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, duration: duration, supportsStreaming: supportsStreaming, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendVideoUN returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVideoUN(chatId string, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, duration int, supportsStreaming, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: VIDEO, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, duration: duration, supportsStreaming: supportsStreaming, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendAudio returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to all types of chat except channels. To send a audio to a channel use "SendAudioUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.*/
func (bot *AdvancedBot) ASendAudio(chatId, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, performer, title string, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: AUDIO, bot: bot.bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, performer: performer, title: title, duration: duration, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendAudioUN returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.*/
func (bot *AdvancedBot) ASendAudioUN(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, performer, title string, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: AUDIO, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, performer: performer, title: title, duration: duration, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendDocument returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to all types of chat except channels. To send a audio to a channel use "SendDocumentUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendDocument(chatId, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, disableContentTypeDetection, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: DOCUMENT, bot: bot.bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, disableContentTypeDetection: disableContentTypeDetection, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendDocumentUN returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendDocumentUN(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, disableContentTypeDetection, allowSendingWithoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: DOCUMENT, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, disableContentTypeDetection: disableContentTypeDetection, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*ASendAnimation returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to all types of chat except channels. To send a audio to a channel use "SendAnimationUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendAnimation(chatId int, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, width, height, duration int, allowSendingWihtoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: ANIMATION, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, width: width, height: height, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*ASendAnimationUN returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to channels
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendAnimationUN(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, width, height, duration int, allowSendingWihtoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: ANIMATION, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, width: width, height: height, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*ASendVoice returns a MediaSender which has several methods for sending a voice. This method is only used for sending a voice to all types of chat except channels. To send a voice to a channel use "SendVoiceUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVoice(chatId int, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, allowSendingWihtoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: VOICE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*ASendVoiceUN returns a MediaSender which has several methods for sending a voice. This method is only used for sending a voice to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVoiceUN(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, allowSendingWihtoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: VOICE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*ASendVideoNote returns a MediaSender which has several methods for sending a video note. This method is only used for sending a video note to all types of chat except channels. To send a video note to a channel use "SendVideoNoteUN" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVideoNote(chatId int, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, length, duration int, allowSendingWihtoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup, length: length, duration: duration}
}

/*ASendVideoNoteUN returns an MediaSender which has several methods for sending a video note. This method is only used for sending a video note to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVideoNoteUN(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, length, duration int, allowSendingWihtoutReply bool, keyboard MarkUps) *MediaSender {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup, length: length, duration: duration}
}

/*ACreateAlbum creates a MediaGroup for grouping media messages.
To ignore replyTo argument, pass 0.*/
func (bot *AdvancedBot) ACreateAlbum(replyTo int, allowSendingWihtoutReply bool, keyboard MarkUps) *MediaGroup {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &MediaGroup{replyTo: replyTo, bot: bot.bot, media: make([]objs.InputMedia, 0), files: make([]*os.File, 0), allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*ASendVenue sends a venue to all types of chat but channels. To send it to channels use "SendVenueUN" method.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVenue(chatId, replyTo int, latitude, longitude float32, title, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType string, silent bool, allowSendingWihtoutReply, protectContent bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendVenue(
		chatId, "", latitude, longitude, title, address, foursquareId, foursquareType,
		googlePlaceId, googlePlaceType, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*ASendVenueUN sends a venue to a channel.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVenueUN(chatId string, replyTo int, latitude, longitude float32, title, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType string, silent, protectContent bool, allowSendingWihtoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendVenue(
		0, chatId, latitude, longitude, title, address, foursquareId, foursquareType,
		googlePlaceId, googlePlaceType, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*ASendContact sends a contact to all types of chat but channels. To send it to channels use "SendContactUN" method.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendContact(chatId, replyTo int, phoneNumber, firstName, lastName, vCard string, silent, protectContent bool, allowSendingWihtoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendContact(
		chatId, "", phoneNumber, firstName, lastName, vCard, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*ASendContactUN sends a contact to a channel.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendContactUN(chatId string, replyTo int, phoneNumber, firstName, lastName, vCard string, silent, protectContent bool, allowSendingWihtoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendContact(
		0, chatId, phoneNumber, firstName, lastName, vCard, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*ASendDice sends a dice message to all types of chat but channels. To send it to channels use "SendDiceUN" method.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned*/
func (bot *AdvancedBot) ASendDice(chatId, replyTo int, emoji string, silent, protectContent bool, allowSendingWihtoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendDice(
		chatId, "", emoji, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*ASendDiceUN sends a dice message to a channel.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned*/
func (bot *AdvancedBot) ASendDiceUN(chatId string, replyTo int, emoji string, silent, protectContent bool, allowSendingWihtoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendDice(
		0, chatId, emoji, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*ACreateLiveLocation creates a live location which has several methods for managing it.*/
func (bot *AdvancedBot) ACreateLiveLocation(latitude, longitude, accuracy float32, livePeriod, heading, proximtyAlertRadius, replyTo int, allowSendingWihtoutReply bool, keyboard MarkUps) *LiveLocation {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return &LiveLocation{bot: bot.bot, replyTo: replyTo, allowSendingWihoutReply: allowSendingWihtoutReply, latitude: latitude, longitude: longitude, livePeriod: livePeriod, horizontalAccuracy: accuracy, heading: heading, proximityAlertRadius: proximtyAlertRadius, replyMarkUp: replyMarkup}
}

/*ASendLocation sends a location (not live) to all types of chats but channels. To send it to channel use "SendLocationUN" method.

You can not use this methods to send a live location. To send a live location use "ACreateLiveLocation" method.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendLocation(chatId int, silent, protectContent bool, latitude, longitude, accuracy float32, replyTo int, allowSendingWihtoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendLocation(
		chatId, "", latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*ASendLocationUN sends a location (not live) to a channel.

You can not use this methods to send a live location. To send a live location use "ACreateLiveLocation" method.

If "silent" argument is true, the message will be sent without notification.

If "protectContent" argument is true, the message can't be forwarded or saved.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendLocationUN(chatId string, silent, protectContent bool, latitude, longitude, accuracy float32, replyTo int, allowSendingWihtoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendLocation(
		0, chatId, latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, allowSendingWihtoutReply, protectContent, replyMarkup,
	)
}

/*AAnswerCallbackQuery can be used to send answers to callback queries sent from inline keyboards. The answer will be displayed to the user as a notification at the top of the chat screen or as an alert. On success, True is returned.

Alternatively, the user can be redirected to the specified Game URL. For this option to work, you must first create a game for your bot via @Botfather and accept the terms. Otherwise, you may use links like t.me/your_bot?start=XXXX that open your bot with a parameter.*/
func (bot *AdvancedBot) AAnswerCallbackQuery(callbackQueryId, text string, showAlert bool, url string, cacheTime int) (*objs.LogicalResult, error) {
	return bot.bot.apiInterface.AnswerCallbackQuery(callbackQueryId, text, url, showAlert, cacheTime)
}

/*AAnswerInlineQuery returns an InlineQueryResponder which has several methods for answering an inline query.

--------------------------

Official telegram doc :

Use this method to send answers to an inline query. On success, True is returned.
No more than 50 results per query are allowed.*/
func (bot *AdvancedBot) AAnswerInlineQuery(id string, cacheTime int, isPersonal bool, nextOffset, switchPmText, switchPmParameter string) *InlineQueryResponder {
	return &InlineQueryResponder{bot: bot.bot, id: id, cacheTime: cacheTime, results: make([]objs.InlineQueryResult, 0), isPersonal: isPersonal, nextOffset: nextOffset, switchPmText: switchPmText, switchPmParameter: switchPmParameter}
}

/*ACreateInvoice returns an InvoiceSender which has several methods for creating and sending an invoice.

This method is suitable for sending this invoice to a chat that has an id, to send the invoice to channels use "ACreateInvoiceUN" method.*/
func (bot *AdvancedBot) ACreateInvoice(chatId int, title, description, payload, providerToken, currency string, prices []objs.LabeledPrice, maxTipAmount int, suggestedTipAmounts []int, startParameter, providerData, photoURL string, photoSize, photoWidth, photoHeight int, needName, needPhoneNumber, needEmail, needSippingAddress, sendPhoneNumberToProvider, sendEmailToProvider, isFlexible, bool, allowSendingWithoutReply bool, keyboard *inlineKeyboard) *Invoice {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return &Invoice{
		chatIdInt: chatId, chatIdString: "", title: title, description: description, providerToken: providerToken, currency: currency, prices: make([]objs.LabeledPrice, 0),
		bot: bot.bot, replyMarkup: replyMarkup, suggestedTipAmounts: suggestedTipAmounts, photoURL: photoURL, startParameter: startParameter, providerData: providerData, payload: payload,
		photoSize: photoSize, photoWidth: photoWidth, photoHeight: photoHeight, maxTipAmount: maxTipAmount, allowSendingWithoutReply: allowSendingWithoutReply, needName: needName, needPhoneNumber: needPhoneNumber,
		needEmail: needEmail, needShippingAddress: needSippingAddress, sendPhoneNumberToProvider: sendPhoneNumberToProvider, sendEmailToProvider: sendEmailToProvider, isFlexible: isFlexible,
	}
}

/*ACreateInvoiceUN returns an InvoiceSender which has several methods for creating and sending an invoice.*/
func (bot *AdvancedBot) ACreateInvoiceUN(chatId string, title, description, payload, providerToken, currency string, prices []objs.LabeledPrice, maxTipAmount int, suggestedTipAmounts []int, startParameter, providerData, photoURL string, photoSize, photoWidth, photoHeight int, needName, needPhoneNumber, needEmail, needSippingAddress, sendPhoneNumberToProvider, sendEmailToProvider, isFlexible, bool, allowSendingWithoutReply bool, keyboard *inlineKeyboard) *Invoice {
	var replyMarkup objs.InlineKeyboardMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toInlineKeyboardMarkup()
	}
	return &Invoice{
		chatIdInt: 0, chatIdString: chatId, title: title, description: description, providerToken: providerToken, currency: currency, prices: make([]objs.LabeledPrice, 0),
		bot: bot.bot, replyMarkup: replyMarkup, suggestedTipAmounts: suggestedTipAmounts, photoURL: photoURL, startParameter: startParameter, providerData: providerData, payload: payload,
		photoSize: photoSize, photoWidth: photoWidth, photoHeight: photoHeight, maxTipAmount: maxTipAmount, allowSendingWithoutReply: allowSendingWithoutReply, needName: needName, needPhoneNumber: needPhoneNumber,
		needEmail: needEmail, needShippingAddress: needSippingAddress, sendPhoneNumberToProvider: sendPhoneNumberToProvider, sendEmailToProvider: sendEmailToProvider, isFlexible: isFlexible,
	}
}

/*ASendGame sends a game to the chat.

-----------------------

Official telegram doc :

Use this method to send a game. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendGame(chatId int, gameShortName string, silent bool, replyTo int, allowSendingWithoutReply bool, keyboard MarkUps) (*objs.SendMethodsResult, error) {
	var replyMarkup objs.ReplyMarkup
	if keyboard != nil {
		replyMarkup = keyboard.toMarkUp()
	}
	return bot.bot.apiInterface.SendGame(
		chatId, gameShortName, silent, replyTo, allowSendingWithoutReply, replyMarkup,
	)
}

/*ASetGameScore sets the score of the given user.

-----------------------

Official telegram doc :

Use this method to set the score of the specified user in a game message. On success, if the message is not an inline message, the Message is returned, otherwise True is returned. Returns an error, if the new score is not greater than the user's current score in the chat and force is False.

"score" is new score, must be non-negative.

"force" : Pass True, if the high score is allowed to decrease. This can be useful when fixing mistakes or banning cheaters.

"disableEditMessage" : Pass True, if the game message should not be automatically edited to include the current scoreboard.

"inlineMessageId" : Required if chat_id and message_id are not specified. Identifier of the inline message.
*/
func (bot *AdvancedBot) ASetGameScore(userId, score, chatId, messageId int, force, disableEditMessage bool, inlineMessageId string) (*objs.DefaultResult, error) {
	return bot.bot.apiInterface.SetGameScore(
		userId, score, force, disableEditMessage, chatId, messageId, inlineMessageId,
	)
}

/*SetPassportDataErrors informs a user that some of the Telegram Passport elements they provided contains errors. The user will not be able to re-submit their Passport to you until the errors are fixed (the contents of the field for which you returned the error must change). Returns True on success.

Use this if the data submitted by the user doesn't satisfy the standards your service requires for any reason. For example, if a birthday date seems invalid, a submitted document is blurry, a scan shows evidence of tampering, etc. Supply some details in the error message to make sure the user knows how to correct the issues.*/
func (bot *AdvancedBot) SetPassportDataErrors(userId int, errors []objs.PassportElementError) (*objs.LogicalResult, error) {
	return bot.bot.apiInterface.SetPassportDataErrors(
		userId, errors,
	)
}

/*RegisterChannel can be used to register special channels. Sepcial channels can be used to get only certain types of updates through them. For example you can register a channel to only receive messages or register a channel to only receive edited messages in a certain chat.

"chatId" argument specifies a chat which this channel will be dedicated to. If an empty string is passed it means that no chat is specified and channel will work for all chats. chatIds that are integer should be converted to string.

"mediaType" argument specifies a media type that channel will be dedicated to. If an empty string is passed it means no media type is in mind and channel will work for all media types. mediaType can have the following values :

1. Empty string

2. message

3. edited_message

4. channel_post

5. edited_channel_post

6. inline_query

7. chosen_inline_query

8. callback_query

9. shipping_query

10. pre_checkout_query

11. poll_answer

12. my_chat_member

13. chat_member

14. chat_join_request

bot "chatId" and "mediaType" arguments can be used together to create a channel that will be updated only if a certain type of update is received for a ceratin chat.

Examples :

1. RegisterChannel("123456","message") : The returned channel will be updated when a message (text,photo,video ...) is received from a chat with "123456" as it's chat id.

2. RegiterChannel("","message") : The returned channel will be updated everytime a message is received from any chat.

3. RegisterChannel("123456","") : The returned channel will be updated everytime an update of anykind is received for the specified chat.

4. RegisterChannel("","") : The returned is the global update channel which will be updated everytime an update is received. You can get this channel by calling `getUpdateChannel()` method too.
*/
func (bot *AdvancedBot) RegisterChannel(chatId, mediaType string) (*chan *objs.Update, error) {
	if chatId == "" {
		chatId = "global"
	}
	if mediaType == "" {
		mediaType = "all"
	}
	if mediaType != "all" && mediaType != "message" && mediaType != "edited_message" && mediaType != "channel_post" && mediaType != "edited_channel_post" && mediaType != "inline_query" && mediaType != "chosen_inline_result" && mediaType != "callback_query" && mediaType != "shipping_query" && mediaType != "pre_checkout_query" && mediaType != "poll_answer" && mediaType != "my_chat_member" && mediaType != "chat_member" && mediaType != "chat_join_request" {
		return nil, errors.New("unknown media type : " + mediaType)
	}
	return bot.getChannel(chatId, mediaType), nil
}

/*UnRegisterChannel can be used to unregister a channel for the given arguments*/
func (bot *AdvancedBot) UnRegisterChannel(chatId, mediaType string) {
	if bot.bot.channelsMap[chatId] != nil {
		bot.bot.channelsMap[chatId][mediaType] = nil
		if len(bot.bot.channelsMap[chatId]) == 0 {
			delete(bot.bot.channelsMap, chatId)
		}
	}
}

func (bot *AdvancedBot) getChannel(chatId, media string) *chan *objs.Update {
	if bot.bot.channelsMap[chatId] == nil {
		bot.bot.channelsMap[chatId] = make(map[string]*chan *objs.Update)
	}
	out := bot.bot.channelsMap[chatId][media]
	if out == nil {
		tmp := make(chan *objs.Update)
		out = &tmp
		bot.bot.channelsMap[chatId][media] = out
	}
	return out
}
