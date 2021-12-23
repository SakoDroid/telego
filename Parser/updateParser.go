package parser

import (
	"encoding/json"
	"strconv"

	errs "github.com/SakoDroid/telebot/Errors"
	objs "github.com/SakoDroid/telebot/objects"
)

//Parses the recived update and returns the last update offset.
func ParseUpdate(body []byte, uc *chan *objs.Update, cu *chan *objs.ChatUpdate) (int, error) {
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
	return parse(ur, uc, cu)
}

func parse(ur *objs.UpdateResult, uc *chan *objs.Update, cu *chan *objs.ChatUpdate) (int, error) {
	lastOffset := 0
	for _, val := range ur.Result {
		if val.Update_id > lastOffset {
			lastOffset = val.Update_id
		}
		if !processChat(val, cu) {
			*uc <- val
		}
	}
	return lastOffset, nil
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
	}
	if chat == nil {
		return false
	} else {
		*chatUpdateChannel <- createChatUpdate(chat, update)
	}
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
