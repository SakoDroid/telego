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
	cfgs "github.com/SakoDroid/telebot/configs"
	logger "github.com/SakoDroid/telebot/logger"
	objs "github.com/SakoDroid/telebot/objects"
)

var interfaceCreated = false

type BotAPIInterface struct {
	botConfigs           *cfgs.BotConfigs
	updateRoutineRunning bool
	updateChannel        *chan *objs.Update
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
	def := &objs.DefaultResult{}
	err2 := json.Unmarshal(body, def)
	if err2 != nil {
		return err2
	}
	if !def.Ok {
		return &errs.MethodNotSentError{Method: "getUpdates", Reason: "server returned false for \"ok\" field."}
	}
	ur := &objs.UpdateResult{}
	err := json.Unmarshal(body, ur)
	if err != nil {
		return err
	}
	if !ur.Ok {
		return &errs.UpdateNotOk{Offset: bai.lastOffset}
	}
	for _, val := range ur.Result {
		if val.Update_id > bai.lastOffset {
			bai.lastOffset = val.Update_id
		}
		(*bai.updateChannel) <- &val
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
	temp := &BotAPIInterface{botConfigs: botCfg, updateChannel: &ch}
	return temp, nil
}
