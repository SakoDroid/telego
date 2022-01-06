package objects

/*This object represents an incoming update.
At most one of the optional parameters can be present in any given update.*/
type Update struct {
	/*The update's unique identifier. Update identifiers start from a certain positive number and increase sequentially. This ID becomes especially handy if you're using Webhooks, since it allows you to ignore repeated updates or to restore the correct update sequence, should they get out of order. If there are no new updates for at least a week, then identifier of the next update will be chosen randomly instead of sequentially.*/
	Update_id int `json:"update_id"`
	/*Optional. New incoming message of any kind — text, photo, sticker, etc.*/
	Message *Message `json:"message,omitempty"`
	/*Optional. New version of a message that is known to the bot and was edited*/
	EditedMessage *Message `json:"edited_message,omitempty"`
	/*Optional. New incoming channel post of any kind — text, photo, sticker, etc.*/
	ChannelPost *Message `json:"channel_post,omitempty"`
	/*Optional. New version of a channel post that is known to the bot and was edited*/
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`
	/*Optional. New incoming inline query*/
	InlineQuery *InlineQuery `json:"inline_query,omitempty"`
	/*Optional. The result of an inline query that was chosen by a user and sent to their chat partner. */
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	/*Optional. New incoming callback query*/
	CallbackQuery *CallbackQuery `json:"callback_query,omitempty"`
	/*Optional. New incoming shipping query. Only for invoices with flexible price*/
	ShippingQuery *ShippingQuery `json:"shipping_query,omitempty"`
	/*Optional. New incoming pre-checkout query. Contains full information about checkout*/
	PreCheckoutQuery *PreCheckoutQuery `json:"pre_checkout_query,omitempty"`
	/*Optional. New poll state. Bots receive only updates about stopped polls and polls, which are sent by the bot*/
	Poll *Poll `json:"poll,omitempty"`
	/*Optional. A user changed their answer in a non-anonymous poll. Bots receive new votes only in polls that were sent by the bot itself.*/
	PollAnswer *PollAnswer `json:"poll_answer,omitempty"`
	/*Optional. The bot's chat member status was updated in a chat. For private chats, this update is received only when the bot is blocked or unblocked by the user.*/
	MyChatMember *ChatMemberUpdated `json:"my_chat_member,omitempty"`
	/*Optional. A chat member's status was updated in a chat. The bot must be an administrator in the chat and must explicitly specify “chat_member” in the list of allowed_updates to receive these updates.*/
	ChatMember *ChatMemberUpdated `json:"chat_member,omitempty"`
	/*Optional. A request to join the chat has been sent. The bot must have the can_invite_users administrator right in the chat to receive these updates.*/
	ChatJoinRequest *ChatJoinRequest `json:"chat_join_request,omitempty"`
}

/*Returnes the populated field of this update*/
func (u *Update) GetType() string {
	if u.Message != nil {
		return "message"
	}
	if u.EditedMessage != nil {
		return "edited_message"
	}
	if u.ChannelPost != nil {
		return "channel_post"
	}
	if u.EditedChannelPost != nil {
		return "edited_channel_post"
	}
	if u.InlineQuery != nil {
		return "inline_query"
	}
	if u.ChosenInlineResult != nil {
		return "chosen_inline_result"
	}
	if u.CallbackQuery != nil {
		return "callback_query"
	}
	if u.ShippingQuery != nil {
		return "shipping_query"
	}
	if u.PreCheckoutQuery != nil {
		return "pre_checkout_query"
	}
	if u.Poll != nil {
		return "poll"
	}
	if u.PollAnswer != nil {
		return "poll_answer"
	}
	if u.MyChatMember != nil {
		return "my_chat_member"
	}
	if u.ChatMember != nil {
		return "chat_member"
	}
	if u.ChatJoinRequest != nil {
		return "chat_join_request"
	}
	return ""
}

/*This object represents a message.*/
type Message struct {
	/*Unique message identifier inside this chat*/
	MessageId int `json:"message_id"`
	/*Optional. Sender, empty for messages sent to channels*/
	From *User `json:"from,omitempty"`
	/*Optional. Sender of the message, sent on behalf of a chat. The channel itself for channel messages. The supergroup itself for messages from anonymous group administrators. The linked channel for messages automatically forwarded to the discussion group*/
	SenderChat *Chat `json:"sender_chat,omitempty"`
	/*Date the message was sent in Unix time*/
	Date int `json:"date"`
	/*Conversation the message belongs to*/
	Chat *Chat `json:"chat"`
	/*Optional. For forwarded messages, sender of the original message*/
	ForwardFrom *User `json:"forward_from,omitempty"`
	/*Optional. For messages forwarded from channels or from anonymous administrators, information about the original sender chat*/
	ForwardFromChat *Chat `json:"forward_from_chat,omitempty"`
	/*Optional. For messages forwarded from channels, identifier of the original message in the channel*/
	ForwardFromMessageId int `json:"forward_from_message_id,omitempty"`
	/*Optional. For messages forwarded from channels, signature of the post author if present*/
	ForwardSignature string `json:"forward_signature,omitempty"`
	/*Optional. Sender's name for messages forwarded from users who disallow adding a link to their account in forwarded messages*/
	ForwardSenderName string `json:"forward_sender_name,omitempty"`
	/*Optional. For forwarded messages, date the original message was sent in Unix time*/
	ForwardDate int `json:"forward_date,omitempty"`
	/*Optional. True, if the message is a channel post that was automatically forwarded to the connected discussion group*/
	IsAutomaticForward bool `json:"is_automatic_forward,omitempty"`
	/*Optional. For replies, the original message. Note that the Message object in this field will not contain further reply_to_message fields even if it itself is a reply.*/
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`
	/*Optional. Bot through which the message was sent*/
	ViaBot User `json:"via_bot,omitempty"`
	/*Optional. Date the message was last edited in Unix time*/
	EditDate int `json:"edit_date,omitempty"`
	/*Optional. True, if the message can't be forwarded*/
	HasProtectedContent bool `json:"has_protected_content,omitempty"`
	/*Optional. The unique identifier of a media message group this message belongs to*/
	MediaGroupId string `json:"media_group_id,omitempty"`
	/*Optional. Signature of the post author for messages in channels, or the custom title of an anonymous group administrator*/
	AuthorSignature string `json:"author_signature,omitempty"`
	/*Optional. For text messages, the actual UTF-8 text of the message, 0-4096 characters*/
	Text string `json:"text,omitempty"`
	/*Optional. For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text*/
	Entities []MessageEntity `json:"entities,omitempty"`
	/*Optional. Message is an animation, information about the animation. For backward compatibility, when this field is set, the document field will also be set*/
	Animation *Animation `json:"animation,omitempty"`
	/*Optional. Message is an audio file, information about the file*/
	Audio *Audio `json:"audio,omitempty"`
	/*Optional. Message is a general file, information about the file*/
	Document *Document `json:"document,omitempty"`
	/*Optional. Message is a photo, available sizes of the photo*/
	Photo []PhotoSize `json:"photo,omitempty"`
	/*Optional. Message is a sticker, information about the sticker*/
	Sticker *Sticker `json:"sticker,omitempty"`
	/*Optional. Message is a video, information about the video*/
	Video *Video `json:"video,omitempty"`
	/*Optional. Message is a video note, information about the video message*/
	VideoNote *VideoNote `json:"video_note,omitempty"`
	/*Optional. Message is a voice message, information about the file*/
	Vocie *Voice `json:"voice,omitempty"`
	/*Optional. Caption for the animation, audio, document, photo, video or voice, 0-1024 characters*/
	Caption string `json:"caption,omitempty"`
	/*Optional. For messages with a caption, special entities like usernames, URLs, bot commands, etc. that appear in the caption*/
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`
	/*Optional. Message is a shared contact, information about the contact*/
	Contact *Contact `json:"contact,omitempty"`
	/*Optional. Message is a dice with random value*/
	Dice *Dice `json:"dice,omitempty"`
	/*Optional. Message is a game, information about the game.*/
	Game *Game `json:"game,omitempty"`
	/*Optional. Message is a native poll, information about the poll*/
	Poll *Poll `json:"poll,omitempty"`
	/*Optional. Message is a venue, information about the venue. For backward compatibility, when this field is set, the location field will also be set*/
	Venue *Venue `json:"venue,omitempty"`
	/*Optional. Message is a shared location, information about the location*/
	Location *Location `json:"location,omitempty"`
	/*Optional. New members that were added to the group or supergroup and information about them (the bot itself may be one of these members)*/
	NewChatMembers []User `json:"new_chat_members,omitempty"`
	/*Optional. A member was removed from the group, information about them (this member may be the bot itself)*/
	LeftChatMember *User `json:"left_chat_member,omitempty"`
	/*Optional. A chat title was changed to this value*/
	NewChatTitle string `json:"new_chat_title,omitempty"`
	/*Optional. A chat photo was change to this value*/
	NewChatPhoto []PhotoSize `json:"new_chat_photo,omitempty"`
	/*Optional. Service message: the chat photo was deleted*/
	DeleteChatPhoto bool `json:"delete_chat_photo,omitempty"`
	/*Optional. Service message: the group has been created*/
	GroupChatCreated bool `json:"group_chat_created,omitempty"`
	/*Optional. Service message: the supergroup has been created. This field can't be received in a message coming through updates, because bot can't be a member of a supergroup when it is created. It can only be found in reply_to_message if someone replies to a very first message in a directly created supergroup.*/
	SupergroupChatCreated bool `json:"supergroup_chat_created,omitempty"`
	/*Optional. Service message: the channel has been created. This field can't be received in a message coming through updates, because bot can't be a member of a channel when it is created. It can only be found in reply_to_message if someone replies to a very first message in a channel.*/
	ChannelChatCreated bool `json:"channel_chat_created,omitempty"`
	/*Optional. Service message: auto-delete timer settings changed in the chat*/
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	/*Optional. The group has been migrated to a supergroup with the specified identifier. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.*/
	MigrateToChatId int `json:"migrate_to_chat_id,omitempty"`
	/*Optional. The supergroup has been migrated from a group with the specified identifier. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.*/
	MigrateFromChatId int `json:"migrate_from_chat_id,omitempty"`
	/*Optional. Specified message was pinned. Note that the Message object in this field will not contain further reply_to_message fields even if it is itself a reply.*/
	PinnedMessage *Message `json:"pinned_message,omitempty"`
	/*Optional. Message is an invoice for a payment, information about the invoice*/
	Invoice *Invoice `json:"invoice,omitempty"`
	/*Optional. Message is a service message about a successful payment, information about the payment.*/
	SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`
	/*Optional. The domain name of the website on which the user has logged in.*/
	ConnectedWebsite string `json:"connected_website,omitempty"`
	/*Optional. Telegram Passport data*/
	PassportData *PassportData `json:"passport_data,omitempty"`
	/*Optional. Service message. A user in the chat triggered another user's proximity alert while sharing Live Location.*/
	ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"`
	/*Optional. Service message: voice chat scheduled*/
	VoiceChatScheduled *VoiceChatScheduled `json:"voice_chat_scheduled,omitempty"`
	/*Optional. Service message: voice chat started*/
	VoiceChatStarted *VoiceChatStarted `json:"voice_chat_started,omitempty"`
	/*Optional. Service message: voice chat ended*/
	VoiceChatEnded *VoiceChatEnded `json:"voice_chat_ended,omitempty"`
	/*Optional. Service message: new participants invited to a voice chat*/
	VoiceChatParticipantsInvited *VoiceChatParticipantsInvited `json:"voice_chat_participants_invited,omitempty"`
	/*Optional. Inline keyboard attached to the message. login_url buttons are represented as ordinary url buttons.*/
	ReplyMakrup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
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

func (rm *ForceReply) blah() {}

/*Contains information about why a request was unsuccessful.*/
type ResponseParameters struct {
	/*Optional. The group has been migrated to a supergroup with the specified identifier. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.*/
	MigrateToChatId int `json:"migrate_to_chat_id"`
	/*ptional. In case of exceeding flood control, the number of seconds left to wait before the request can be repeated*/
	RetryAfter int `json:"retry_after"`
}

/*Contains information about the current status of a webhook.*/
type WebhookInfo struct {
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

/*Not related to telegram bot api*/
type ChatUpdate struct {
	ChatId string
	Update *Update
}
