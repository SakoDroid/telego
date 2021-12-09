package objects

type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_tijme_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
}

type KeyboardButton struct {
	Text            string                 `json:"text"`
	RequestContact  bool                   `json:"request_contact,omitempty"`
	RequestLocation bool                   `json:"request_location,omitempty"`
	RequestPoll     KeyboardButtonPollType `json:"request_poll,omitempty"`
}

type KeyboardButtonPollType struct {
	Type string `json:"type,omitempty"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text                         string       `json:"text"`
	URL                          string       `json:"url,omitempty"`
	LoginURL                     LoginUrl     `json:"login_url,omitempty"`
	CallbackData                 string       `json:"callback_data,omitempty"`
	SwitchInlineQuery            string       `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string       `json:"switch_inline_query_current_chat,omitempty"`
	CallbackGame                 CallbackGame `json:"callbacl_game,omitempty"`
	Pay                          bool         `json:"pay,omitempty"`
}

type LoginUrl struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text,omitempty"`
	BotUsername        string `json:"bot_username,omitempty"`
	RequestWriteAccess bool   `json:"request_write_access,omitempty"`
}

type CallbackQuery struct {
	Id              string  `json:"id"`
	From            User    `json:"from"`
	Message         Message `json:"message,omitempty"`
	InlineMessageId string  `json:"inline_message_id,omitempty"`
	ChatInstance    string  `json:"chat_instance,omitempty"`
	Data            string  `json:"data,omitempty"`
	GameShortName   string  `json:"game_short_name,omitempty"`
}
