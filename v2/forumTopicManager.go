package telego

import objs "github.com/SakoDroid/telego/v2/objects"

// ForumTopicManager is a special object for managing forum topics
type ForumTopicManager struct {
	bot             *Bot
	messageThreadId int
	chatId          int
	chatIdString    string
}

// Edit edits name and icon of a topic in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have can_manage_topics administrator rights, unless it is the creator of the topic. Returns True on success.
func (f *ForumTopicManager) Edit(name, iconCustomEmojiId string) (*objs.Result[bool], error) {
	return f.bot.apiInterface.EditForumTopic(f.chatId, f.chatIdString, name, iconCustomEmojiId, f.messageThreadId)
}

// Close closes an open topic in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights, unless it is the creator of the topic. Returns True on success.
func (f *ForumTopicManager) Close() (*objs.Result[bool], error) {
	return f.bot.apiInterface.CloseForumTopic(f.chatId, f.chatIdString, f.messageThreadId)
}

// Reopen reopens a closed topic in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have the can_manage_topics administrator rights, unless it is the creator of the topic. Returns True on success.
func (f *ForumTopicManager) Reopen() (*objs.Result[bool], error) {
	return f.bot.apiInterface.ReopenForumTopic(f.chatId, f.chatIdString, f.messageThreadId)
}

// Delete deletes a forum topic along with all its messages in a forum supergroup chat. The bot must be an administrator in the chat for this to work and must have the can_delete_messages administrator rights. Returns True on success.
func (f *ForumTopicManager) Delete() (*objs.Result[bool], error) {
	return f.bot.apiInterface.DeleteForumTopic(f.chatId, f.chatIdString, f.messageThreadId)
}

// UnpinAllMesages clears the list of pinned messages in a forum topic. The bot must be an administrator in the chat for this to work and must have the can_pin_messages administrator right in the supergroup. Returns True on success.
func (f *ForumTopicManager) UnpinAllMesages() (*objs.Result[bool], error) {
	return f.bot.apiInterface.UnpinAllForumTopicMessages(f.chatId, f.chatIdString, f.messageThreadId)
}
