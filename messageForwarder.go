package telego

import (
	objs "github.com/SakoDroid/telego/objects"
)

//MessageForwarder is a tool for forwarding messages.
type MessageForwarder struct {
	bot                          *Bot
	disableNotif, protectContent bool
	messageId                    int
}

/*ForwardFromUserToUser forwards the given message from a user to another user. chatId is the user that message is being forwarded to and fromChatId is the user that message is being forwarded to.*/
func (mf *MessageForwarder) ForwardFromUserToUser(chatId, fromChatId int) (*objs.SendMethodsResult, error) {
	return mf.bot.apiInterface.ForwardMessage(chatId, fromChatId, "", "", mf.disableNotif, mf.protectContent, mf.messageId)
}

/*ForwardFromUserToChannel forwards the given message from a user to a channel. chatId is the channel that message is being forwarded to and fromChatId is the user that message is being forwarded to.*/
func (mf *MessageForwarder) ForwardFromUserToChannel(chatId string, fromChatId int) (*objs.SendMethodsResult, error) {
	return mf.bot.apiInterface.ForwardMessage(0, fromChatId, chatId, "", mf.disableNotif, mf.protectContent, mf.messageId)
}

/*ForwardFromChannelToUser forwards the given message from a channel to a user. chatId is the user that message is being forwarded to and fromChatId is the channel that message is being forwarded to.*/
func (mf *MessageForwarder) ForwardFromChannelToUser(chatId int, fromChatId string) (*objs.SendMethodsResult, error) {
	return mf.bot.apiInterface.ForwardMessage(chatId, 0, "", fromChatId, mf.disableNotif, mf.protectContent, mf.messageId)
}

/*ForwardFromChannelToChannel forwards the given message from a channel to another channel. chatId is the channel that message is being forwarded to and fromChatId is the channel that message is being forwarded to.*/
func (mf *MessageForwarder) ForwardFromChannelToChannel(chatId, fromChatId string) (*objs.SendMethodsResult, error) {
	return mf.bot.apiInterface.ForwardMessage(0, 0, chatId, fromChatId, mf.disableNotif, mf.protectContent, mf.messageId)
}
