package parser

import (
	"encoding/json"

	errs "github.com/SakoDroid/telebot/Errors"
	objs "github.com/SakoDroid/telebot/objects"
)

//Parses the recived update and returns the last update offset.
func ParseUpdate(body []byte, uc *chan *objs.Update, pu *chan *objs.Update, messageChannel *chan *objs.Message, editedMessageChannel *chan *objs.Message, channelPostChannel *chan *objs.Message, editedChannelPostChannel *chan *objs.Message, inlineQueryChannel *chan *objs.InlineQuery, chosenInlineResultChannel *chan *objs.ChosenInlineResult, callbackQueryChannel *chan *objs.CallbackQuery, shippingQueryChannel *chan *objs.ShippingQuery, preCheckoutQueryChannel *chan *objs.PreCheckoutQuery, myChatMemberChannel *chan *objs.ChatMemberUpdated, chatMemberChannel *chan *objs.ChatMemberUpdated, chatJoinRequestChannel *chan *objs.ChatJoinRequest) (int, error) {
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
	return parse(
		ur, uc, pu, messageChannel,
		editedMessageChannel,
		channelPostChannel,
		editedChannelPostChannel,
		inlineQueryChannel,
		chosenInlineResultChannel,
		callbackQueryChannel,
		shippingQueryChannel,
		preCheckoutQueryChannel,
		myChatMemberChannel,
		chatMemberChannel,
		chatJoinRequestChannel,
	)
}

func parse(ur *objs.UpdateResult, uc *chan *objs.Update, pu *chan *objs.Update, messageChannel *chan *objs.Message, editedMessageChannel *chan *objs.Message, channelPostChannel *chan *objs.Message, editedChannelPostChannel *chan *objs.Message, inlineQueryChannel *chan *objs.InlineQuery, chosenInlineResultChannel *chan *objs.ChosenInlineResult, callbackQueryChannel *chan *objs.CallbackQuery, shippingQueryChannel *chan *objs.ShippingQuery, preCheckoutQueryChannel *chan *objs.PreCheckoutQuery, myChatMemberChannel *chan *objs.ChatMemberUpdated, chatMemberChannel *chan *objs.ChatMemberUpdated, chatJoinRequestChannel *chan *objs.ChatJoinRequest) (int, error) {
	lastOffset := 0
	for _, val := range ur.Result {
		if val.Update_id > lastOffset {
			lastOffset = val.Update_id
		}
		processTheUpdate(
			val, uc, pu, messageChannel, editedMessageChannel, channelPostChannel,
			editedChannelPostChannel, inlineQueryChannel, chosenInlineResultChannel,
			callbackQueryChannel, shippingQueryChannel, preCheckoutQueryChannel,
			myChatMemberChannel, chatMemberChannel, chatJoinRequestChannel,
		)
	}
	return lastOffset, nil
}

func processTheUpdate(update *objs.Update, uc *chan *objs.Update, pu *chan *objs.Update, messageChannel *chan *objs.Message, editedMessageChannel *chan *objs.Message, channelPostChannel *chan *objs.Message, editedChannelPostChannel *chan *objs.Message, inlineQueryChannel *chan *objs.InlineQuery, chosenInlineResultChannel *chan *objs.ChosenInlineResult, callbackQueryChannel *chan *objs.CallbackQuery, shippingQueryChannel *chan *objs.ShippingQuery, preCheckoutQueryChannel *chan *objs.PreCheckoutQuery, myChatMemberChannel *chan *objs.ChatMemberUpdated, chatMemberChannel *chan *objs.ChatMemberUpdated, chatJoinRequestChannel *chan *objs.ChatJoinRequest) {
	switch {
	case messageChannel != nil && update.Message != nil:
		*messageChannel <- update.Message
	case editedMessageChannel != nil && update.EditedMessage != nil:
		*editedMessageChannel <- update.EditedMessage
	case channelPostChannel != nil && update.ChannelPost != nil:
		*channelPostChannel <- update.ChannelPost
	case editedChannelPostChannel != nil && update.EditedChannelPost != nil:
		*editedChannelPostChannel <- update.EditedChannelPost
	case inlineQueryChannel != nil && update.InlineQuery != nil:
		*inlineQueryChannel <- update.InlineQuery
	case chosenInlineResultChannel != nil && update.ChosenInlineResult != nil:
		*chosenInlineResultChannel <- update.ChosenInlineResult
	case callbackQueryChannel != nil && update.CallbackQuery != nil:
		*callbackQueryChannel <- update.CallbackQuery
	case shippingQueryChannel != nil && update.ShippingQuery != nil:
		*shippingQueryChannel <- update.ShippingQuery
	case preCheckoutQueryChannel != nil && update.PreCheckoutQuery != nil:
		*preCheckoutQueryChannel <- update.PreCheckoutQuery
	case update.Poll != nil:
		*pu <- update
	case myChatMemberChannel != nil && update.MyChatMember != nil:
		*myChatMemberChannel <- update.MyChatMember
	case chatMemberChannel != nil && update.ChatMember != nil:
		*chatMemberChannel <- update.ChatMember
	case chatJoinRequestChannel != nil && update.ChatJoinRequest != nil:
		*chatJoinRequestChannel <- update.ChatJoinRequest
	default:
		*uc <- update
	}
}
