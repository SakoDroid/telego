package telebot

import (
	"os"

	objs "github.com/SakoDroid/telebot/objects"
)

/*An advanced type of bot which will give you alot more customization for the bot.
Methods which are uniquely for advanced bot start with 'A' .*/
type AdvancedBot struct {
	bot *Bot
}

/*Send a text message to a chat (not channel, use SendMessageToChannel method for sending messages to channles) and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.*/
func (bot *AdvancedBot) ASendMessage(chatId int, text, parseMode string, replyTo int, silent bool, entites []objs.MessageEntity, disabelWebPagePreview, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendMessage(chatId, "", text, parseMode, entites, disabelWebPagePreview, silent, allowSendingWithoutReply, replyTo, replyMarkup)
}

/*Send a text message to a channel and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.*/
func (bot *AdvancedBot) ASendMesssageToChannel(chatId, text, parseMode string, replyTo int, silent bool, entites []objs.MessageEntity, disabelWebPagePreview, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendMessage(0, chatId, text, parseMode, entites, disabelWebPagePreview, silent, allowSendingWithoutReply, replyTo, replyMarkup)
}

/*Returns a MessageCopier which has several methods for copying a message*/
func (bot *AdvancedBot) ACopyMessage(messageId int, disableNotif bool, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MessageCopier {
	return &MessageCopier{bot: bot.bot, messageId: messageId, disableNotif: disableNotif, caption: caption, parseMode: parseMode, captionEntities: captionEntites, allowSendingWihtouReply: allowSendingWithoutReply, replyTo: replyTo, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to all types of chat except channels. To send a photo to a channel use "SendPhotoToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")*/
func (bot *AdvancedBot) ASendPhoto(chatId, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: PHOTO, bot: bot.bot, chatIdInt: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a photo. This method is only used for sending a photo to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")
*/
func (bot *AdvancedBot) ASendPhotoToChannel(chatId string, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: PHOTO, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to all types of chat except channels. To send a video to a channel use "SendVideoToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVideo(chatId int, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, duration int, supportsStreaming, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: VIDEO, bot: bot.bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, duration: duration, supportsStreaming: supportsStreaming, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a video. This method is only used for sending a video to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send video files, Telegram clients support mp4 videos (other formats may be sent as Document). On success, the sent Message is returned. Bots can currently send video files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVideoToChannel(chatId string, replyTo int, caption, parseMode string, captionEntites []objs.MessageEntity, duration int, supportsStreaming, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: VIDEO, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntites, duration: duration, supportsStreaming: supportsStreaming, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to all types of chat except channels. To send a audio to a channel use "SendAudioToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.*/
func (bot *AdvancedBot) ASendAudio(chatId, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, performer, title string, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: AUDIO, bot: bot.bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, performer: performer, title: title, duration: duration, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a audio. This method is only used for sending a audio to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display them in the music player. Your audio must be in the .MP3 or .M4A format. On success, the sent Message is returned. Bots can currently send audio files of up to 50 MB in size, this limit may be changed in the future.

For sending voice messages, use the sendVoice method instead.*/
func (bot *AdvancedBot) ASendAudioToChannel(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, performer, title string, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: AUDIO, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, performer: performer, title: title, duration: duration, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to all types of chat except channels. To send a audio to a channel use "SendDocumentToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendDocument(chatId, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, disableContentTypeDetection, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: DOCUMENT, bot: bot.bot, chatIdInt: chatId, chatidString: "", replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, disableContentTypeDetection: disableContentTypeDetection, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a document. This method is only used for sending a document to a channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send general files. On success, the sent Message is returned. Bots can currently send files of any type of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendDocumentToChannel(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, disableContentTypeDetection, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: DOCUMENT, bot: bot.bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode, captionEntities: captionEntities, disableContentTypeDetection: disableContentTypeDetection, allowSendingWihoutReply: allowSendingWithoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to all types of chat except channels. To send a audio to a channel use "SendAnimationToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendAnimation(chatId int, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, width, height, duration int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: ANIMATION, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, width: width, height: height, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending an animation. This method is only used for sending an animation to channels
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send animation files (GIF or H.264/MPEG-4 AVC video without sound). On success, the sent Message is returned. Bots can currently send animation files of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendAnimationToChannel(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, width, height, duration int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: ANIMATION, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, width: width, height: height, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a voice. This method is only used for sending a voice to all types of chat except channels. To send a voice to a channel use "SendVoiceToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVoice(chatId int, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: VOICE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a voice. This method is only used for sending a voice to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

Use this method to send audio files, if you want Telegram clients to display the file as a playable voice message. For this to work, your audio must be in an .OGG file encoded with OPUS (other formats may be sent as Audio or Document). On success, the sent Message is returned. Bots can currently send voice messages of up to 50 MB in size, this limit may be changed in the future.*/
func (bot *AdvancedBot) ASendVoiceToChannel(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, duration int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: VOICE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, duration: duration, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*Returns a MediaSender which has several methods for sending a video note. This method is only used for sending a video note to all types of chat except channels. To send a video note to a channel use "SendVideoNoteToChannel" method.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVideoNote(chatId int, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, length, duration int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: chatId, chatidString: "", replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup, length: length, duration: duration}
}

/*Returns an MediaSender which has several methods for sending a video note. This method is only used for sending a video note to channels.
To ignore int arguments pass 0 and to ignore string arguments pass empty string ("")

---------------------------------

Official telegram doc :

As of v.4.0, Telegram clients support rounded square mp4 videos of up to 1 minute long. Use this method to send video messages. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVideoNoteToChannel(chatId string, replyTo int, caption, parseMode string, captionEntities []objs.MessageEntity, length, duration int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *MediaSender {
	return &MediaSender{mediaType: VIDEONOTE, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, bot: bot.bot, caption: caption, parseMode: parseMode, captionEntities: captionEntities, allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup, length: length, duration: duration}
}

/*To ignore replyTo argument, pass 0.*/
func (bot *AdvancedBot) ACreateAlbum(replyTo int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *MediaGroup {
	return &MediaGroup{replyTo: replyTo, bot: bot.bot, media: make([]objs.InputMedia, 0), files: make([]*os.File, 0), allowSendingWihoutReply: allowSendingWihtoutReply, replyMarkup: replyMarkup}
}

/*Sends a venue to all types of chat but channels. To send it to channels use "SendVenueToChannel" method.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVenue(chatId, replyTo int, latitude, longitude float32, title, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType string, silent bool, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendVenue(
		chatId, "", latitude, longitude, title, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Sends a venue to a channel.

---------------------------------

Official telegram doc :

Use this method to send information about a venue. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendVenueTOChannel(chatId string, replyTo int, latitude, longitude float32, title, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType string, silent bool, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendVenue(
		0, chatId, latitude, longitude, title, address, foursquareId, foursquareType, googlePlaceId, googlePlaceType, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Sends a contact to all types of chat but channels. To send it to channels use "SendContactToChannel" method.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendContact(chatId, replyTo int, phoneNumber, firstName, lastName, vCard string, silent bool, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendContact(
		chatId, "", phoneNumber, firstName, lastName, vCard, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Sends a contact to a channel.

---------------------------------

Official telegram doc :

Use this method to send phone contacts. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendContactToChannel(chatId string, replyTo int, phoneNumber, firstName, lastName, vCard string, silent bool, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendContact(
		0, chatId, phoneNumber, firstName, lastName, vCard, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Sends a dice message to all types of chat but channels. To send it to channels use "SendDiceToChannel" method.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned*/
func (bot *AdvancedBot) ASendDice(chatId, replyTo int, emoji string, silent bool, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendDice(
		chatId, "", emoji, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Sends a dice message to a channel.

Available emojies : “🎲”, “🎯”, “🏀”, “⚽”, “🎳”, or “🎰”.

---------------------------------

Official telegram doc :

Use this method to send an animated emoji that will display a random value. On success, the sent Message is returned*/
func (bot *AdvancedBot) ASendDiceToChannel(chatId string, replyTo int, emoji string, silent bool, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendDice(
		0, chatId, emoji, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Creates a live location which has several methods for managing it.*/
func (bot *AdvancedBot) ACreateLiveLocation(latitude, longitude, accuracy float32, livePeriod, heading, proximtyAlertRadius, replyTo int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) *LiveLocation {
	return &LiveLocation{bot: bot.bot, replyTo: replyTo, allowSendingWihoutReply: allowSendingWihtoutReply, latitude: latitude, longitude: longitude, livePeriod: livePeriod, horizontalAccuracy: accuracy, heading: heading, proximityAlertRadius: proximtyAlertRadius, replyMarkUp: replyMarkup}
}

/*Sends a location (not live) to all types of chats but channels. To send it to channel use "SendLocationToChannel" method.

You can not use this methods to send a live location. To send a live location use "ACreateLiveLocation" method.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendLocation(chatId int, silent bool, latitude, longitude, accuracy float32, replyTo int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendLocation(
		chatId, "", latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Sends a location (not live) to a channel.

You can not use this methods to send a live location. To send a live location use "ACreateLiveLocation" method.

---------------------------------

Official telegram doc :

Use this method to send point on the map. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendLocationToChannel(chatId string, silent bool, latitude, longitude, accuracy float32, replyTo int, allowSendingWihtoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendLocation(
		0, chatId, latitude, longitude, accuracy, 0, 0, 0, replyTo, silent, allowSendingWihtoutReply, replyMarkup,
	)
}

/*Use this method to send answers to callback queries sent from inline keyboards. The answer will be displayed to the user as a notification at the top of the chat screen or as an alert. On success, True is returned.

Alternatively, the user can be redirected to the specified Game URL. For this option to work, you must first create a game for your bot via @Botfather and accept the terms. Otherwise, you may use links like t.me/your_bot?start=XXXX that open your bot with a parameter.*/
func (bot *AdvancedBot) AAnswerCallbackQuery(callbackQueryId, text string, showAlert bool, url string, cacheTime int) (*objs.LogicalResult, error) {
	return bot.bot.apiInterface.AnswerCallbackQuery(callbackQueryId, text, url, showAlert, cacheTime)
}

/*Returns an InlineQueryResponder which has several methods for answering an inline query.

--------------------------

Official telegram doc :

Use this method to send answers to an inline query. On success, True is returned.
No more than 50 results per query are allowed.*/
func (bot *AdvancedBot) AAnswerInlineQuery(id string, cacheTime int, isPersonal bool, nextOffset, switchPmText, switchPmParameter string) *InlineQueryResponder {
	return &InlineQueryResponder{bot: bot.bot, id: id, cacheTime: cacheTime, results: make([]objs.InlineQueryResult, 0), isPersonal: isPersonal, nextOffset: nextOffset, switchPmText: switchPmText, switchPmParameter: switchPmParameter}
}

/*Returnes an InvoiceSender which has several methods for creating and sending an invoice.

This method is suitable for sending this invoice to a chat that has an id, to send the invoice to channels use "ACreateInvoiceUN" method.*/
func (bot *AdvancedBot) ACreateInvoice(chatId int, title, description, payload, providerToken, currency string, prices []objs.LabeledPrice, maxTipAmount int, suggestedTipAmounts []int, startParameter, providerData, photoURL string, photoSize, photoWidth, photoHeight int, needName, needPhoneNumber, needEmail, needSippingAddress, sendPhoneNumberToProvider, sendEmailToProvider, isFlexible, bool, allowSendingWithoutReply bool, replyMarkup objs.InlineKeyboardMarkup) *InvoiceSender {
	return &InvoiceSender{
		chatIdInt: chatId, chatIdString: "", title: title, description: description, providerToken: providerToken, currency: currency, prices: make([]objs.LabeledPrice, 0),
		bot: bot.bot, replyMarkup: replyMarkup, suggestedTipAmounts: suggestedTipAmounts, photoURL: photoURL, startParameter: startParameter, providerData: providerData, payload: payload,
		photoSize: photoSize, photoWidth: photoWidth, photoHeight: photoHeight, maxTipAmount: maxTipAmount, allowSendingWithoutReply: allowSendingWithoutReply, needName: needName, needPhoneNumber: needPhoneNumber,
		needEmail: needEmail, needShippingAddress: needSippingAddress, sendPhoneNumberToProvider: sendPhoneNumberToProvider, sendEmailToProvider: sendEmailToProvider, isFlexible: isFlexible,
	}
}

/*Returnes an InvoiceSender which has several methods for creating and sending an invoice.*/
func (bot *AdvancedBot) ACreateInvoiceUN(chatId string, title, description, payload, providerToken, currency string, prices []objs.LabeledPrice, maxTipAmount int, suggestedTipAmounts []int, startParameter, providerData, photoURL string, photoSize, photoWidth, photoHeight int, needName, needPhoneNumber, needEmail, needSippingAddress, sendPhoneNumberToProvider, sendEmailToProvider, isFlexible, bool, allowSendingWithoutReply bool, replyMarkup objs.InlineKeyboardMarkup) *InvoiceSender {
	return &InvoiceSender{
		chatIdInt: 0, chatIdString: chatId, title: title, description: description, providerToken: providerToken, currency: currency, prices: make([]objs.LabeledPrice, 0),
		bot: bot.bot, replyMarkup: replyMarkup, suggestedTipAmounts: suggestedTipAmounts, photoURL: photoURL, startParameter: startParameter, providerData: providerData, payload: payload,
		photoSize: photoSize, photoWidth: photoWidth, photoHeight: photoHeight, maxTipAmount: maxTipAmount, allowSendingWithoutReply: allowSendingWithoutReply, needName: needName, needPhoneNumber: needPhoneNumber,
		needEmail: needEmail, needShippingAddress: needSippingAddress, sendPhoneNumberToProvider: sendPhoneNumberToProvider, sendEmailToProvider: sendEmailToProvider, isFlexible: isFlexible,
	}
}

/*Sends a game to the chat.

-----------------------

Official telegram doc :

Use this method to send a game. On success, the sent Message is returned.*/
func (bot *AdvancedBot) ASendGame(chatId int, gameShortName string, silent bool, replyTo int, allowSendingWithoutReply bool, replyMarkup objs.InlineKeyboardMarkup) (*objs.SendMethodsResult, error) {
	return bot.bot.apiInterface.SendGame(
		chatId, gameShortName, silent, replyTo, allowSendingWithoutReply, &replyMarkup,
	)
}

/*Sets the score of the given user.

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

/*
Informs a user that some of the Telegram Passport elements they provided contains errors. The user will not be able to re-submit their Passport to you until the errors are fixed (the contents of the field for which you returned the error must change). Returns True on success.

Use this if the data submitted by the user doesn't satisfy the standards your service requires for any reason. For example, if a birthday date seems invalid, a submitted document is blurry, a scan shows evidence of tampering, etc. Supply some details in the error message to make sure the user knows how to correct the issues.*/
func (bot *AdvancedBot) SetPassportDataErrors(userId int, errors []objs.PassportElementError) (*objs.LogicalResult, error) {
	return bot.bot.apiInterface.SetPassportDataErrors(
		userId, errors,
	)
}

/*Register a special channel for message updates only. Everytime an update is received that contains message filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterMessageChannel() *chan *objs.Message {
	if bot.bot.messageChannel == nil {
		ch := make(chan *objs.Message)
		bot.bot.messageChannel = &ch
	}
	return bot.bot.messageChannel
}

/*Register a special channel for edited message updates only. Everytime an update is received that contains edited_message filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterEditedMessageChannel() *chan *objs.Message {
	if bot.bot.editedMessageChannel == nil {
		ch := make(chan *objs.Message)
		bot.bot.editedMessageChannel = &ch
	}
	return bot.bot.editedMessageChannel
}

/*Register a special channel for channel post updates only. Everytime an update is received that contains channel_post filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterChannellPostChannel() *chan *objs.Message {
	if bot.bot.channelPostChannel == nil {
		ch := make(chan *objs.Message)
		bot.bot.channelPostChannel = &ch
	}
	return bot.bot.channelPostChannel
}

/*Register a special channel for edited channel post updates only. Everytime an update is received that contains edited_channel_post filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterEditedChannelPostChannel() *chan *objs.Message {
	if bot.bot.editedChannelPostChannel == nil {
		ch := make(chan *objs.Message)
		bot.bot.editedChannelPostChannel = &ch
	}
	return bot.bot.editedChannelPostChannel
}

/*Register a special channel for inline query updates only. Everytime an update is received that contains inline_query filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterInlineQueryChannel() *chan *objs.InlineQuery {
	if bot.bot.inlineQueryChannel == nil {
		ch := make(chan *objs.InlineQuery)
		bot.bot.inlineQueryChannel = &ch
	}
	return bot.bot.inlineQueryChannel
}

/*Register a special channel for chosen inline result updates only. Everytime an update is received that contains chosen_inline_result filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterChosenInlineResultChannel() *chan *objs.ChosenInlineResult {
	if bot.bot.chosenInlineResultChannel == nil {
		ch := make(chan *objs.ChosenInlineResult)
		bot.bot.chosenInlineResultChannel = &ch
	}
	return bot.bot.chosenInlineResultChannel
}

/*Register a special channel for callback query updates only. Everytime an update is received that contains callback_query filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterCallbackQueryChannel() *chan *objs.CallbackQuery {
	if bot.bot.callbackQueryChannel == nil {
		ch := make(chan *objs.CallbackQuery)
		bot.bot.callbackQueryChannel = &ch
	}
	return bot.bot.callbackQueryChannel
}

/*Register a special channel for shipping query updates only. Everytime an update is received that contains shipping_query filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterShippingQueryChannel() *chan *objs.ShippingQuery {
	if bot.bot.shippingQueryChannel == nil {
		ch := make(chan *objs.ShippingQuery)
		bot.bot.shippingQueryChannel = &ch
	}
	return bot.bot.shippingQueryChannel
}

/*Register a special channel for pre checkout query updates only. Everytime an update is received that contains pre_checkout_query filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterPreCheckoutQueryChannel() *chan *objs.PreCheckoutQuery {
	if bot.bot.preCheckoutQueryChannel == nil {
		ch := make(chan *objs.PreCheckoutQuery)
		bot.bot.preCheckoutQueryChannel = &ch
	}
	return bot.bot.preCheckoutQueryChannel
}

/*Register a special channel for my chat member updates only. Everytime an update is received that contains my_chat_member filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterMyChatMemberChannel() *chan *objs.ChatMemberUpdated {
	if bot.bot.myChatMemberChannel == nil {
		ch := make(chan *objs.ChatMemberUpdated)
		bot.bot.myChatMemberChannel = &ch
	}
	return bot.bot.myChatMemberChannel
}

/*Register a special channel for chat member updates only. Everytime an update is received that contains chat_member filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterChatMemberChannel() *chan *objs.ChatMemberUpdated {
	if bot.bot.chatMemberChannel == nil {
		ch := make(chan *objs.ChatMemberUpdated)
		bot.bot.chatMemberChannel = &ch
	}
	return bot.bot.chatMemberChannel
}

/*Register a special channel for chat join request updates only. Everytime an update is received that contains chat_join_request filed, the message field is passed into this channel. (The update won't be passed into general update channel anymore)

**Note : If a channe has been already set, this method will return it and won't set a new channe for this update type.*/
func (bot *AdvancedBot) RegisterChatJoinRequestChannel() *chan *objs.ChatJoinRequest {
	if bot.bot.chatJoinRequestChannel == nil {
		ch := make(chan *objs.ChatJoinRequest)
		bot.bot.chatJoinRequestChannel = &ch
	}
	return bot.bot.chatJoinRequestChannel
}
