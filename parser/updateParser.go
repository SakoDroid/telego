package parser

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SakoDroid/telego/configs"
	errs "github.com/SakoDroid/telego/errors"
	"github.com/SakoDroid/telego/logger"
	objs "github.com/SakoDroid/telego/objects"
)

//ParseUpdate parses the received update and returns the last update offset.
func ParseUpdate(body []byte, uc *chan *objs.Update, cu *chan *objs.ChatUpdate, cfg *configs.BotConfigs) (int, error) {
	def := &objs.DefaultResult{}
	err2 := json.Unmarshal(body, def)
	if err2 != nil {
		return 0, err2
	}
	if !def.Ok {
		return 0, &errs.MethodNotSentError{Method: "getUpdates", Reason: "server returned false for \"ok\" field."}
	}
	ur := &objs.UpdateResult{}
	err := json.Unmarshal(body, ur)
	if err != nil {
		return 0, err
	}
	return parse(ur, uc, cu, cfg)
}

func parse(ur *objs.UpdateResult, uc *chan *objs.Update, cu *chan *objs.ChatUpdate, cfg *configs.BotConfigs) (int, error) {
	lastOffset := 0
	for _, val := range ur.Result {
		if val.Update_id > lastOffset {
			lastOffset = val.Update_id
		}
		ParseSingleUpdate(val, uc, cu, cfg)
	}
	return lastOffset, nil
}

//ParseSingleUpdate processes the given update object.
func ParseSingleUpdate(up *objs.Update, uc *chan *objs.Update, cu *chan *objs.ChatUpdate, cfg *configs.BotConfigs) {
	userId, isUserBlocked := isUserBlocked(up, cfg)
	if !isUserBlocked {
		logger.Log("Update", "\t\t\t\t", up.GetType(), "Parsed", logger.HEADER, logger.OKCYAN, logger.OKGREEN)
		if !checkHandlers(up) && !processChat(up, cu) {
			*uc <- up
		}
	} else {
		logger.Log("Update", "\t\t\t\t", up.GetType(), fmt.Sprintf("User %d is blocked", userId), logger.HEADER, logger.OKCYAN, logger.FAIL)
	}
}

func processChat(update *objs.Update, chatUpdateChannel *chan *objs.ChatUpdate) bool {
	var chat *objs.Chat
	switch {
	case update.Message != nil:
		chat = update.Message.Chat
	case update.EditedMessage != nil:
		chat = update.EditedMessage.Chat
	case update.ChannelPost != nil:
		chat = update.ChannelPost.Chat
	case update.EditedChannelPost != nil:
		chat = update.EditedChannelPost.Chat
	case update.MyChatMember != nil:
		chat = update.MyChatMember.Chat
	case update.ChatMember != nil:
		chat = update.ChatMember.Chat
	case update.ChatJoinRequest != nil:
		chat = update.ChatJoinRequest.Chat
	case update.CallbackQuery != nil:
		chat = update.CallbackQuery.Message.Chat
	}
	if chat == nil {
		return false
	}
	*chatUpdateChannel <- createChatUpdate(chat, update)
	return true
}

func createChatUpdate(chat *objs.Chat, update *objs.Update) *objs.ChatUpdate {
	out := objs.ChatUpdate{Update: update}
	if chat.Type == "channel" {
		out.ChatId = chat.Username
	} else {
		out.ChatId = strconv.Itoa(chat.Id)
	}
	return &out
}

func isUserBlocked(up *objs.Update, cfg *configs.BotConfigs) (int, bool) {
	switch up.GetType() {
	case "message":
		return checkBlocked(up.Message.From, cfg)
	case "edited_message":
		return checkBlocked(up.EditedMessage.From, cfg)
	case "inline_query":
		return checkBlocked(up.InlineQuery.From, cfg)
	case "chosen_inline_result":
		return checkBlocked(&up.ChosenInlineResult.From, cfg)
	case "callback_query":
		return checkBlocked(&up.CallbackQuery.From, cfg)
	case "shipping_query":
		return checkBlocked(up.ShippingQuery.From, cfg)
	case "pre_checkout_query":
		return checkBlocked(up.PreCheckoutQuery.From, cfg)
	case "poll_answer":
		return checkBlocked(up.PollAnswer.User, cfg)
	default:
		return 0, false
	}
}

func checkBlocked(user *objs.User, cfg *configs.BotConfigs) (int, bool) {
	for _, us := range cfg.BlockedUsers {
		if us.UserID == user.Id {
			return us.UserID, true
		}
	}
	return 0, false
}
