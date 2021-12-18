package TBA

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	errs "github.com/SakoDroid/telebot/Errors"
	up "github.com/SakoDroid/telebot/Parser"
	cfgs "github.com/SakoDroid/telebot/configs"
	logger "github.com/SakoDroid/telebot/logger"
	objs "github.com/SakoDroid/telebot/objects"
)

var interfaceCreated = false

type BotAPIInterface struct {
	botConfigs           *cfgs.BotConfigs
	updateRoutineRunning bool
	updateChannel        *chan *objs.Update
	pollUpdateChannel    *chan *objs.Poll
	updateRoutineChannel chan bool
	lastOffset           int
}

/*Starts the update routine to receive updates from api sever*/
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

/*Stops the update routine*/
func (bai *BotAPIInterface) StopUpdateRoutine() {
	if bai.updateRoutineRunning {
		bai.updateRoutineRunning = false
		bai.updateRoutineChannel <- true
	}
}

/*Returns the update channel*/
func (bai *BotAPIInterface) GetUpdateChannel() *chan *objs.Update {
	return bai.updateChannel
}

/*Returns the poll update channel*/
func (bai *BotAPIInterface) GetPollUpdateChannel() *chan *objs.Poll {
	return bai.pollUpdateChannel
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
	of, err := up.ParseUpdate(body, bai.updateChannel, bai.pollUpdateChannel)
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

/*Sends a message to the user. chatIdInt is used for all chats but channles and chatidString is used for channels (in form of @channleusername) and only of them has be populated, otherwise ChatIdProblem error will be returned.
"chatId" and "text" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) SendMessage(chatIdInt int, chatIdString, text, parseMode string, entities []objs.MessageEntity, disable_web_page_preview, disable_notification, allow_sending_without_reply bool, reply_to_message_id int, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		def := bai.fixTheDefaultArguments(chatIdInt, reply_to_message_id, chatIdString, disable_notification, allow_sending_without_reply, reply_markup)
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

/*Forwards a message from a user or channel to a user or channel. If the source or destination (or both) of the forwarded message is a channel, only string chat ids should be given to the function, and if it is user only int chat ids should be given.
"chatId", "fromChatId" and "messageId" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) ForwardMessage(chatIdInt, fromChatIdInt int, chatIdString, fromChatIdString string, disableNotif bool, messageId int) (*objs.SendMethodsResult, error) {
	if (chatIdInt != 0 && chatIdString != "") && (fromChatIdInt != 0 && fromChatIdString != "") {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) && bai.isChatIdOk(fromChatIdInt, fromChatIdString) {
		fm := &objs.ForwardMessageArgs{
			DisableNotification: disableNotif,
			MessageId:           messageId,
		}
		if chatIdInt == 0 {
			bt, _ := json.Marshal(chatIdString)
			fm.ChatId = bt
		} else {
			bt, _ := json.Marshal(chatIdInt)
			fm.ChatId = bt
		}
		if fromChatIdInt == 0 {
			bt, _ := json.Marshal(fromChatIdString)
			fm.FromChatId = bt
		} else {
			bt, _ := json.Marshal(fromChatIdInt)
			fm.FromChatId = bt
		}
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

/*Sends a photo (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "photo" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) SendPhoto(chatIdInt int, chatIdString, photo string, photoFile *os.File, caption, parseMode string, reply_to_message_id int, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup, captionEntities []objs.MessageEntity) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendPhotoArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a video (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "video" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendVideo(chatIdInt int, chatIdString, video string, videoFile *os.File, caption, parseMode string, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply bool, captionEntities []objs.MessageEntity, duration int, supportsStreaming bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVideoArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends an audio (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "audio" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0,to ignore string arguments pass "")*/
func (bai *BotAPIInterface) SendAudio(chatIdInt int, chatIdString, audio string, audioFile *os.File, caption, parseMode string, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply bool, captionEntities []objs.MessageEntity, duration int, performer, title string, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendAudioArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a document (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "document" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendDocument(chatIdInt int, chatIdString, document string, documentFile *os.File, caption, parseMode string, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply bool, captionEntities []objs.MessageEntity, DisableContentTypeDetection bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendDocumentArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends an animation (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "animation" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendAnimation(chatIdInt int, chatIdString, animation string, animationFile *os.File, caption, parseMode string, width, height, duration int, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply bool, captionEntities []objs.MessageEntity, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendAnimationArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a voice (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "voice" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendVoice(chatIdInt int, chatIdString, voice string, voiceFile *os.File, caption, parseMode string, duration int, reply_to_message_id int, disable_notification, allow_sending_without_reply bool, captionEntities []objs.MessageEntity, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVoiceArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a video note (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "videoNote" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)
Note that sending video note by URL is not supported by telegram.*/
func (bai *BotAPIInterface) SendVideoNote(chatIdInt int, chatIdString, videoNote string, videoNoteFile *os.File, caption, parseMode string, length, duration int, reply_to_message_id int, thumb string, thumbFile *os.File, disable_notification, allow_sending_without_reply bool, captionEntities []objs.MessageEntity, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVideoNoteArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends an album of media (file,url,telegramId) to a channel (chatIdString) or a chat (chatIdInt)
"chatId" and "media" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendMediaGroup(chatIdInt int, chatIdString string, reply_to_message_id int, media []objs.InputMedia, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup, files ...*os.File) (*objs.SendMediaGroupMethodResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendMediaGroupArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a location to a channel (chatIdString) or a chat (chatIdInt)
"chatId","latitude" and "longitude" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendLocation(chatIdInt int, chatIdString string, latitude, longitude, horizontalAccuracy float32, livePeriod, heading, proximityAlertRadius, reply_to_message_id int, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendLocationArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Edits a live location sent to a channel (chatIdString) or a chat (chatIdInt)
"chatId","latitude" and "longitude" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) EditMessageLiveLocation(chatIdInt int, chatIdString, inlineMessageId string, messageId int, latitude, longitude, horizontalAccuracy float32, heading, proximityAlertRadius int, reply_markup objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
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
		if chatIdInt == 0 {
			bt, _ := json.Marshal(chatIdString)
			args.ChatId = bt
		} else {
			bt, _ := json.Marshal(chatIdInt)
			args.ChatId = bt
		}
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

/*Stops a live location sent to a channel (chatIdString) or a chat (chatIdInt)
"chatId" argument is required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) StopMessageLiveLocation(chatIdInt int, chatIdString, inlineMessageId string, messageId int, replyMarkup *objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.StopMessageLiveLocationArgs{
			InlineMessageId: inlineMessageId,
			MessageId:       messageId,
			ReplyMarkup:     *replyMarkup,
		}
		if chatIdInt == 0 {
			bt, _ := json.Marshal(chatIdString)
			args.ChatId = bt
		} else {
			bt, _ := json.Marshal(chatIdInt)
			args.ChatId = bt
		}
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

/*Sends a venue to a channel (chatIdString) or a chat (chatIdInt)
"chatId","latitude","longitude","title" and "address" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendVenue(chatIdInt int, chatIdString string, latitude, longitude float32, title, address, fourSquareId, fourSquareType, googlePlaceId, googlePlaceType string, reply_to_message_id int, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendVenueArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a contact to a channel (chatIdString) or a chat (chatIdInt)
"chatId","phoneNumber" and "firstName" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendContact(chatIdInt int, chatIdString, phoneNumber, firstName, lastName, vCard string, reply_to_message_id int, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendContactArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a poll to a channel (chatIdString) or a chat (chatIdInt)
"chatId","phoneNumber" and "firstName" arguments are required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendPoll(chatIdInt int, chatIdString, question string, options []string, isClosed, isAnonymous bool, pollType string, allowMultipleAnswers bool, correctOptionIndex int, explanation, explanationParseMode string, explanationEntities []objs.MessageEntity, openPeriod, closeDate int, reply_to_message_id int, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendPollArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a dice message to a channel (chatIdString) or a chat (chatIdInt)
"chatId" argument is required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendDice(chatIdInt int, chatIdString, emoji string, reply_to_message_id int, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendDiceArgs{
			DefaultSendMethodsArguments: bai.fixTheDefaultArguments(
				chatIdInt, reply_to_message_id, chatIdString, disable_notification,
				allow_sending_without_reply, reply_markup,
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

/*Sends a chat action message to a channel (chatIdString) or a chat (chatIdInt)
"chatId" argument is required. other arguments are optional for bot api. (to ignore int arguments, pass 0)*/
func (bai *BotAPIInterface) SendChatAction(chatIdInt int, chatIdString, chatAction string) (*objs.SendMethodsResult, error) {
	if chatIdInt != 0 && chatIdString != "" {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) {
		args := &objs.SendChatActionArgs{
			Action: chatAction,
		}
		if chatIdInt == 0 {
			bt, _ := json.Marshal(chatIdString)
			args.ChatId = bt
		} else {
			bt, _ := json.Marshal(chatIdInt)
			args.ChatId = bt
		}
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

/*Gets the user profile photos*/
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

/*Gets the file based on the given file id and returns the file object. */
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

/*Downloads a file from telegram servers and saves it into the given file.

This method closes the given file. If the file is nil, this method will create a file based on the name of the file stored in telegram servers.*/
func (bai *BotAPIInterface) DownloadFile(fileObject objs.File, file *os.File) error {
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

/*Bans a chat member*/
func (bai *BotAPIInterface) BanChatMember(chatIdInt int, chatIdString string, userId, untilDate int, revokeMessages bool) (*objs.LogicalResult, error) {
	args := &objs.BanChatMemberArgs{
		UserId:         userId,
		UntilDate:      untilDate,
		RevokeMessages: revokeMessages,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Unbans a chat member*/
func (bai *BotAPIInterface) UnbanChatMember(chatIdInt int, chatIdString string, userId int, onlyIfBanned bool) (*objs.LogicalResult, error) {
	args := &objs.UnbanChatMemberArgsArgs{
		UserId:       userId,
		OnlyIfBanned: onlyIfBanned,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Restricts a chat member*/
func (bai *BotAPIInterface) RestrictChatMember(chatIdInt int, chatIdString string, userId int, permissions objs.ChatPermissions, untilDate int) (*objs.LogicalResult, error) {
	args := &objs.RestrictChatMemberArgs{
		UserId:     userId,
		Permission: permissions,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Promotes a chat member*/
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
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Sets a custom title for the administrator.*/
func (bai *BotAPIInterface) SetChatAdministratorCustomTitle(chatIdInt int, chatIdString string, userId int, customTitle string) (*objs.LogicalResult, error) {
	args := &objs.SetChatAdministratorCustomTitleArgs{
		UserId:      userId,
		CustomTitle: customTitle,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Bans or unbans a channel in the group..*/
func (bai *BotAPIInterface) BanOrUnbanChatSenderChat(chatIdInt int, chatIdString string, senderChatId int, ban bool) (*objs.LogicalResult, error) {
	args := &objs.BanChatSenderChatArgs{
		SenderChatId: senderChatId,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Sets default permissions for all users in the chat.*/
func (bai *BotAPIInterface) SetChatPermissions(chatIdInt int, chatIdString string, permissions objs.ChatPermissions) (*objs.LogicalResult, error) {
	args := &objs.SetChatPermissionsArgs{
		Permissions: permissions,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Exports the chat invite link and returns the new invite link as string.*/
func (bai *BotAPIInterface) ExportChatInviteLink(chatIdInt int, chatIdString string) (*objs.StringResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Creates a new invite link for the chat.*/
func (bai *BotAPIInterface) CreateChatInviteLink(chatIdInt int, chatIdString, name string, expireDate, memberLimit int, createsJoinRequest bool) (*objs.ChatInviteLinkResult, error) {
	args := &objs.CreateChatInviteLinkArgs{
		Name:               name,
		ExpireDate:         expireDate,
		MemberLimit:        memberLimit,
		CreatesjoinRequest: createsJoinRequest,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Edits an existing invite link for the chat.*/
func (bai *BotAPIInterface) EditChatInviteLink(chatIdInt int, chatIdString, inviteLink, name string, expireDate, memberLimit int, createsJoinRequest bool) (*objs.ChatInviteLinkResult, error) {
	args := &objs.EditChatInviteLinkArgs{
		InviteLink:         inviteLink,
		Name:               name,
		ExpireDate:         expireDate,
		MemberLimit:        memberLimit,
		CreatesjoinRequest: createsJoinRequest,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Revokes the given invite link.*/
func (bai *BotAPIInterface) RevokeChatInviteLink(chatIdInt int, chatIdString, inviteLink string) (*objs.ChatInviteLinkResult, error) {
	args := &objs.RevokeChatInviteLinkArgs{
		InviteLink: inviteLink,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Approves a request from the given user to join the chat.*/
func (bai *BotAPIInterface) ApproveChatJoinRequest(chatIdInt int, chatIdString string, userId int) (*objs.LogicalResult, error) {
	args := &objs.ApproveChatJoinRequestArgs{
		UserId: userId,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Declines a request from the given user to join the chat.*/
func (bai *BotAPIInterface) DeclineChatJoinRequest(chatIdInt int, chatIdString string, userId int) (*objs.LogicalResult, error) {
	args := &objs.DeclineChatJoinRequestArgs{
		UserId: userId,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Sets the chat photo to given file.*/
func (bai *BotAPIInterface) SetChatPhoto(chatIdInt int, chatIdString string, file *os.File) (*objs.LogicalResult, error) {
	args := &objs.SetChatPhotoArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Deletes chat photo.*/
func (bai *BotAPIInterface) DeleteChatPhoto(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Sets the chat title.*/
func (bai *BotAPIInterface) SetChatTitle(chatIdInt int, chatIdString, title string) (*objs.LogicalResult, error) {
	args := &objs.SetChatTitleArgs{
		Title: title,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Sets the chat description.*/
func (bai *BotAPIInterface) SetChatDescription(chatIdInt int, chatIdString, descriptions string) (*objs.LogicalResult, error) {
	args := &objs.SetChatDescriptionArgs{
		Description: descriptions,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Pines the message in the chat.*/
func (bai *BotAPIInterface) PinChatMessage(chatIdInt int, chatIdString string, messageId int, disableNotification bool) (*objs.LogicalResult, error) {
	args := &objs.PinChatMessageArgs{
		MessageId:           messageId,
		DisableNotification: disableNotification,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Unpines the pinned message in the chat.*/
func (bai *BotAPIInterface) UnpinChatMessage(chatIdInt int, chatIdString string, messageId int) (*objs.LogicalResult, error) {
	args := &objs.UnpinChatMessageArgs{
		MessageId: messageId,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Unpines all the pinned messages in the chat.*/
func (bai *BotAPIInterface) UnpinAllChatMessages(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*The bot will leave the chat if this method is called.*/
func (bai *BotAPIInterface) LeaveChat(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*A Chat object containing the information of the chat will be returned*/
func (bai *BotAPIInterface) GetChat(chatIdInt int, chatIdString string) (*objs.ChatResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Returnes an array of ChatMember containing the informations of the chat administrators.*/
func (bai *BotAPIInterface) GetChatAdministrators(chatIdInt int, chatIdString string) (*objs.ChatAdministratorsResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Returns the number of the memebrs of the chat.*/
func (bai *BotAPIInterface) GetChatMemberCount(chatIdInt int, chatIdString string) (*objs.IntResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Returns the information of the member in a ChatMember object.*/
func (bai *BotAPIInterface) GetChatMember(chatIdInt int, chatIdString string, userId int) (*objs.DefaultResult, error) {
	args := &objs.GetChatMemberArgs{
		UserId: userId,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Sets the sticker set of the chat.*/
func (bai *BotAPIInterface) SetChatStickerSet(chatIdInt int, chatIdString, stickerSetName string) (*objs.LogicalResult, error) {
	args := &objs.SetChatStcikerSet{
		StickerSetName: stickerSetName,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Deletes the sticker set of the chat..*/
func (bai *BotAPIInterface) DeleteChatStickerSet(chatIdInt int, chatIdString string) (*objs.LogicalResult, error) {
	args := &objs.DefaultChatArgs{}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Answers a callback query*/
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

/*Sets the commands of the bot*/
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

/*Deletes the commands of the bot*/
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

/*Gets the commands of the bot*/
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

/*Edits the text of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessageText(chatIdInt int, chatIdString string, messageId int, inlineMessageId, text, parseMode string, entities []objs.MessageEntity, disableWebPagePreview bool, replyMakrup objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
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
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Edits the caption of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessageCaption(chatIdInt int, chatIdString string, messageId int, inlineMessageId, caption, parseMode string, captionEntities []objs.MessageEntity, replyMakrup objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
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
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Edits the media of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessageMedia(chatIdInt int, chatIdString string, messageId int, inlineMessageId string, media objs.InputMedia, replyMakrup objs.InlineKeyboardMarkup, file ...*os.File) (*objs.DefaultResult, error) {
	args := &objs.EditMessageMediaArgs{
		EditMessageDefaultArgs: objs.EditMessageDefaultArgs{
			MessageId:       messageId,
			InlineMessageId: inlineMessageId,
			ReplyMarkup:     replyMakrup,
		},
		Media: media,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Edits the reply makrup of the given message in the given chat.*/
func (bai *BotAPIInterface) EditMessagereplyMarkup(chatIdInt int, chatIdString string, messageId int, inlineMessageId string, replyMakrup objs.InlineKeyboardMarkup) (*objs.DefaultResult, error) {
	args := &objs.EditMessageReplyMakrupArgs{
		EditMessageDefaultArgs: objs.EditMessageDefaultArgs{
			MessageId:       messageId,
			InlineMessageId: inlineMessageId,
			ReplyMarkup:     replyMakrup,
		},
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Stops the poll.*/
func (bai *BotAPIInterface) StopPoll(chatIdInt int, chatIdString string, messageId int, replyMakrup objs.InlineKeyboardMarkup) (*objs.PollResult, error) {
	args := &objs.StopPollArgs{
		MessageId:   messageId,
		ReplyMarkup: replyMakrup,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Deletes the given message int the given chat.*/
func (bai *BotAPIInterface) DeleteMessage(chatIdInt int, chatIdString string, messageId int) (*objs.LogicalResult, error) {
	args := &objs.DeleteMessageArgs{
		MessageId: messageId,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
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

/*Sends a sticker to the given chat id.*/
func (bai *BotAPIInterface) SendSticker(chatIdInt int, chatIdString, sticker string, disableNotif, allowSendingWithoutreply bool, replyTo int, replyMarkup objs.ReplyMarkup, files ...*os.File) (*objs.SendMethodsResult, error) {
	args := &objs.SendStickerArgs{
		DefaultSendMethodsArguments: objs.DefaultSendMethodsArguments{
			DisableNotification:      disableNotif,
			AllowSendingWithoutReply: allowSendingWithoutreply,
			ReplyToMessageId:         replyTo,
			ReplyMarkup:              replyMarkup,
		},
		Sticker: sticker,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		args.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		args.ChatId = bt
	}
	res, err := bai.SendCustom("sendSticker", args, true, files...)
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

/*Gets the sticker set by the given name*/
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

/*Uploads the given file as an sticker on the telegram servers.*/
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

/*Creates a new sticker set with the given arguments*/
func (bai *BotAPIInterface) CreateNewStickerSet(userId int, name, title, pngSticker, tgsSticker, emojies string, containsMasks bool, maskPosition objs.MaskPosition, file *os.File) (*objs.LogicalResult, error) {
	args := &objs.CreateNewStickerSetArgs{
		UserId:        userId,
		Name:          name,
		Title:         title,
		PngSticker:    pngSticker,
		TgsSticker:    tgsSticker,
		Emojies:       emojies,
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

/*Adds a new sticker to the given set.*/
func (bai *BotAPIInterface) AddStickerToSet(userId int, name, pngSticker, tgsSticker, emojies string, maskPosition objs.MaskPosition, file *os.File) (*objs.LogicalResult, error) {
	args := &objs.AddStickerSetArgs{
		UserId:       userId,
		Name:         name,
		PngSticker:   pngSticker,
		TgsSticker:   tgsSticker,
		Emojies:      emojies,
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

/*Sets the position of a sticker in an sticker set*/
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

/*Deletes the given sticker from a set created by the bot*/
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

/*Sets the thumbnail for the given sticker*/
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

/*Copies a message from a user or channel and sends it to a user or channel. If the source or destination (or both) of the forwarded message is a channel, only string chat ids should be given to the function, and if it is user only int chat ids should be given.
"chatId", "fromChatId" and "messageId" arguments are required. other arguments are optional for bot api.*/
func (bai *BotAPIInterface) CopyMessage(chatIdInt, fromChatIdInt int, chatIdString, fromChatIdString string, messageId int, disableNotif bool, caption, parseMode string, replyTo int, allowSendingWihtoutReply bool, replyMarkUp objs.ReplyMarkup, captionEntities []objs.MessageEntity) (*objs.SendMethodsResult, error) {
	if (chatIdInt != 0 && chatIdString != "") && (fromChatIdInt != 0 && fromChatIdString != "") {
		return nil, &errs.ChatIdProblem{}
	}
	if bai.isChatIdOk(chatIdInt, chatIdString) && bai.isChatIdOk(fromChatIdInt, fromChatIdString) {
		fm := objs.ForwardMessageArgs{
			DisableNotification: disableNotif,
			MessageId:           messageId,
		}
		if chatIdInt == 0 {
			bt, _ := json.Marshal(chatIdString)
			fm.ChatId = bt
		} else {
			bt, _ := json.Marshal(chatIdInt)
			fm.ChatId = bt
		}
		if fromChatIdInt == 0 {
			bt, _ := json.Marshal(fromChatIdString)
			fm.FromChatId = bt
		} else {
			bt, _ := json.Marshal(fromChatIdInt)
			fm.FromChatId = bt
		}
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

/*Calls the given method on api server wtih the given arguments. "MP" options indicates that the request should be made in multipart/formadata form. If this method sends a file to the api server the "MP" option should be true*/
func (bai *BotAPIInterface) SendCustom(methdName string, args objs.MethodArguments, MP bool, files ...*os.File) ([]byte, error) {
	cl := httpSenderClient{botApi: bai.botConfigs.BotAPI, apiKey: bai.botConfigs.APIKey}
	var res []byte
	var err2 error
	if MP {
		res, err2 = cl.sendHttpReqMultiPart(methdName, args, files...)
	} else {
		res, err2 = cl.sendHttpReqJson(methdName, args)
	}
	if err2 != nil {
		return nil, err2
	}
	return res, nil
}

func (bai *BotAPIInterface) fixTheDefaultArguments(chatIdInt, reply_to_message_id int, chatIdString string, disable_notification, allow_sending_without_reply bool, reply_markup objs.ReplyMarkup) objs.DefaultSendMethodsArguments {
	def := objs.DefaultSendMethodsArguments{
		DisableNotification:      disable_notification,
		AllowSendingWithoutReply: allow_sending_without_reply,
	}
	if chatIdInt == 0 {
		bt, _ := json.Marshal(chatIdString)
		def.ChatId = bt
	} else {
		bt, _ := json.Marshal(chatIdInt)
		def.ChatId = bt
	}
	if reply_to_message_id != 0 {
		def.ReplyToMessageId = reply_to_message_id
	}
	if reply_markup != nil {
		def.ReplyMarkup = reply_markup
	}
	return def
}

/*This method returns an iterface to communicate with the bot api.
If the updateFrequency argument is not nil, the update routine begins automtically*/
func CreateInterface(botCfg *cfgs.BotConfigs) (*BotAPIInterface, error) {
	if interfaceCreated {
		return nil, &errs.BotInterfaceAlreadyCreated{}
	}
	interfaceCreated = true
	ch := make(chan *objs.Update)
	ch2 := make(chan *objs.Poll)
	temp := &BotAPIInterface{botConfigs: botCfg, updateChannel: &ch, pollUpdateChannel: &ch2}
	return temp, nil
}
