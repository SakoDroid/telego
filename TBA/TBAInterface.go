package TBA

import (
	"encoding/json"
	"errors"
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
				continue
			}
			err = bai.parseUpdateresults(res)
			if err != nil {
				logger.Logger.Println("Error parsing the result of the update. " + err.Error())
			}
		}
		time.Sleep(bai.botConfigs.UpdateConfigs.UpdateFrequency)
	}
}

func (bai *BotAPIInterface) parseUpdateresults(body []byte) error {
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

		return bai.SendCustom("sendMessage", args)
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
		return bai.SendCustom("forwardMessage", fm)
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString or fromChatIdInt or fromChatIdString", MethodName: "sendMessage"}
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
		return bai.SendCustom("copyMessage", cp)
	} else {
		return nil, &errs.RequiredArgumentError{ArgName: "chatIdInt or chatIdString or fromChatIdInt or fromChatIdString", MethodName: "sendMessage"}
	}
}

func (bai *BotAPIInterface) SendCustom(methdName string, args objs.MethodArguments) (*objs.SendMethodsResult, error) {
	cl := httpSenderClient{botApi: bai.botConfigs.BotAPI, apiKey: bai.botConfigs.APIKey}
	res, err2 := cl.sendHttpReqJson(methdName, args)
	if err2 != nil {
		return nil, err2
	}
	msg := &objs.SendMethodsResult{}
	err3 := json.Unmarshal(res, msg)
	if err3 != nil {
		return nil, err3
	}
	return msg, nil
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