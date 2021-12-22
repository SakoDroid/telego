package objects

import "encoding/json"

type Chat struct {
	/*Unique identifier for this chat. This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it. But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.*/
	Id int `json:"id"`
	/*Type of chat, can be either “private”, “group”, “supergroup” or “channel”*/
	Type string `json:"type"`
	/*Optional. Title, for supergroups, channels and group chats*/
	Title string `json:"title,omitempty"`
	/*Optional. Username, for private chats, supergroups and channels if available*/
	Username string `json:"username,omitempty"`
	/*Optional. First name of the other party in a private chat*/
	FirstName string `json:"first_name,omitempty"`
	/*Optional. Last name of the other party in a private chat*/
	LastName string `json:"last_name,omitempty"`
	/*Optional. Chat photo. */
	Photo *ChatPhoto `json:"photo,omitempty"`
	/*Optional. Bio of the other party in a private chat.*/
	Bio string `json:"bio,omitempty"`
	/*Optional. True, if privacy settings of the other party in the private chat allows to use tg://user?id=<user_id> links only in chats with the user.*/
	HasPrivateForwards bool `json:"has_private_forwards,omitempty"`
	/*Optional. Description, for groups, supergroups and channel chats.*/
	Description string `json:"description,omitempty"`
	/*Optional. Primary invite link, for groups, supergroups and channel chats.*/
	InviteLink string `json:"invite_link,omitempty"`
	/*Optional. The most recent pinned message (by sending date).*/
	PinnedMessage *Message `json:"pinned_message,omitempty"`
	/*Optional. Default chat member permissions, for groups and supergroups.*/
	Permissions *ChatPermissions `json:"permissions,omitempty"`
	/*Optional. For supergroups, the minimum allowed delay between consecutive messages sent by each unpriviledged user; in seconds.*/
	SlowModeDelay int `json:"slow_mode_delay,omitempty"`
	/*Optional. The time after which all messages sent to the chat will be automatically deleted; in seconds.*/
	MessageAutoDeletTime int `json:"message_auto_delete_time,omitempty"`
	/*Optional. True, if messages from the chat can't be forwarded to other chats*/
	HasProtectedContent bool `json:"has_protected_content,omitempty"`
	/*Optional. For supergroups, name of group sticker set.*/
	StickerSetName string `json:"sticker_set_name,omitempty"`
	/*Optional. True, if the bot can change the group sticker set.*/
	CanSetStickerSet bool `json:"can_set_sticker_set,omitempty"`
	/*Optional. Unique identifier for the linked chat, i.e. the discussion group identifier for a channel and vice versa; for supergroups and channel chats. This identifier may be greater than 32 bits and some programming languages may have difficulty/silent defects in interpreting it. But it is smaller than 52 bits, so a signed 64 bit integer or double-precision float type are safe for storing this identifier.*/
	LinkedChatId int `json:"linked_chat_id,omitempty"`
	/*Optional. For supergroups, the location to which the supergroup is connected. */
	Location *ChatLocation `json:"location,omitempty"`
}

type ChatPhoto struct {
	/*File identifier of small (160x160) chat photo. This file_id can be used only for photo download and only for as long as the photo is not changed.*/
	SmallFileId string `json:"small_file_id"`
	/*Unique file identifier of small (160x160) chat photo, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.*/
	SmallFileUniqueId string `json:"small_file_unique_id"`
	/*File identifier of big (640x640) chat photo. This file_id can be used only for photo download and only for as long as the photo is not changed.*/
	BigFileId string `json:"big_file_id"`
	/*Unique file identifier of big (640x640) chat photo, which is supposed to be the same over time and for different bots. Can't be used to download or reuse the file.*/
	BigFileUniqueId string `json:"big_file_unique_id"`
}

type ChatInviteLink struct {
	/*The invite link. If the link was created by another chat administrator, then the second part of the link will be replaced with “…”.*/
	InviteLink string `json:"invite_link"`
	/*Creator of the link*/
	Creator *User `json:"user"`
	/*True, if users joining the chat via the link need to be approved by chat administrators*/
	CreatesJoinRequest bool `json:"creates_join_request"`
	/*True, if the link is primary*/
	IsPrimary bool `json:"is_primary"`
	/*True, if the link is revoked*/
	IsRevoked bool `json:"is_revoked"`
	/*Optional. Invite link name*/
	Name string `json:"name,omitempty"`
	/*Optional. Point in time (Unix timestamp) when the link will expire or has been expired*/
	ExpireDate int `json:"expire_date,omitempty"`
	/*Optional. Maximum number of users that can be members of the chat simultaneously after joining the chat via this invite link; 1-99999*/
	MemberLimit int `json:"member_limit,omitempty"`
	/*Optional. Number of pending join requests created using this link*/
	PendingJoinRequestCount int `json:"pending_join_request_count,omitempty"`
}

/*This object represents changes in the status of a chat member.*/
type ChatMemberUpdated struct {
	/*Chat the user belongs to*/
	Chat *Chat `json:"chat"`
	/*Performer of the action, which resulted in the change*/
	From *User `json:"from"`
	/*Date the change was done in Unix time*/
	Date int `json:"date"`
	/*Previous information about the chat member*/
	OldChatMember json.RawMessage `json:"old_chat_member"`
	/*New information about the chat member*/
	NewChatMember json.RawMessage `json:"new_chat_member"`
	/*Optional. Chat invite link, which was used by the user to join the chat; for joining by invite link events only.*/
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

/*Represents a join request sent to a chat.*/
type ChatJoinRequest struct {
	/*Chat to which the request was sent*/
	Chat *Chat `json:"chat"`
	/*User that sent the join request*/
	From *User `json:"from"`
	/*Date the request was sent in Unix time*/
	Date int `json:"date"`
	/*Optional. Bio of the user.*/
	Bio string `json:"bio,omitempty"`
	/*Optional. Chat invite link that was used by the user to send the join request*/
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

/*Describes actions that a non-administrator user is allowed to take in a chat.*/
type ChatPermissions struct {
	/*Optional. True, if the user is allowed to send text messages, contacts, locations and venues*/
	CanSendMessages bool `json:"can_send_messages,omitempty"`
	/*Optional. True, if the user is allowed to send audios, documents, photos, videos, video notes and voice notes, implies can_send_messages*/
	CanSendMediaMessages bool `json:"can_send_media_messages,omitempty"`
	/*Optional. True, if the user is allowed to send polls, implies can_send_messages*/
	CanSendPolls bool `json:"can_send_polls,omitempty"`
	/*Optional. True, if the user is allowed to send animations, games, stickers and use inline bots, implies can_send_media_messages*/
	CanSendOtherMessages bool `json:"can_send_other_messages,omitempty"`
	/*Optional. True, if the user is allowed to add web page previews to their messages, implies can_send_media_messages*/
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	/*Optional. True, if the user is allowed to change the chat title, photo and other settings. Ignored in public supergroups*/
	CanChangeInfo bool `json:"can_change_info,omitempty"`
	/*Optional. True, if the user is allowed to invite new users to the chat*/
	CanInviteUsers bool `json:"can_invite_users,omitempty"`
	/*Optional. True, if the user is allowed to pin messages. Ignored in public supergroups*/
	CanPinMessages bool `json:"can_pin_messages,omitempty"`
}

/*Represents a location to which a chat is connected.*/
type ChatLocation struct {
	/*The location to which the supergroup is connected. Can't be a live location.*/
	Location *Location `json:"location"`
	/*Location address; 1-64 characters, as defined by the chat owner*/
	Address string `json:"address"`
}
