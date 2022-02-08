package tba

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	cfgs "github.com/SakoDroid/telego/configs"
	errs "github.com/SakoDroid/telego/errors"
	logger "github.com/SakoDroid/telego/logger"
	objs "github.com/SakoDroid/telego/objects"
	up "github.com/SakoDroid/telego/parser"
)

var interfaceCreated = false

//BotAPIInterface is the interface which connects the telegram bot API to the bot.
type BotAPIInterface struct {
	botConfigs           *cfgs.BotConfigs
	updateRoutineRunning bool
	updateChannel        *chan *objs.Update
	chatUpadateChannel   *chan *objs.ChatUpdate
	updateRoutineChannel chan bool
	lastOffset           int
}

/*StartUpdateRoutine starts the update routine to receive updates from api sever*/
func (bai *BotAPIInterface) StartUpdateRoutine() error {
	if !bai.botConfigs.Webhook {
		if bai.updateRoutineRunning {
			return &errs.UpdateRoutineAlreadyStarted{}
		}
		bai.updateRoutineRunning = true
		bai.updateRoutineChannel = make(chan bool)
		go bai.startReceiving()
		return nil
	} else {
		return errors.New("webhook option is true")
	}
}

/*StopUpdateRoutine stops the update routine*/
func (bai *BotAPIInterface) StopUpdateRoutine() {
	if bai.updateRoutineRunning {
		bai.updateRoutineRunning = false
		bai.updateRoutineChannel <- true
	}
}

/*GetUpdateChannel returns the update channel*/
func (bai *BotAPIInterface) GetUpdateChannel() *chan *objs.Update {
	return bai.updateChannel
}

/*GetChatUpdateChannel returnes the chat update channel*/
func (bai *BotAPIInterface) GetChatUpdateChannel() *chan *objs.ChatUpdate {
	return bai.chatUpadateChannel
}

func (bai *BotAPIInterface) startReceiving() {
	cl := httpSenderClient{botApi: bai.botConfigs.BotAPI, apiKey: bai.botConfigs.APIKey}
loop:
	for {
		time.Sleep(bai.botConfigs.UpdateConfigs.UpdateFrequency)
		select {
		case <-bai.updateRoutineChannel:
			break loop
		default:
			args := objs.GetUpdatesArgs{Offset: bai.lastOffset + 1, Limit: bai.botConfigs.UpdateConfigs.Limit, Timeout: bai.botConfigs.UpdateConfigs.Timeout}
			if bai.botConfigs.UpdateConfigs.AllowedUpdates != nil {
				args.AllowedUpdates = bai.botConfigs.UpdateConfigs.AllowedUpdates
			}
			res, err := cl.sendHttpReqJson("getUpdates", &args)
			if err != nil {
				logger.Logger.Println("Error receiving updates.", err)
				continue loop
			}
			err = bai.parseUpdateresults(res)
			if err != nil {
				logger.Logger.Println("Error parsing the result of the update. " + err.Error())
			}
		}
	}
}

func (bai *BotAPIInterface) parseUpdateresults(body []byte) error {
	of, err := up.ParseUpdate(
		body, bai.updateChannel, bai.chatUpadateChannel, bai.botConfigs,
	)
	if err != nil {
		return err
	}
	if of > bai.lastOffset {
		bai.lastOffset = of
	}
	return nil
}

func (bai *BotAPIInterface) isChatIdOk(chatIdInt int, chatIdString string) bool {
	if chatIdInt == 0 {
		return chatIdString != ""
	} else {
		return chatIdString == ""
	}
}

/*GetMe gets the bot info*/
func (bai *BotAPIInterface) GetMe() (*objs.UserResult, error) {
	res, err := bai.SendCustom("getMe", nil, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.UserResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SendMessage sends a message to the user. chatIdInt is used for all chats but channles and chatidString is used for channels (in form of @channleusername) and only of them has be populated, otherwise ChatIdProblem error will be returned.
"chatId" and "text" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) SendMessage(chatIdInt int, chatIdString, text, parseMode string, entities []objs.MessageEntity, disable_web_page_preview, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_to_message_id int, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		def := bai.fixTheDefaultArguments(chatIdInt, reply_to_message_id, chatIdString, disable_notification, allow_sending_without_reply, ProtectContent, reply_markup)
		args := &objs.SendMessageArgs{
			Text:                        text,
			DisableWebPagePreview:       disable_web_page_preview,
			DefaultSendMethodsArguments: def,
			ParseMode:                   parseMode,
			Entities:                    entities,
		}

		res, err := bai.SendCustom("sendMessage", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendMessage"}
	}
}

/*ForwardMessage forwards a message from a user or channel to a user or channel. If the source or destination (or both) of the forwarded message is a channel, only string chat ids should be given to the function, and if it is user only int chat ids should be given.
"chatId", "fromChatId" and "messageId" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) ForwardMessage(chatIdInt, fromChatIdInt int, chatIdString, fromChatIdString string, disableNotif, ProtectContent bool, messageId int) (*objs.SendMethodsResult, error) {
	if (chatIdInt != 0 && chatIdString != "") && (fromChatIdInt != 0 && fromChatIdString != "") {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) && bai.isChatIdOk(fromChatIdInt, fromChatIdString) {
		fm := &objs.ForwardMessageArgs{
			DisableNotification: disableNotif,
			MessageId:           messageId,
			ProtectContent:      ProtectContent,
		}
		fm.ChatId = bai.fixChatId(chatIdInt, chatIdString)
		fm.FromChatId = bai.fixChatId(fromChatIdInt, fromChatIdString)
		res, err := bai.SendCustom("forwardMessage", fm, false, nil, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString or fromChatIdInt or fromChatIdString", MethodName: "forwardMessage"}
	}
}

/*SendPhoto sends a photo (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "photo" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) SendPhoto(chatIdInt int, chatIdString, photo string, photoFile *os.File, caption, parseMode string, reply_to_message_id int, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup, captionEntities []objs.MessageEntity) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendPhotoArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Photo:           photo,
			Caption:         caption,
			ParseMode:       parseMode,
			CaptionEntities: captionEntities,
		}
		var res []byte
		var err error
		if photoFile != nil {
			res, err = bai.SendCustom("sendPhoto", args, true, photoFile, nil)
		} else {
			res, err = bai.SendCustom("sendPhoto", args, false, nil, nil)
		}
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendPhoto"}
	}
}

/*SendVideo sends a video (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "video" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendVideo(chatIdInt int, chatIdString, video string, videoFile *os.File, caption, parseMode string, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply, ProtectContent bool, captionEntities []objs.MessageEntity, duration int, supportsStreaming bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVideoArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Video:             video,
			Caption:           caption,
			Thumb:             thumb,
			ParseMode:         parseMode,
			CaptionEntities:   captionEntities,
			Duration:          duration,
			SupportsStreaming: supportsStreaming,
		}
		res, err := bai.SendCustom("sendVideo", args, true, videoFile, thumbFile)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendVideo"}
	}
}

/*SendAudio sends an audio (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "audio" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0,to ignore string arguments pass "")*/
func (bai *BotAPIInterface) SendAudio(chatIdInt int, chatIdString, audio string, audioFile *os.File, caption, parseMode string, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply, ProtectContent bool, captionEntities []objs.MessageEntity, duration int, performer, title string, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendAudioArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Audio:           audio,
			Caption:         caption,
			Thumb:           thumb,
			ParseMode:       parseMode,
			CaptionEntities: captionEntities,
			Duration:        duration,
			Performer:       performer,
			Title:           title,
		}
		res, err := bai.SendCustom("sendAudio", args, true, audioFile, thumbFile)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendAudio"}
	}
}

/*sSendDocument sends a document (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "document" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendDocument(chatIdInt int, chatIdString, document string, documentFile *os.File, caption, parseMode string, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply, ProtectContent bool, captionEntities []objs.MessageEntity, DisableContentTypeDetection bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendDocumentArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Document:                    document,
			Caption:                     caption,
			Thumb:                       thumb,
			ParseMode:                   parseMode,
			CaptionEntities:             captionEntities,
			DisableContentTypeDetection: DisableContentTypeDetection,
		}
		res, err := bai.SendCustom("sendDocument", args, true, documentFile, thumbFile)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendDocument"}
	}
}

/*SendAnimation sends an animation (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "animation" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendAnimation(chatIdInt int, chatIdString, animation string, animationFile *os.File, caption, parseMode string, width, height, duration int, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply, ProtectContent bool, captionEntities []objs.MessageEntity, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendAnimationArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Animation:       animation,
			Caption:         caption,
			Thumb:           thumb,
			ParseMode:       parseMode,
			CaptionEntities: captionEntities,
			Width:           width,
			Height:          height,
			Duration:        duration,
		}
		res, err := bai.SendCustom("sendAnimation", args, true, animationFile, thumbFile)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendAnimation"}
	}
}

/*sSendVoice sends a voice (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "voice" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendVoice(chatIdInt int, chatIdString, voice string, voiceFile *os.File, caption, parseMode string, duration int, reply_to_message_id int, disable_notification, allow_sending_without_reply, ProtectContent bool, captionEntities []objs.MessageEntity, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVoiceArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Voice:           voice,
			Caption:         caption,
			ParseMode:       parseMode,
			CaptionEntities: captionEntities,
			Duration:        duration,
		}
		res, err := bai.SendCustom("sendVoice", args, true, voiceFile)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendVoice"}
	}
}

/*SendVideoNote sends a video note (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "videoNote" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)
Note that sending video note by URL is not supported by telegram.*/
func (bai *BotAPIInterface) SendVideoNote(chatIdInt int, chatIdString, videoNote string, videoNoteFile *os.File, caption, parseMode string, length, duration int, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply, ProtectContent bool, captionEntities []objs.MessageEntity, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVideoNoteArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			VideoNote:       videoNote,
			Caption:         caption,
			Thumb:           thumb,
			ParseMode:       parseMode,
			CaptionEntities: captionEntities,
			Length:          length,
			Duration:        duration,
		}
		res, err := bai.SendCustom("sendVideoNote", args, true, videoNoteFile, thumbFile)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendVideoNote"}
	}
}

/*SendMediaGroup sends an album of media (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "media" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendMediaGroup(chatIdInt int, chatIdString string, reply_to_message_id int, media []objs.InputMedia, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup, files ...*os.File) (*objs.SendMediaGroupMethodResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendMediaGroupArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Media: media,
		}
		res, err := bai.SendCustom("sendMediaGroup", args, true, files...)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMediaGroupMethodResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendMediaGRoup"}
	}
}

/*SendLocation sends a location to a channel (chatIdString) or a chat (chatIdInt)
"chatId","latitude" and "longitude" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendLocation(chatIdInt int, chatIdString string, latitude, longitude, horizontalAccuracy float32, livePeriod, heading, proximityAlertRadius, reply_to_message_id int, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendLocationArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Latitude:             latitude,
			Longitude:            longitude,
			HorizontalAccuracy:   horizontalAccuracy,
			LivePeriod:           livePeriod,
			Heading:              heading,
			ProximityAlertRadius: proximityAlertRadius,
		}
		res, err := bai.SendCustom("sendLocation", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendLocation"}
	}
}

/*EditMessageLiveLocation edits a live location sent to a channel (chatIdString) or a chat (chatIdInt)
"chatId","latitude" and "longitude" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) EditMessageLiveLocation(chatIdInt int, chatIdString, inlineMessageId string, messageId int, latitude, longitude, horizontalAccuracy float32, heading, proximityAlertRadius int, reply_markup *objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.EditMessageLiveLocationArgs{
			InlineMessageId:      inlineMessageId,
			MessageId:            messageId,
			Latitude:             latitude,
			Longitude:            longitude,
			HorizontalAccuracy:   horizontalAccuracy,
			Heading:              heading,
			ProximityAlertRadius: proximityAlertRadius,
			ReplyMarkup:          reply_markup,
		}
		args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
		res, err := bai.SendCustom("editMessageLiveLocation", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.DefaultResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "editMessageLiveLocation"}
	}
}

/*StopMessageLiveLocation stops a live location sent to a channel (chatIdString) or a chat (chatIdInt)
"chatId" argument is required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) StopMessageLiveLocation(chatIdInt int, chatIdString, inlineMessageId string, messageId int, replyMarkup *objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.StopMessageLiveLocationArgs{
			InlineMessageId: inlineMessageId,
			MessageId:       messageId,
			ReplyMarkup:     replyMarkup,
		}
		args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
		res, err := bai.SendCustom("stopMessageLiveLocation", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.DefaultResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "stopMessageLiveLocation"}
	}
}

/*SendVenue sends a venue to a channel (chatIdString) or a chat (chatIdInt)
"chatId","latitude","longitude","title" and "address" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendVenue(chatIdInt int, chatIdString string, latitude, longitude float32, title, address, fourSquareId, fourSquareType, googlePlaceId, googlePlaceType string, reply_to_message_id int, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVenueArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Latitude:        latitude,
			Longitude:       longitude,
			Title:           title,
			Address:         address,
			FoursquareId:    fourSquareId,
			FoursquareType:  fourSquareType,
			GooglePlaceId:   googlePlaceId,
			GooglePlaceType: googlePlaceType,
		}
		res, err := bai.SendCustom("sendVnue", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendcontact"}
	}
}

/*SendContact sends a contact to a channel (chatIdString) or a chat (chatIdInt)
"chatId","phoneNumber" and "firstName" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendContact(chatIdInt int, chatIdString, phoneNumber, firstName, lastName, vCard string, reply_to_message_id int, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendContactArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			PhoneNumber: phoneNumber,
			FirstName:   firstName,
			LastName:    lastName,
			Vcard:       vCard,
		}
		res, err := bai.SendCustom("sendContact", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendContact"}
	}
}

/*SendPoll sends a poll to a channel (chatIdString) or a chat (chatIdInt)
"chatId","phoneNumber" and "firstName" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendPoll(chatIdInt int, chatIdString, question string, options []string, isClosed, isAnonymous bool, pollType string, allowMultipleAnswers bool, correctOptionIndex int, explanation, explanationParseMode string, explanationEntities []objs.MessageEntity, openPeriod, closeDate int, reply_to_message_id int, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendPollArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Question:              question,
			Options:               options,
			IsClosed:              isClosed,
			IsAnonymous:           isAnonymous,
			Type:                  pollType,
			AllowsMultipleAnswers: allowMultipleAnswers,
			CorrectOptionId:       correctOptionIndex,
			Explanation:           explanation,
			ExplanationParseMode:  explanationParseMode,
			ExplanationEntities:   explanationEntities,
			OpenPeriod:            openPeriod,
			CloseDate:             closeDate,
		}
		res, err := bai.SendCustom("sendPoll", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendPoll"}
	}
}

/*SendDice sends a dice message to a channel (chatIdString) or a chat (chatIdInt)
"chatId" argument is required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendDice(chatIdInt int, chatIdString, emoji string, reply_to_message_id int, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendDiceArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, ProtectContent, reply_markup,
			),
			Emoji: emoji,
		}
		res, err := bai.SendCustom("sendDice", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendDice"}
	}
}

/*SendChatAction sends a chat action message to a channel (chatIdString) or a chat (chatIdInt)
"chatId" argument is required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendChatAction(chatIdInt int, chatIdString, chatAction string) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendChatActionArgs{
			Action: chatAction,
		}
		args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
		res, err := bai.SendCustom("sendChatAction", args, false, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString", MethodName: "sendChatAction"}
	}
}

/*GetUserProfilePhotos gets the user profile photos*/
func (bai *BotAPIInterface) GetUserProfilePhotos(userId, offset, limit int) (*objs.ProfilePhototsResult, error) {
	args := &objs.GetUserProfilePhototsArgs{UserId: userId, Offset: offset, Limit: limit}
	res, err := bai.SendCustom("getUserProfilePhotos", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.ProfilePhototsResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetFile gets the file based on the given file id and returns the file object. */
func (bai *BotAPIInterface) GetFile(fileId string) (*objs.GetFileResult, error) {
	args := &objs.GetFileArgs{FileId: fileId}
	res, err := bai.SendCustom("getFile", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.GetFileResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DownloadFile downloads a file from telegram servers and saves it into the given file.

This method closes the given file. If the file is nil, this method will create a file based on the name of the file stored in telegram servers.*/
func (bai *BotAPIInterface) DownloadFile(fileObject *objs.File, file *os.File) error {
	url := "https://api.telegram.org/file/bot" + bai.botConfigs.APIKey + "/" + fileObject.FilePath
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if file == nil {
		ar := strings.Split(fileObject.FilePath, "/")
		name := ar[len(ar)-1]
		var er error
		file, er = os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0666)
		if er != nil {
			return er
		}
	}
	if res.StatusCode < 300 {
		_, err2 := io.Copy(file, res.Body)
		if err2 != nil {
			return err2
		}
		err3 := file.Close()
		return err3
	} else {
		return &errs.MethodNotSentError{Method: "getFile", Reason: "server returned status code " + strconv.Itoa(res.StatusCode)}
	}
}

/*BanChatMember bans a chat member*/
func (bai *BotAPIInterface) BanChatMember(chatIdInt int, chatIdString string, userId, untilDate int, revokeMessages bool) (*objs.LogicalResult, error) {
	args := &objs.BanChatMemberArgs{
		UserId:         userId,
		UntilDate:      untilDate,
		RevokeMessages: revokeMessages,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("banChatMember", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*UnbanChatMember unbans a chat member*/
func (bai *BotAPIInterface) UnbanChatMember(chatIdInt int, chatIdString string, userId int, onlyIfBanned bool) (*objs.LogicalResult, error) {
	args := &objs.UnbanChatMemberArgsArgs{
		UserId:       userId,
		OnlyIfBanned: onlyIfBanned,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("unbanChatMember", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*RestrictChatMember restricts a chat member*/
func (bai *BotAPIInterface) RestrictChatMember(chatIdInt int, chatIdString string, userId int, permissions objs.ChatPermissions, untilDate int) (*objs.LogicalResult, error) {
	args := &objs.RestrictChatMemberArgs{
		UserId:     userId,
		Permission: permissions,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("restrictChatMember", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*PromoteChatMember promotes a chat member*/
func (bai *BotAPIInterface) PromoteChatMember(chatIdInt int, chatIdString string, userId int, isAnonymous, canManageChat, canPostmessages, canEditMessages, canDeleteMessages, canManageVoiceChats, canRestrictMembers, canPromoteMembers, canChangeInfo, canInviteUsers, canPinMessages bool) (*objs.LogicalResult, error) {
	args := &objs.PromoteChatMemberArgs{
		UserId:              userId,
		IsAnonymous:         isAnonymous,
		CanManageChat:       canManageChat,
		CanPostMessages:     canPostmessages,
		CanEditMessages:     canEditMessages,
		CanDeleteMessages:   canDeleteMessages,
		CanManageVoiceChats: canManageVoiceChats,
		CanRestrictMembers:  canRestrictMembers,
		CanPromoteMembers:   canPromoteMembers,
		CanChangeInfo:       canChangeInfo,
		CanInviteUsers:      canInviteUsers,
		CanPinMessages:      canPinMessages,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("promoteChatMember", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetChatAdministratorCustomTitle sets a custom title for the administrator.*/
func (bai *BotAPIInterface) SetChatAdministratorCustomTitle(chatIdInt int, chatIdString string, userId int, customTitle string) (*objs.LogicalResult, error) {
	args := &objs.SetChatAdministratorCustomTitleArgs{
		UserId:      userId,
		CustomTitle: customTitle,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("setChatAdministratorCustomTitle", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*BanOrUnbanChatSenderChat bans or unbans a channel in the group..*/
func (bai *BotAPIInterface) BanOrUnbanChatSenderChat(chatIdInt int, chatIdString string, senderChatId int, ban bool) (*objs.LogicalResult, error) {
	args := &objs.BanChatSenderChatArgs{
		SenderChatId: senderChatId,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	var method string
	if ban {
		method = "banChatSenderChat"
	} else {
		method = "unbanChatSenderChat"
	}
	res, err := bai.SendCustom(method, args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetChatPermissions sets default permissions for all users in the chat.*/
func (bai *BotAPIInterface) SetChatPermissions(chatIdInt int, chatIdString string, permissions objs.ChatPermissions) (*objs.LogicalResult, error) {
	args := &objs.SetChatPermissionsArgs{
		Permissions: permissions,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("setChatPermissions", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*ExportChatInviteLink exports the chat invite link and returns the new invite link as string.*/
func (bai *BotAPIInterface) ExportChatInviteLink(chatIdInt int, chatIdString string) (*objs.StringResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("exprotChatInviteLink", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.StringResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*CreateChatInviteLink creates a new invite link for the chat.*/
func (bai *BotAPIInterface) CreateChatInviteLink(chatIdInt int, chatIdString, name string, expireDate, memberLimit int, createsJoinRequest bool) (*objs.ChatInviteLinkResult, error) {
	args := &objs.CreateChatInviteLinkArgs{
		Name:               name,
		ExpireDate:         expireDate,
		MemberLimit:        memberLimit,
		CreatesjoinRequest: createsJoinRequest,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("createChatInviteLink", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.ChatInviteLinkResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*EditChatInviteLink edits an existing invite link for the chat.*/
func (bai *BotAPIInterface) EditChatInviteLink(chatIdInt int, chatIdString, inviteLink, name string, expireDate, memberLimit int, createsJoinRequest bool) (*objs.ChatInviteLinkResult, error) {
	args := &objs.EditChatInviteLinkArgs{
		InviteLink:         inviteLink,
		Name:               name,
		ExpireDate:         expireDate,
		MemberLimit:        memberLimit,
		CreatesjoinRequest: createsJoinRequest,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("editChatInviteLink", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.ChatInviteLinkResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*RevokeChatInviteLink revokes the given invite link.*/
func (bai *BotAPIInterface) RevokeChatInviteLink(chatIdInt int, chatIdString, inviteLink string) (*objs.ChatInviteLinkResult, error) {
	args := &objs.RevokeChatInviteLinkArgs{
		InviteLink: inviteLink,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("revokeChatInviteLink", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.ChatInviteLinkResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*ApproveChatJoinRequest approves a request from the given user to join the chat.*/
func (bai *BotAPIInterface) ApproveChatJoinRequest(chatIdInt int, chatIdString string, userId int) (*objs.LogicalResult, error) {
	args := &objs.ApproveChatJoinRequestArgs{
		UserId: userId,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("approveChatJoinRequest", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DeclineChatJoinRequest declines a request from the given user to join the chat.*/
func (bai *BotAPIInterface) DeclineChatJoinRequest(chatIdInt int, chatIdString string, userId int) (*objs.LogicalResult, error) {
	args := &objs.DeclineChatJoinRequestArgs{
		UserId: userId,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("declineChatJoinRequest", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetChatPhoto sets the chat photo to given file.*/
func (bai *BotAPIInterface) SetChatPhoto(chatIdInt int, chatIdString string, file *os.File) (*objs.LogicalResult, error) {
	args := &objs.SetChatPhotoArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	stats, er := file.Stat()
	if er != nil {
		return nil, er
	}
	args.Photo = "attach://" + stats.Name()
	res, err := bai.SendCustom("setChatPhoto", args, true, file)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DeleteChatPhoto deletes chat photo.*/
func (bai *BotAPIInterface) DeleteChatPhoto(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("deleteChatPhoto", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetChatTitle sets the chat title.*/
func (bai *BotAPIInterface) SetChatTitle(chatIdInt int, chatIdString, title string) (*objs.LogicalResult, error) {
	args := &objs.SetChatTitleArgs{
		Title: title,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("setChatTitle", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetChatDescription sets the chat description.*/
func (bai *BotAPIInterface) SetChatDescription(chatIdInt int, chatIdString, descriptions string) (*objs.LogicalResult, error) {
	args := &objs.SetChatDescriptionArgs{
		Description: descriptions,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("setChatDescription", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*PinChatMessage pins the message in the chat.*/
func (bai *BotAPIInterface) PinChatMessage(chatIdInt int, chatIdString string, messageId int, disableNotification bool) (*objs.LogicalResult, error) {
	args := &objs.PinChatMessageArgs{
		MessageId:           messageId,
		DisableNotification: disableNotification,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("pinChatMessage", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*UnpinChatMessage unpins the pinned message in the chat.*/
func (bai *BotAPIInterface) UnpinChatMessage(chatIdInt int, chatIdString string, messageId int) (*objs.LogicalResult, error) {
	args := &objs.UnpinChatMessageArgs{
		MessageId: messageId,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("unpinChatMessage", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*UnpinAllChatMessages unpins all the pinned messages in the chat.*/
func (bai *BotAPIInterface) UnpinAllChatMessages(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("unpinAllChatMessages", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*LeaveChat, the bot will leave the chat if this method is called.*/
func (bai *BotAPIInterface) LeaveChat(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("leaveChat", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetChat : a Chat object containing the information of the chat will be returned*/
func (bai *BotAPIInterface) GetChat(chatIdInt int, chatIdString string) (*objs.ChatResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("getChat", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.ChatResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetChatAdministrators returns an array of ChatMember containing the informations of the chat administrators.*/
func (bai *BotAPIInterface) GetChatAdministrators(chatIdInt int, chatIdString string) (*objs.ChatAdministratorsResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("getChatAdministrators", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.ChatAdministratorsResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/* GetChatMemberCount returns the number of the memebrs of the chat.*/
func (bai *BotAPIInterface) GetChatMemberCount(chatIdInt int, chatIdString string) (*objs.IntResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("getChatMemberCount", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.IntResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetChatMember returns the information of the member in a ChatMember object.*/
func (bai *BotAPIInterface) GetChatMember(chatIdInt int, chatIdString string, userId int) (*objs.DefaultResult, error) {
	args := &objs.GetChatMemberArgs{
		UserId: userId,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("getChatMember", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.DefaultResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetChatStickerSet sets the sticker set of the chat.*/
func (bai *BotAPIInterface) SetChatStickerSet(chatIdInt int, chatIdString, stickerSetName string) (*objs.LogicalResult, error) {
	args := &objs.SetChatStcikerSet{
		StickerSetName: stickerSetName,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("setChatStickerSet", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DeleteChatStickerSet deletes the sticker set of the chat..*/
func (bai *BotAPIInterface) DeleteChatStickerSet(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("deleteChatStickerSet", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*AnswerCallbackQuery answers a callback query*/
func (bai *BotAPIInterface) AnswerCallbackQuery(callbackQueryId, text, url string, showAlert bool, CacheTime int) (*objs.LogicalResult, error) {
	args := &objs.AnswerCallbackQueryArgs{
		CallbackQueyId: callbackQueryId,
		Text:           text,
		ShowAlert:      showAlert,
		URL:            url,
		CacheTime:      CacheTime,
	}
	res, err := bai.SendCustom("answerCallbackQuery", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetMyCommands sets the commands of the bot*/
func (bai *BotAPIInterface) SetMyCommands(commands []objs.BotCommand, scope objs.BotCommandScope, languageCode string) (*objs.LogicalResult, error) {
	args := &objs.SetMyCommandsArgs{
		Commands: commands,
		MyCommandsDefault: objs.MyCommandsDefault{
			Scope:        scope,
			LanguageCode: languageCode,
		},
	}
	res, err := bai.SendCustom("setMyCommands", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DeleteMyCommands deletes the commands of the bot*/
func (bai *BotAPIInterface) DeleteMyCommands(scope objs.BotCommandScope, languageCode string) (*objs.LogicalResult, error) {
	args := &objs.MyCommandsDefault{
		Scope:        scope,
		LanguageCode: languageCode,
	}
	res, err := bai.SendCustom("deleteMyCommands", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetMyCommands gets the commands of the bot*/
func (bai *BotAPIInterface) GetMyCommands(scope objs.BotCommandScope, languageCode string) (*objs.GetCommandsResult, error) {
	args := &objs.MyCommandsDefault{
		Scope:        scope,
		LanguageCode: languageCode,
	}
	res, err := bai.SendCustom("getMyCommands", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.GetCommandsResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*EditMessageText edits the text of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessageText(chatIdInt int, chatIdString string, messageId int, inlineMessageId, text, parseMode string, entities []objs.MessageEntity, disableWebPagePreview bool, replyMakrup *objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	args := &objs.EditMessageTextArgs{
		EditMessageDefaultArgs: objs.EditMessageDefaultArgs{
			MessageId:       messageId,
			InlineMessageId: inlineMessageId,
			ReplyMarkup:     replyMakrup,
		},
		Text:                  text,
		ParseMode:             parseMode,
		Entities:              entities,
		DisablewebpagePreview: disableWebPagePreview,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("editMessageText", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.DefaultResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*EditMessageCaption edits the caption of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessageCaption(chatIdInt int, chatIdString string, messageId int, inlineMessageId, caption, parseMode string, captionEntities []objs.MessageEntity, replyMakrup *objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	args := &objs.EditMessageCaptionArgs{
		EditMessageDefaultArgs: objs.EditMessageDefaultArgs{
			MessageId:       messageId,
			InlineMessageId: inlineMessageId,
			ReplyMarkup:     replyMakrup,
		},
		Caption:         caption,
		ParseMode:       parseMode,
		CaptionEntities: captionEntities,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("editMessageCaption", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.DefaultResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*EditMessageMedia edits the media of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessageMedia(chatIdInt int, chatIdString string, messageId int, inlineMessageId string, media objs.InputMedia, replyMakrup *objs.InlineKeyboardMarkup, file ...*os.File) (*objs.DefaultResult, error) {
	args := &objs.EditMessageMediaArgs{
		EditMessageDefaultArgs: objs.EditMessageDefaultArgs{
			MessageId:       messageId,
			InlineMessageId: inlineMessageId,
			ReplyMarkup:     replyMakrup,
		},
		Media: media,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("editMessageMedia", args, true, file...)
	if err != nil {
		return nil, err
	}
	msg := &objs.DefaultResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*EditMessagereplyMarkup edits the reply makrup of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessagereplyMarkup(chatIdInt int, chatIdString string, messageId int, inlineMessageId string, replyMakrup *objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	args := &objs.EditMessageReplyMakrupArgs{
		EditMessageDefaultArgs: objs.EditMessageDefaultArgs{
			MessageId:       messageId,
			InlineMessageId: inlineMessageId,
			ReplyMarkup:     replyMakrup,
		},
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("editMessageReplyMarkup", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.DefaultResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*StopPoll stops the poll.*/
func (bai *BotAPIInterface) StopPoll(chatIdInt int, chatIdString string, messageId int, replyMakrup *objs.InlineKeyboardMarkup) (*objs.PollResult, error) {
	args := &objs.StopPollArgs{
		MessageId:   messageId,
		ReplyMarkup: replyMakrup,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("stopPoll", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.PollResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DeleteMessage deletes the given message int the given chat.*/
func (bai *BotAPIInterface) DeleteMessage(chatIdInt int, chatIdString string, messageId int) (*objs.LogicalResult, error) {
	args := &objs.DeleteMessageArgs{
		MessageId: messageId,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("deleteMessage", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SendSticker sends an sticker to the given chat id.*/
func (bai *BotAPIInterface) SendSticker(chatIdInt int, chatIdString, sticker string, disableNotif, allowSendingWithoutreply, protectContent bool, replyTo int, replyMarkup objs.ReplyMarkup, file *os.File) (*objs.SendMethodsResult, error) {
	args := &objs.SendStickerArgs{
		DefaultSendMethodsArguments: objs.DefaultSendMethodsArguments{
			DisableNotification:      disableNotif,
			AllowSendingWithoutReply: allowSendingWithoutreply,
			ReplyToMessageId:         replyTo,
			ReplyMarkup:              replyMarkup,
			ProtectContent:           protectContent,
		},
		Sticker: sticker,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("sendSticker", args, true, file)
	if err != nil {
		return nil, err
	}
	msg := &objs.SendMethodsResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetStickerSet gets the sticker set by the given name*/
func (bai *BotAPIInterface) GetStickerSet(name string) (*objs.StickerSetResult, error) {
	args := &objs.GetStickerSetArgs{
		Name: name,
	}
	res, err := bai.SendCustom("getStickerSet", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.StickerSetResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*UploadStickerFile uploads the given file as an sticker on the telegram servers.*/
func (bai *BotAPIInterface) UploadStickerFile(userId int, pngSticker string, file *os.File) (*objs.GetFileResult, error) {
	args := &objs.UploadStickerFileArgs{
		UserId:     userId,
		PngSticker: pngSticker,
	}
	res, err := bai.SendCustom("uploadStickerFile", args, true, file)
	if err != nil {
		return nil, err
	}
	msg := &objs.GetFileResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*CreateNewStickerSet creates a new sticker set with the given arguments*/
func (bai *BotAPIInterface) CreateNewStickerSet(userId int, name, title, pngSticker, tgsSticker, webmSticker, emojies string, containsMasks bool, maskPosition *objs.MaskPosition, file *os.File) (*objs.LogicalResult, error) {
	args := &objs.CreateNewStickerSetArgs{
		UserId:        userId,
		Name:          name,
		Title:         title,
		Emojis:        emojies,
		PngSticker:    pngSticker,
		TgsSticker:    tgsSticker,
		WebmSticker:   webmSticker,
		ContainsMasks: containsMasks,
		MaskPosition:  maskPosition,
	}
	res, err := bai.SendCustom("createNewStickerSet", args, true, file)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*AddStickerToSet adds a new sticker to the given set.*/
func (bai *BotAPIInterface) AddStickerToSet(userId int, name, pngSticker, tgsSticker, webmSticker, emojies string, maskPosition *objs.MaskPosition, file *os.File) (*objs.LogicalResult, error) {
	args := &objs.AddStickerSetArgs{
		UserId:       userId,
		Name:         name,
		PngSticker:   pngSticker,
		TgsSticker:   tgsSticker,
		WebmSticker:  webmSticker,
		Emojis:       emojies,
		MaskPosition: maskPosition,
	}
	res, err := bai.SendCustom("addStickerToSet", args, true, file)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetStickerPositionInSet sets the position of a sticker in an sticker set*/
func (bai *BotAPIInterface) SetStickerPositionInSet(sticker string, position int) (*objs.LogicalResult, error) {
	args := &objs.SetStickerPositionInSetArgs{
		Sticker:  sticker,
		Position: position,
	}
	res, err := bai.SendCustom("setStickerPositionInSet", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DeleteStickerFromSet deletes the given sticker from a set created by the bot*/
func (bai *BotAPIInterface) DeleteStickerFromSet(sticker string) (*objs.LogicalResult, error) {
	args := &objs.DeleteStickerFromSetArgs{
		Sticker: sticker,
	}
	res, err := bai.SendCustom("deleteStickerFromSet", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetStickerSetThumb sets the thumbnail for the given sticker*/
func (bai *BotAPIInterface) SetStickerSetThumb(name, thumb string, userId int, file *os.File) (*objs.LogicalResult, error) {
	args := &objs.SetStickerSetThumbArgs{
		Name:   name,
		Thumb:  thumb,
		UserId: userId,
	}
	res, err := bai.SendCustom("setStickerSetThumb", args, true, file)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*AnswerInlineQuery answers an inline query with the given parameters*/
func (bai *BotAPIInterface) AnswerInlineQuery(inlineQueryId string, results []objs.InlineQueryResult, cacheTime int, isPersonal bool, nextOffset, switchPmText, switchPmParameter string) (*objs.LogicalResult, error) {
	args := &objs.AnswerInlineQueryArgs{
		InlineQueryId:     inlineQueryId,
		Results:           results,
		CacheTime:         cacheTime,
		IsPersonal:        isPersonal,
		NextOffset:        nextOffset,
		SwitchPmText:      switchPmText,
		SwitchPmParameter: switchPmParameter,
	}
	res, err := bai.SendCustom("answerInlineQuery", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SendInvoice sends an invoice*/
func (bai *BotAPIInterface) SendInvoice(chatIdInt int, chatIdString, title, description, payload, providerToken, currency string, prices []objs.LabeledPrice, maxTipAmount int, suggestedTipAmounts []int, startParameter, providerData, photoURL string, photoSize, photoWidth, photoHeight int, needName, needPhoneNumber, needEmail, needSippingAddress, sendPhoneNumberToProvider, sendEmailToProvider, isFlexible, disableNotif bool, replyToMessageId int, allowSendingWithoutReply bool, replyMarkup objs.InlineKeyboardMarkup) (*objs.SendMethodsResult, error) {
	args := &objs.SendInvoiceArgs{
		DefaultSendMethodsArguments: objs.DefaultSendMethodsArguments{
			DisableNotification:      disableNotif,
			AllowSendingWithoutReply: allowSendingWithoutReply,
			ReplyToMessageId:         replyToMessageId,
			ReplyMarkup:              &replyMarkup,
		},
		Title:                     title,
		Description:               description,
		Payload:                   payload,
		ProviderToken:             providerToken,
		Currency:                  currency,
		Prices:                    prices,
		MaxTipAmount:              maxTipAmount,
		SuggestedTipAmounts:       suggestedTipAmounts,
		StartParameter:            startParameter,
		ProviderData:              providerData,
		PhotoURL:                  photoURL,
		PhotoSize:                 photoSize,
		PhotoWidth:                photoWidth,
		PhotoHeight:               photoHeight,
		NeedName:                  needName,
		NeedPhoneNumber:           needPhoneNumber,
		NeedEmail:                 needEmail,
		NeedShippingAddress:       needSippingAddress,
		SendPhoneNumberToProvider: sendPhoneNumberToProvider,
		SendEmailToProvider:       sendEmailToProvider,
		IsFlexible:                isFlexible,
	}
	args.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	res, err := bai.SendCustom("sendInvoice", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.SendMethodsResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*AnswerShippingQuery answers a shipping query*/
func (bai *BotAPIInterface) AnswerShippingQuery(shippingQueryId string, ok bool, shippingOptions []objs.ShippingOption, errorMessage string) (*objs.LogicalResult, error) {
	args := &objs.AnswerShippingQueryArgs{
		ShippingQueryId: shippingQueryId,
		OK:              ok,
		ShippingOptions: shippingOptions,
		ErrorMessage:    errorMessage,
	}
	res, err := bai.SendCustom("answerShippingQuery", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*AnswerPreCheckoutQuery answers a pre checkout query*/
func (bai *BotAPIInterface) AnswerPreCheckoutQuery(preCheckoutQueryId string, ok bool, errorMessage string) (*objs.LogicalResult, error) {
	args := &objs.AnswerPreCheckoutQueryArgs{
		PreCheckoutQueryId: preCheckoutQueryId,
		Ok:                 ok,
		ErrorMessage:       errorMessage,
	}
	res, err := bai.SendCustom("answerPreCheckoutQuery", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*CopyMessage copies a message from a user or channel and sends it to a user or channel. If the source or destination (or both) of the forwarded message is a channel, only string chat ids should be given to the function, and if it is user only int chat ids should be given.
"chatId", "fromChatId" and "messageId" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) CopyMessage(chatIdInt, fromChatIdInt int, chatIdString, fromChatIdString string, messageId int, disableNotif bool, caption, parseMode string, replyTo int, allowSendingWihtoutReply, ProtectContent bool, replyMarkUp objs.ReplyMarkup, captionEntities []objs.MessageEntity) (*objs.SendMethodsResult, error) {
	if (chatIdInt != 0 && chatIdString != "") && (fromChatIdInt != 0 && fromChatIdString != "") {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) && bai.isChatIdOk(fromChatIdInt, fromChatIdString) {
		fm := objs.ForwardMessageArgs{
			DisableNotification: disableNotif,
			MessageId:           messageId,
			ProtectContent:      ProtectContent,
		}
		fm.ChatId = bai.fixChatId(chatIdInt, chatIdString)
		fm.FromChatId = bai.fixChatId(fromChatIdInt, fromChatIdString)
		cp := &objs.CopyMessageArgs{
			ForwardMessageArgs:       fm,
			Caption:                  caption,
			ParseMode:                parseMode,
			AllowSendingWithoutReply: allowSendingWihtoutReply,
			ReplyMarkup:              replyMarkUp,
			CaptionEntities:          captionEntities,
		}
		if replyTo != 0 {
			cp.ReplyToMessageId = replyTo
		}
		res, err := bai.SendCustom("copyMessage", cp, false, nil, nil)
		if err != nil {
			return nil, err
		}
		msg := &objs.SendMethodsResult{}
		err3 := json.Unmarshal(res, msg)
		if err3 != nil {
			return nil, err3
		}
		return msg, nil
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString or fromChatIdInt or fromChatIdString", MethodName: "copyMessage"}
	}
}

/*SetPassportDataErrors sets passport data errors*/
func (bai *BotAPIInterface) SetPassportDataErrors(userId int, errors []objs.PassportElementError) (*objs.LogicalResult, error) {
	args := &objs.SetPassportDataErrorsArgs{
		UserId: userId, Errors: errors,
	}
	res, err := bai.SendCustom("setPassportDataErrors", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SendGame sends a game*/
func (bai *BotAPIInterface) SendGame(chatId int, gameShortName string, disableNotif bool, replyTo int, allowSendingWithoutReply bool, replyMarkup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	args := &objs.SendGameArgs{
		DefaultSendMethodsArguments: objs.DefaultSendMethodsArguments{
			ReplyToMessageId:         replyTo,
			DisableNotification:      disableNotif,
			ReplyMarkup:              replyMarkup,
			AllowSendingWithoutReply: allowSendingWithoutReply,
		},
		GameShortName: gameShortName,
	}
	bt, _ := json.Marshal(chatId)
	args.ChatId = bt
	res, err := bai.SendCustom("sendGame", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.SendMethodsResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetGameScore sets the game high score*/
func (bai *BotAPIInterface) SetGameScore(userId, score int, force, disableEditMessage bool, chatId, messageId int, inlineMessageId string) (*objs.DefaultResult, error) {
	args := &objs.SetGameScoreArgs{
		UserId:             userId,
		Score:              score,
		Force:              force,
		DisableEditMessage: disableEditMessage,
		ChatId:             chatId,
		MessageId:          messageId,
		InlineMessageId:    inlineMessageId,
	}
	res, err := bai.SendCustom("setGameScore", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.DefaultResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetGameHighScores gets the high scores of the user*/
func (bai *BotAPIInterface) GetGameHighScores(userId, chatId, messageId int, inlineMessageId string) (*objs.GameHighScoresResult, error) {
	args := &objs.GetGameHighScoresArgs{
		UserId:          userId,
		ChatId:          chatId,
		MessageId:       messageId,
		InlineMessageId: inlineMessageId,
	}
	res, err := bai.SendCustom("getGameHighScores", args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.GameHighScoresResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*GetWebhookInfo returns the web hook info of the bot.*/
func (bai *BotAPIInterface) GetWebhookInfo() (*objs.WebhookInfoResult, error) {
	res, err := bai.SendCustom("getWebhookInfo", nil, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.WebhookInfoResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SetWebhook sets a webhook for the bot.*/
func (bai *BotAPIInterface) SetWebhook(url, ip string, maxCnc int, allowedUpdates []string, dropPendingUpdates bool, keyFile *os.File) (*objs.LogicalResult, error) {
	stat, errs := keyFile.Stat()
	if errs != nil {
		return nil, errs
	}
	args := objs.SetWebhookArgs{
		URL:                url,
		IPAddress:          ip,
		Certificate:        "attach://" + stat.Name(),
		MaxConnections:     maxCnc,
		AllowedUpdates:     allowedUpdates,
		DropPendingUpdates: dropPendingUpdates,
	}
	res, err := bai.SendCustom("setWebhook", &args, true, keyFile)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*DeleteWebhook deletes the webhook for this bot*/
func (bai *BotAPIInterface) DeleteWebhook(dropPendingUpdates bool) (*objs.LogicalResult, error) {
	args := objs.DeleteWebhookArgs{
		DropPendingUpdates: dropPendingUpdates,
	}
	res, err := bai.SendCustom("deleteWebhook", &args, false, nil)
	if err != nil {
		return nil, err
	}
	msg := &objs.LogicalResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
}

/*SendCustom calls the given method on api server with the given arguments. "MP" options indicates that the request should be made in multipart/formdata form. If this method sends a file to the api server the "MP" option should be true*/
func (bai *BotAPIInterface) SendCustom(methodName string, args objs.MethodArguments, MP bool, files ...*os.File) ([]byte, error) {
	start := time.Now().UnixMicro()
	cl := httpSenderClient{botApi: bai.botConfigs.BotAPI, apiKey: bai.botConfigs.APIKey}
	var res []byte
	var err2 error
	if MP {
		res, err2 = cl.sendHttpReqMultiPart(methodName, args, files...)
	} else {
		res, err2 = cl.sendHttpReqJson(methodName, args)
	}
	done := time.Now().UnixMicro()
	if err2 != nil {
		logger.Log(methodName, "\t\t\t", "Error  ", strconv.FormatInt((done-start), 10)+"µs", logger.BOLD+logger.OKBLUE, logger.FAIL, "")
		return nil, err2
	}
	logger.Log(methodName, "\t\t\t", "Success", strconv.FormatInt((done-start), 10)+"µs", logger.BOLD+logger.OKBLUE, logger.OKGREEN, "")
	return bai.preParseResult(res, methodName)
}

func (bai *BotAPIInterface) fixTheDefaultArguments(chatIdInt, reply_to_message_id int, chatIdString string, disable_notification, allow_sending_without_reply, ProtectContent bool, reply_markup objs.ReplyMarkup) objs.DefaultSendMethodsArguments {
	def := objs.DefaultSendMethodsArguments{
		DisableNotification:      disable_notification,
		AllowSendingWithoutReply: allow_sending_without_reply,
		ProtectContent:           ProtectContent,
		ReplyToMessageId:         reply_to_message_id,
		ReplyMarkup:              reply_markup,
	}
	def.ChatId = bai.fixChatId(chatIdInt, chatIdString)
	return def
}

func (bai *BotAPIInterface) preParseResult(res []byte, method string) ([]byte, error) {
	def := &objs.DefaultResult{}
	err := json.Unmarshal(res, def)
	if err != nil {
		return nil, err
	}
	if !def.Ok {
		fr := &objs.FailureResult{}
		err := json.Unmarshal(res, fr)
		if err != nil {
			return nil, err
		}
		return nil, &errs.MethodNotSentError{Method: method, Reason: "server returned false ok filed", FailureResult: fr}
	}
	return res, nil
}

func (bai *BotAPIInterface) fixChatId(chatIdInt int, chatIdString string) []byte {
	if chatIdInt == 0 {
		if !strings.HasPrefix(chatIdString, "@") {
			chatIdString = "@" + chatIdString
		}
		bt, _ := json.Marshal(chatIdString)
		return bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		return bt
	}
}

/*CreateInterface returns an iterface to communicate with the bot api.
If the updateFrequency argument is not nil, the update routine begins automtically*/
func CreateInterface(botCfg *cfgs.BotConfigs) (*BotAPIInterface, error) {
	if interfaceCreated {
		return nil, &errs.BotInterfaceAlreadyCreated{}
	}
	interfaceCreated = true
	ch := make(chan *objs.Update)
	ch3 := make(chan *objs.ChatUpdate)
	temp := &BotAPIInterface{botConfigs: botCfg, updateChannel: &ch, chatUpadateChannel: &ch3}
	return temp, nil
}
