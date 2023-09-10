package telego

import objs "github.com/SakoDroid/telego/objects"

// forumTopicManager is a special object for managing forum topics
type forumTopicManager struct {
	bot             *Bot
	messageThreadId int
	chatId          int
	chatIdString    string
}

// Edit edits name and icon of a topic in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have can_manage_topics administrator rights, unless it is the creator of the topic. Returns True on success.
func (f *forumTopicManager) Edit(name, iconCustomEmojiId string) (*objs.Result[bool], error) {
	return f.bot.apiInterface.EditForumTopic(f.chatId, f.chatIdString, name, iconCustomEmojiId, f.messageThreadId)
}

// Close closes an open topic in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights, unless it is the creator of the topic. Returns True on success.
func (f *forumTopicManager) Close() (*objs.Result[bool], error) {
	return f.bot.apiInterface.CloseForumTopic(f.chatId, f.chatIdString, f.messageThreadId)
}

// Reopen reopens a closed topic in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights, unless it is the creator of the topic. Returns True on success.
func (f *forumTopicManager) Reopen() (*objs.Result[bool], error) {
	return f.bot.apiInterface.ReopenForumTopic(f.chatId, f.chatIdString, f.messageThreadId)
}

// Delete deletes a forum topic along with all its messages in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have the can_delete_messages administrator rights. Returns True on success.
func (f *forumTopicManager) Delete() (*objs.Result[bool], error) {
	return f.bot.apiInterface.DeleteForumTopic(f.chatId, f.chatIdString, f.messageThreadId)
}

// UnpinAllMesages clears the list of pinned messages in a forum topic. The bot must be an administrator in the chat for this to work and must have the can_pin_messages administrator right in the supergroup. Returns True on success.
func (f *forumTopicManager) UnpinAllMesages() (*objs.Result[bool], error) {
	return f.bot.apiInterface.UnpinAllForumTopicMessages(f.chatId, f.chatIdString, f.messageThreadId)
}
