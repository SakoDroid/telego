package objects

type ChatMember interface {
	blah()
}

/*Represents a chat member that has no additional privileges or restrictions*/
type ChatMemberMember struct {
	/*The member's status in the chat.*/
	Status string `json:"status"`
	/*Information about the user*/
	User User `json:"user"`
}

func (*ChatMemberMember) blah() {}

/*Represents a chat member that has some additional privileges.*/
type ChatMemberOwner struct {
	ChatMemberMember
	/*True, if the bot is allowed to edit administrator privileges of that user*/
	CanBeEdited bool `json:"can_be_edited"`
	/*True, if the user's presence in the chat is hidden*/
	IsAnonymous bool `json:"is_anonymous"`
	/*True, if the administrator can access the chat event log, chat statistics, message statistics in channels, see channel members, see anonymous administrators in supergroups and ignore slow mode. Implied by any other administrator privilege*/
	CanManageChat bool `json:"can_manage_chat"`
	/*True, if the administrator can delete messages of other users*/
	CanDeleteMessages bool `json:"can_delete_messages"`
	/*True, if the administrator can manage voice chats*/
	CanManageVoiceChats bool `json:"can_manage_voice_chats"`
	/*True, if the administrator can restrict, ban or unban chat members*/
	CanRestrictMembers bool `json:"can_restrict_members"`
	/*True, if the administrator can add new administrators with a subset of their own privileges or demote administrators that he has promoted, directly or indirectly (promoted by administrators that were appointed by the user)*/
	CanPromoteMembers bool `json:"can_promote_members"`
	/*True, if the user is allowed to change the chat title, photo and other settings*/
	CanChangeInfo bool `json:"can_change_info"`
	/*True, if the user is allowed to invite new users to the chat*/
	CanInviteUsers bool `json:"can_invite_users"`
	/*Optional. True, if the administrator can post in the channel; channels only*/
	CanPostMessages bool `json:"can_post_messages,omitempty"`
	/*Optional. True, if the administrator can edit messages of other users and can pin messages; channels only*/
	CanEditMessages bool `json:"can_edit_messages,omitempty"`
	/*Optional. True, if the user is allowed to pin messages; groups and supergroups only*/
	CanPinMessages bool `json:"can_pin_messages,omitempty"`
	/*Optional. Custom title for this user*/
	CustomTitle string `json:"custom_title,omitempty"`
}

func (*ChatMemberOwner) blah() {}

type ChatMemberRestricted struct {
	ChatMemberMember
	/*True, if the user is a member of the chat at the moment of the request*/
	IsMember bool `json:"is_member"`
	/*True, if the user is allowed to change the chat title, photo and other settings*/
	CanChangeInfo bool `json:"can_change_info"`
	/*True, if the user is allowed to invite new users to the chat*/
	CanInviteUsers bool `json:"can_invite_users"`
	/*True, if the user is allowed to pin messages*/
	CanPinMessages bool `json:"can_pin_messages"`
	/*True, if the user is allowed to send text messages, contacts, locations and venues*/
	CanSendMessages bool `json:"can_send_messages"`
	/*True, if the user is allowed to send audios, documents, photos, videos, video notes and voice notes*/
	CanSendMediaMessages bool `json:"can_send_media_messages"`
	/*True, if the user is allowed to send polls*/
	CanSendPolls bool `json:"can_send_polls"`
	/*True, if the user is allowed to send animations, games, stickers and use inline bots*/
	CanSendOtherMessages bool `json:"can_send_other_messages"`
	/*True, if the user is allowed to add web page previews to their messages*/
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	/*Date when restrictions will be lifted for this user; unix time. If 0, then the user is banned forever*/
	UntilDate int `json:"until_date"`
}

func (*ChatMemberRestricted) blah() {}

/*Represents a chat member that isn't currently a member of the chat, but may join it themselves.*/
type ChatMemberLeft struct {
	ChatMemberMember
}

func (*ChatMemberLeft) blah() {}

/*Represents a chat member that was banned in the chat and can't return to the chat or view chat messages.*/
type ChatMemberBanned struct {
	ChatMemberMember
	/*Date when restrictions will be lifted for this user; unix time. If 0, then the user is banned forever*/
	UntilDate int `json:"until_date"`
}

func (*ChatMemberBanned) blah() {}
