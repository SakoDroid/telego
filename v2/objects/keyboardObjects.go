package objects

type ReplyMarkup interface {
	blah()
}

/*This object represents a custom keyboard with reply options*/
type ReplyKeyboardMarkup struct {
	/*Array of button rows, each represented by an Array of KeyboardButton objects*/
	Keyboard [][]*KeyboardButton `json:"keyboard"`
	/*Optional. Requests clients to always show the keyboard when the regular keyboard is hidden. Defaults to false, in which case the custom keyboard can be hidden and opened with a keyboard icon.*/
	IsPersistent bool `json:"is_persistent"`
	/*Optional. Requests clients to resize the keyboard vertically for optimal fit (e.g., make the keyboard smaller if there are just two rows of buttons). Defaults to false, in which case the custom keyboard is always of the same height as the app's standard keyboard.*/
	ResizeKeyboard bool `json:"resize_keyboard"`
	/*Optional. Requests clients to hide the keyboard as soon as it's been used. The keyboard will still be available, but clients will automatically display the usual letter-keyboard in the chat – the user can press a special button in the input field to see the custom keyboard again. Defaults to false*/
	OneTimeKeyboard bool `json:"one_tijme_keyboard"`
	/*Optional. The placeholder to be shown in the input field when the keyboard is active; 1-64 characters*/
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	/*Optional. Use this parameter if you want to show the keyboard to specific users only. Targets: 1) users that are @mentioned in the text of the Message object; 2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.

	Example: A user requests to change the bot's language, bot replies to the request with a keyboard to select the new language. Other users in the group don't see the keyboard.*/
	Selective bool `json:"selective"`
}

func (rm *ReplyKeyboardMarkup) blah() {}

/*This object represents one button of the reply keyboard. For simple text buttons String can be used instead of this object to specify text of the button. Optional fields request_contact, request_location, and request_poll are mutually exclusive*/
type KeyboardButton struct {
	/*Text of the button. If none of the optional fields are used, it will be sent as a message when the button is pressed*/
	Text string `json:"text"`
	/*Optional. If specified, pressing the button will open a list of suitable users. Tapping on any user will send their identifier to the bot in a “user_shared” service message. Available in private chats only.*/
	RequestUser *KeyboardButtonRequestUser `json:"request_user,omitempty"`
	/*Optional. If specified, pressing the button will open a list of suitable chats. Tapping on a chat will send its identifier to the bot in a “chat_shared” service message. Available in private chats only.*/
	RequestChat *KeyboardButtonRequestChat `json:"request_chat,omitempty"`
	/*Optional. If True, the user's phone number will be sent as a contact when the button is pressed. Available in private chats only
	Note: request_contact and request_location options will only work in Telegram versions released after 9 April, 2016. Older clients will display unsupported message.*/
	RequestContact bool `json:"request_contact"`
	/*Optional. If True, the user's current location will be sent when the button is pressed. Available in private chats only*/
	RequestLocation bool `json:"request_location"`
	/*Optional. If specified, the user will be asked to create a poll and send it to the bot when the button is pressed. Available in private chats only.
	Note: request_poll option will only work in Telegram versions released after 23 January, 2020. Older clients will display unsupported message.*/
	RequestPoll *KeyboardButtonPollType `json:"request_poll,omitempty"`
	/*A web app info (url) for launching a web app
	If specified, the described Web App will be launched when the button is pressed. The Web App will be able to send a “web_app_data” service message. Available in private chats only.*/
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

/*This object represents type of a poll, which is allowed to be created and sent when the corresponding button is pressed.*/
type KeyboardButtonPollType struct {
	/*Optional. If quiz is passed, the user will be allowed to create only polls in the quiz mode. If regular is passed, only regular polls will be allowed. Otherwise, the user will be allowed to create a poll of any type.*/
	Type string `json:"type,omitempty"`
}

/*This object defines the criteria used to request a suitable user. The identifier of the selected user will be shared with the bot when the corresponding button is pressed.*/
type KeyboardButtonRequestUser struct {
	/*Signed 32-bit identifier of the request, which will be received back in the UserShared object. Must be unique within the message*/
	RequestId int `json:"request_id"`
	/*Optional. Pass True to request a bot, pass False to request a regular user. If not specified, no additional restrictions are applied.*/
	UserIsBot bool `json:"user_is_bot"`
	/*Optional. Pass True to request a premium user, pass False to request a non-premium user. If not specified, no additional restrictions are applied.*/
	UserIsPremium bool `json:"user_is_premium"`
}

type KeyboardButtonRequestChat struct {
	/*Signed 32-bit identifier of the request, which will be received back in the ChatShared object. Must be unique within the message*/
	RequestId int `json:"request_id"`
	/*Pass True to request a channel chat, pass False to request a group or a supergroup chat.*/
	ChatIsChannel bool `json:"chat_is_channel"`
	/*Optional. Pass True to request a forum supergroup, pass False to request a non-forum chat. If not specified, no additional restrictions are applied.*/
	ChatIsForum bool `json:"chat_is_forum"`
	/*Optional. Pass True to request a supergroup or a channel with a username, pass False to request a chat without a username. If not specified, no additional restrictions are applied.*/
	ChatHasUsername bool `json:"chat_has_username"`
	/*Optional. Pass True to request a chat owned by the user. Otherwise, no additional restrictions are applied.*/
	ChatIsCreated bool `json:"chat_is_created"`
	/*Optional. A JSON-serialized object listing the required administrator rights of the user in the chat. The rights must be a superset of bot_administrator_rights. If not specified, no additional restrictions are applied.*/
	UserAdministratorRights *ChatAdministratorRights `json:"user_administrator_rights,omitempty"`
	/*Optional. A JSON-serialized object listing the required administrator rights of the bot in the chat. The rights must be a subset of user_administrator_rights. If not specified, no additional restrictions are applied.*/
	BotAdministratorRights *ChatAdministratorRights `json:"bot_administrator_rights,omitempty"`
	/*Optional. Pass True to request a chat with the bot as a member. Otherwise, no additional restrictions are applied.*/
	BotIsMemeber bool `json:"bot_is_member"`
}

/*Upon receiving a message with this object, Telegram clients will remove the current custom keyboard and display the default letter-keyboard. By default, custom keyboards are displayed until a new keyboard is sent by a bot. An exception is made for one-time keyboards that are hidden immediately after the user presses a button */
type ReplyKeyboardRemove struct {
	/*Requests clients to remove the custom keyboard (user will not be able to summon this keyboard; if you want to hide the keyboard from sight but keep it accessible, use one_time_keyboard in ReplyKeyboardMarkup)*/
	RemoveKeyboard bool `json:"remove_keyboard"`
	/*Optional. Use this parameter if you want to remove the keyboard for specific users only. Targets: 1) users that are @mentioned in the text of the Message object; 2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.

	Example: A user votes in a poll, bot returns confirmation message in reply to the vote and removes the keyboard for that user, while still showing the keyboard with poll options to users who haven't voted yet.*/
	Selective bool `json:"selective"`
}

func (rm *ReplyKeyboardRemove) blah() {}

/*This object represents an inline keyboard that appears right next to the message it belongs to.*/
type InlineKeyboardMarkup struct {
	/*Array of button rows, each represented by an Array of InlineKeyboardButton objects*/
	InlineKeyboard [][]*InlineKeyboardButton `json:"inline_keyboard"`
}

func (rm *InlineKeyboardMarkup) blah() {}

/*This object represents one button of an inline keyboard. You must use exactly one of the optional fields.*/
type InlineKeyboardButton struct {
	/*Label text on the button*/
	Text string `json:"text"`
	/*Optional. HTTP or tg:// url to be opened when the button is pressed. Links tg://user?id=<user_id> can be used to mention a user by their ID without using a username, if this is allowed by their privacy settings.*/
	URL string `json:"url,omitempty"`
	/*Optional. An HTTP URL used to automatically authorize the user. Can be used as a replacement for the Telegram Login Widget.*/
	LoginURL *LoginUrl `json:"login_url,omitempty"`
	/*Optional. Data to be sent in a callback query to the bot when button is pressed, 1-64 bytes*/
	CallbackData string `json:"callback_data,omitempty"`
	/*Optional. If set, pressing the button will prompt the user to select one of their chats, open that chat and insert the bot's username and the specified inline query in the input field. Can be empty, in which case just the bot's username will be inserted.

	Note: This offers an easy way for users to start using your bot in inline mode when they are currently in a private chat with it. Especially useful when combined with switch_pm… actions – in this case the user will be automatically returned to the chat they switched from, skipping the chat selection screen.*/
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`
	/*Optional. If set, pressing the button will insert the bot's username and the specified inline query in the current chat's input field. Can be empty, in which case only the bot's username will be inserted.

	This offers a quick way for the user to open your bot in inline mode in the same chat – good for selecting something from multiple options.*/
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`
	//Optional. If set, pressing the button will prompt the user to select one of their chats of the specified type, open that chat and insert the bot's username and the specified inline query in the input field
	SwitchInlineQueryChosenChat *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`
	/*Optional. Description of the game that will be launched when the user presses the button.

	NOTE: This type of button must always be the first button in the first row.*/
	CallbackGame *CallbackGame `json:"callback_game,omitempty"`
	/*Optional. Specify True, to send a Pay button.

	NOTE: This type of button must always be the first button in the first row and can only be used in invoice messages.*/
	Pay bool `json:"pay,omitempty"`
	/*A web app info (url) for launching a web app
	Description of the Web App that will be launched when the user presses the button. The Web App will be able to send an arbitrary message on behalf of the user using the method answerWebAppQuery. Available only in private chats between a user and the bot.
	*/
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

/*This object represents a parameter of the inline keyboard button used to automatically authorize a user. Serves as a great replacement for the Telegram Login Widget when the user is coming from Telegram. All the user needs to do is tap/click a button and confirm that they want to log in.*/
type LoginUrl struct {
	/*An HTTP URL to be opened with user authorization data added to the query string when the button is pressed. If the user refuses to provide authorization data, the original URL without information about the user will be opened. The data added is the same as described in Receiving authorization data.

	NOTE: You must always check the hash of the received data to verify the authentication and the integrity of the data as described in Checking authorization.*/
	URL string `json:"url"`
	/*Optional. New text of the button in forwarded messages.*/
	ForwardText string `json:"forward_text,omitempty"`
	/*Optional. Username of a bot, which will be used for user authorization. See Setting up a bot for more details. If not specified, the current bot's username will be assumed. The url's domain must be the same as the domain linked with the bot. See Linking your domain to the bot for more details.*/
	BotUsername string `json:"bot_username,omitempty"`
	/*Optional. Pass True to request the permission for your bot to send messages to the user.*/
	RequestWriteAccess bool `json:"request_write_access"`
}

/*This object represents an incoming callback query from a callback button in an inline keyboard. If the button that originated the query was attached to a message sent by the bot, the field message will be present. If the button was attached to a message sent via the bot (in inline mode), the field inline_message_id will be present. Exactly one of the fields data or game_short_name will be present.*/
type CallbackQuery struct {
	/*Unique identifier for this query*/
	Id string `json:"id"`
	/*Sender*/
	From User `json:"from"`
	/*Optional. Message with the callback button that originated the query. Note that message content and message date will not be available if the message is too old*/
	Message Message `json:"message,omitempty"`
	/*Optional. Identifier of the message sent via the bot in inline mode, that originated the query.*/
	InlineMessageId string `json:"inline_message_id,omitempty"`
	/*Global identifier, uniquely corresponding to the chat to which the message with the callback button was sent. Useful for high scores in games.*/
	ChatInstance string `json:"chat_instance,omitempty"`
	/*Optional. Data associated with the callback button. Be aware that a bad client can send arbitrary data in this field.*/
	Data string `json:"data,omitempty"`
	/*Optional. Short name of a Game to be returned, serves as the unique identifier for the game*/
	GameShortName string `json:"game_short_name,omitempty"`
}

// This object represents an inline button that switches the current user to inline mode in a chosen chat, with an optional default inline query.
type SwitchInlineQueryChosenChat struct {
	//Optional. The default inline query to be inserted in the input field. If left empty, only the bot's username will be inserted
	Query string `json:"query,omitempty"`
	//Optional. True, if private chats with users can be chosen
	AllowUserChats bool `json:"allow_user_chats"`
	//Optional. True, if private chats with bots can be chosen
	AllowBotChats bool `json:"allow_bot_chats"`
	//Optional. True, if group and supergroup chats can be chosen
	AllowGroupChats bool `json:"allow_group_chats"`
	//Optional. True, if channel chats can be chosen
	AllowChannelChats bool `json:"allow_channel_chats"`
}
