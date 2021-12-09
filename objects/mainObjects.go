package objects

type Update struct {
	/*The update's unique identifier. Update identifiers start from a certain positive number and increase sequentially. This ID becomes especially handy if you're using Webhooks, since it allows you to ignore repeated updates or to restore the correct update sequence, should they get out of order. If there are no new updates for at least a week, then identifier of the next update will be chosen randomly instead of sequentially.*/
	Update_id int `json:"update_id"`
	/*Optional. New incoming message of any kind — text, photo, sticker, etc.*/
	Message Message `json:"message,omitempty"`
	/*Optional. New version of a message that is known to the bot and was edited*/
	Edited_message Message `json:"edited_message,omitempty"`
}

type Message struct {
}

type User struct {
	/*Unique identifier for this user or bot. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier.*/
	Id int `json:"id"`
	/*True, if this user is a bot*/
	IsBot bool `json:"is_bot"`
	/*User's or bot's first name*/
	FirstName string `json:"first_name"`
	/*Optional. User's or bot's last name*/
	Lastname string `json:"last_name,omitempty"`
	/*Optional. User's or bot's username*/
	Username string `json:"username,omitempty"`
	/*Optional. IETF language tag of the user's language*/
	LanguageCode string `json:"language_code,omitempty"`
	/*Optional. True, if the bot can be invited to groups.*/
	CanJoinGroups bool `json:"can_join_groups,omitempty"`
	/*Optional. True, if privacy mode is disabled for the bot.*/
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages,omitempty"`
	/*Optional. True, if the bot supports inline queries.*/
	SuportsInlineQueries bool `json:"supports_inline_queries,omitempty"`
}

/*Upon receiving a message with this object, Telegram clients will display a reply interface to the user (act as if the user has selected the bot's message and tapped 'Reply'). This can be extremely useful if you want to create user-friendly step-by-step interfaces without having to sacrifice privacy mode.*/
type ForceReply struct {
	/*Shows reply interface to the user, as if they manually selected the bot's message and tapped 'Reply'*/
	ForceReply bool `json:"force_reply"`
	/*Optional. The placeholder to be shown in the input field when the reply is active; 1-64 characters*/
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	/*Optional. Use this parameter if you want to force reply from specific users only. Targets: 1) users that are @mentioned in the text of the Message object; 2) if the bot's message is a reply (has reply_to_message_id), sender of the original message.*/
	Selective bool `json:"selective,omitempty"`
}

/*Contains information about why a request was unsuccessful.*/
type ResponseParameters struct {
	/*Optional. The group has been migrated to a supergroup with the specified identifier. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.*/
	MigrateToChatId int `json:"migrate_to_chat_id"`
	/*ptional. In case of exceeding flood control, the number of seconds left to wait before the request can be repeated*/
	RetryAfter int `json:"retry_after"`
}

/*Contains information about the current status of a webhook.*/
type Webhookinfo struct {
	/*Webhook URL, may be empty if webhook is not set up*/
	URL string `json:"url"`
	/*True, if a custom certificate was provided for webhook certificate checks*/
	HasCustomCertificate bool `json:"has_custom_certificate"`
	/*Number of updates awaiting delivery*/
	PendingUpdateCount int `json:"pending_update_count"`
	/*Optional. Currently used webhook IP address*/
	IPAddress string `json:"ip_address,omitempty"`
	/*Optional. Unix time for the most recent error that happened when trying to deliver an update via webhook*/
	LastErrorDate int `json:"last_error_date,omitempty"`
	/*Optional. Error message in human-readable format for the most recent error that happened when trying to deliver an update via webhook*/
	LastErrorMessage string `json:"last_error_message,omitempty"`
	/*Optional. Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery*/
	MaxConnection int `json:"max_connections,omitempty"`
	/*Optional. A list of update types the bot is subscribed to. Defaults to all update types except chat_member*/
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}
