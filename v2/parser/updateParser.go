package parser

import (
	"fmt"
	"strconv"

	"github.com/SakoDroid/telego/v2/configs"
	"github.com/SakoDroid/telego/v2/logger"
	objs "github.com/SakoDroid/telego/v2/objects"
)

type UpdateParser struct {
	uc                 *chan *objs.Update
	cu                 *chan *objs.ChatUpdate
	cfg                *configs.BotConfigs
	handlers           *handlerTree
	callbackHandlers   threadSafeMap[string, *callbackHandler]
	userSharedHandlers threadSafeMap[int, *chatRequestHandler]
	chatSharedHandlers threadSafeMap[int, *chatRequestHandler]
}

// ExecuteChain executes the chained middlewares
func (u *UpdateParser) ExecuteChain(up *objs.Update) {
	middlewares.executeChain(up)
}

// GetUpdateParserMiddleware returns a middleware that processes the given update object.
func (u *UpdateParser) GetUpdateParserMiddleware(uc *chan *objs.Update, cu *chan *objs.ChatUpdate, cfg *configs.BotConfigs) func(up *objs.Update, next func()) {
	//next is not called because this middleware is always the last middleware.
	return func(up *objs.Update, next func()) {
		userId, isUserBlocked := u.isUserBlocked(up, cfg)
		if !isUserBlocked {
			logger.Log("Update", "\t\t\t\t", up.GetType(), "Parsed", logger.HEADER, logger.OKCYAN, logger.OKGREEN)
			if !u.checkHandlers(up) && !u.processChat(up, cu) {
				*uc <- up
			}
		} else {
			logger.Log("Update", "\t\t\t\t", up.GetType(), fmt.Sprintf("User %d is blocked", userId), logger.HEADER, logger.OKCYAN, logger.FAIL)
		}
	}
}

func (u *UpdateParser) processChat(update *objs.Update, chatUpdateChannel *chan *objs.ChatUpdate) bool {
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
	*chatUpdateChannel <- u.createChatUpdate(chat, update)
	return true
}

func (u *UpdateParser) createChatUpdate(chat *objs.Chat, update *objs.Update) *objs.ChatUpdate {
	out := objs.ChatUpdate{Update: update}
	if chat.Type == "channel" {
		out.ChatId = chat.Username
	} else {
		out.ChatId = strconv.Itoa(chat.Id)
	}
	return &out
}

func (u *UpdateParser) isUserBlocked(up *objs.Update, cfg *configs.BotConfigs) (int, bool) {
	switch up.GetType() {
	case "message":
		return u.checkBlocked(up.Message.From, cfg)
	case "edited_message":
		return u.checkBlocked(up.EditedMessage.From, cfg)
	case "inline_query":
		return u.checkBlocked(up.InlineQuery.From, cfg)
	case "chosen_inline_result":
		return u.checkBlocked(&up.ChosenInlineResult.From, cfg)
	case "callback_query":
		return u.checkBlocked(&up.CallbackQuery.From, cfg)
	case "shipping_query":
		return u.checkBlocked(up.ShippingQuery.From, cfg)
	case "pre_checkout_query":
		return u.checkBlocked(up.PreCheckoutQuery.From, cfg)
	case "poll_answer":
		return u.checkBlocked(up.PollAnswer.User, cfg)
	default:
		return 0, false
	}
}

func (u *UpdateParser) checkBlocked(user *objs.User, cfg *configs.BotConfigs) (int, bool) {
	for _, us := range cfg.BlockedUsers {
		if us.UserID == user.Id {
			return us.UserID, true
		}
	}
	return 0, false
}

func (u *UpdateParser) AddMiddleWare(middleware func(update *objs.Update, next func())) {
	middlewares.addToBegin(middleware)
}

func CreateUpdateParser(uc *chan *objs.Update, cu *chan *objs.ChatUpdate, cfg *configs.BotConfigs) *UpdateParser {
	up := &UpdateParser{
		uc:                 uc,
		cu:                 cu,
		cfg:                cfg,
		handlers:           &handlerTree{},
		callbackHandlers:   threadSafeMap[string, *callbackHandler]{internal: make(map[string]*callbackHandler)},
		userSharedHandlers: threadSafeMap[int, *chatRequestHandler]{internal: make(map[int]*chatRequestHandler)},
		chatSharedHandlers: threadSafeMap[int, *chatRequestHandler]{internal: make(map[int]*chatRequestHandler)},
	}

	up.AddMiddleWare(
		func(update *objs.Update, next func()) {
			userId, isUserBlocked := up.isUserBlocked(update, cfg)
			if !isUserBlocked {
				logger.Log("Update", "\t\t\t\t", update.GetType(), "Parsed", logger.HEADER, logger.OKCYAN, logger.OKGREEN)
				if !up.checkHandlers(update) && !up.processChat(update, cu) {
					*uc <- update
				}
			} else {
				logger.Log("Update", "\t\t\t\t", update.GetType(), fmt.Sprintf("User %d is blocked", userId), logger.HEADER, logger.OKCYAN, logger.FAIL)
			}
		},
	)
	return up
}
