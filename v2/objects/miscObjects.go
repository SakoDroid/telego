package objects

type MenuButton struct {
	Type   string      `json:"type"`
	Text   string      `json:"text,omitempty"`
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

type ChatAdministratorRights struct {
	/*True, if the user's presence in the chat is hidden*/
	IsAnonymous bool `json:"is_anonymous"`
	/*True, if the administrator can access the chat event log, chat statistics, message statistics in channels, see channel members, see anonymous administrators in supergroups and ignore slow mode. Implied by any other administrator privilege*/
	CanManageChat bool `json:"can_manage_chat"`
	/*True, if the administrator can delete messages of other users*/
	CanDeleteMessages bool `json:"can_delete_messages"`
	/*True, if the administrator can manage voice chats*/
	CanManageVideoChats bool `json:"can_manage_video_chats"`
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
	/*Optional. True, if the administrator can post stories in the channel; channels only*/
	CanPostStories bool `json:"can_post_stories"`
	/*Optional. True, if the administrator can edit stories posted by other users; channels only*/
	CanEditStories bool `json:"can_edit_stories"`
	/*Optional. True, if the administrator can delete stories posted by other users; channels only*/
	CanDeleteStories bool `json:"can_delete_stories"`
	/*Optional. True, if the user is allowed to create, rename, close, and reopen forum topics; supergroups only*/
	CanManageTopics bool `json:"can_manage_topics"`
}

type UserShared struct {
	/*Identifier of the request*/
	RequestId int `json:"request_id"`
	/*dentifier of the shared user. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier. The bot may not have access to the user and could be unable to use this identifier, unless the user is already known to the bot by some other means.*/
	UserId int64 `json:"user_id"`
}

type ChatShared struct {
	/*Identifier of the request*/
	RequestId int `json:"request_id"`
	/*Identifier of the shared chat. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a 64-bit integer or double-precision float type are safe for storing this identifier. The bot may not have access to the chat and could be unable to use this identifier, unless the chat is already known to the bot by some other means.*/
	ChatId int64 `json:"chat_id"`
}

type BotName struct {
	//The bot's name
	Name string `json:"name"`
}

type BotDescription struct {
	/*The bot's description*/
	Description string `json:"description"`
}

type BotShortDescription struct {
	/*The bot's short description*/
	ShortDescription string `json:"short_description"`
}
