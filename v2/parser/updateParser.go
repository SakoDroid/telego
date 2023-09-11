package parser

import (
	"fmt"
	"strconv"

	"github.com/SakoDroid/telego/v2/configs"
	"github.com/SakoDroid/telego/v2/logger"
	objs "github.com/SakoDroid/telego/v2/objects"
)

// ExecuteChain executes the chained middlewares
func ExecuteChain(up *objs.Update) {
	middlewares.executeChain(up)
}

// GetUpdateParserMiddleware returns a middleware that processes the given update object.
func GetUpdateParserMiddleware(uc *chan *objs.Update, cu *chan *objs.ChatUpdate, cfg *configs.BotConfigs) func(up *objs.Update, next func()) {
	//next is not called because this middleware is always the last middleware.
	return func(up *objs.Update, next func()) {
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
