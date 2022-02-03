package telego

import (
	"os"

	objs "github.com/SakoDroid/telego/objects"
)

//ChatManager is a tool for managing chats via the bot.
type ChatManager struct {
	bot          *Bot
	chatIdInt    int
	chatIdString string
}

func (cm *ChatManager) fixThePerms(canSendMessages, canSendMediaMessages, canSendPolls, canSendOtherMessages, canAddWebPagePreviews, canChangeInfo, canInviteUsers, canPinMessages bool) objs.ChatPermissions {
	return objs.ChatPermissions{
		CanSendMessages:       canSendMessages,
		CanSendMediaMessages:  canSendMediaMessages,
		CanSendPolls:          canSendPolls,
		CanSendOtherMessages:  canSendOtherMessages,
		CanAddWebPagePreviews: canAddWebPagePreviews,
		CanChangeInfo:         canChangeInfo,
		CanInviteUsers:        canInviteUsers,
		CanPinMessages:        canPinMessages,
	}
}

/*BanMember bans a user in a group, a supergroup or a channel. In the case of supergroups and channels, the user will not be able to return to the chat on their own using invite links, etc., unless unbanned first. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns True on success.*/
func (cm *ChatManager) BanMember(userId, untilDate int, revokeMessages bool) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.BanChatMember(
		cm.chatIdInt, cm.chatIdString, userId, untilDate, revokeMessages,
	)
}

/*UnbanMember ubans a previously banned user in a supergroup or channel. The user will not return to the group or channel automatically, but will be able to join via link, etc. The bot must be an administrator for this to work. By default, this method guarantees that after the call the user is not a member of the chat, but will be able to join it. So if the user is a member of the chat they will also be removed from the chat. If you don't want this, use the parameter only_if_banned. Returns True on success.*/
func (cm *ChatManager) UnbanMember(userId int, onlyIfBanned bool) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.UnbanChatMember(
		cm.chatIdInt, cm.chatIdString, userId, onlyIfBanned,
	)
}

/*
Use this method to restrict a user in a supergroup. The bot must be an administrator in the supergroup for this to work and must have the appropriate administrator rights. Pass True for all permissions to lift restrictions from a user. Returns True on success.*/
func (cm *ChatManager) RestrictMember(userId int, untilDate int, canSendMessages, canSendMediaMessages, canSendPolls, canSendOtherMessages, canAddWebPagePreviews, canChangeInfo, canInviteUsers, canPinMessages bool) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.RestrictChatMember(
		cm.chatIdInt, cm.chatIdString, userId, cm.fixThePerms(
			canSendMessages, canSendMediaMessages, canSendPolls, canSendOtherMessages, canAddWebPagePreviews, canChangeInfo, canInviteUsers, canPinMessages,
		), untilDate,
	)
}

/*PromoteChatMember promotes or demote a user in a supergroup or a channel. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Pass False for all boolean parameters to demote a user. Returns True on success.*/
func (cm *ChatManager) PromoteChatMember(userId int, isAnonymous, canManageChat, canPostmessages, canEditMessages, canDeleteMessages, canManageVoiceChats, canRestrictMembers, canPromoteMembers, canChangeInfo, canInviteUsers, canPinMessages bool) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.PromoteChatMember(
		cm.chatIdInt, cm.chatIdString, userId, isAnonymous, canManageChat,
		canPostmessages, canEditMessages, canDeleteMessages, canManageVoiceChats,
		canRestrictMembers, canPromoteMembers, canChangeInfo, canInviteUsers, canPinMessages,
	)
}

/*SetCustomTitle sets a custom title for an administrator in a supergroup promoted by the bot. Returns True on success.*/
func (cm *ChatManager) SetCustomTitle(userId int, customTitle string) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.SetChatAdministratorCustomTitle(
		cm.chatIdInt, cm.chatIdString, userId, customTitle,
	)
}

/*BanChatSender bans a channel chat in a supergroup or a channel. Until the chat is unbanned, the owner of the banned chat won't be able to send messages on behalf of any of their channels. The bot must be an administrator in the supergroup or channel for this to work and must have the appropriate administrator rights. Returns True on success.*/
func (cm *ChatManager) BanChatSender(senderChatId int) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.BanOrUnbanChatSenderChat(
		cm.chatIdInt, cm.chatIdString, senderChatId, true,
	)
}

/*UnbanChatSender unbans a previously banned channel chat in a supergroup or channel. The bot must be an administrator for this to work and must have the appropriate administrator rights. Returns True on success.*/
func (cm *ChatManager) UnbanChatSender(senderChatId int) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.BanOrUnbanChatSenderChat(
		cm.chatIdInt, cm.chatIdString, senderChatId, false,
	)
}

/*SetGeneralPermissions sets default chat permissions for all members. The bot must be an administrator in the group or a supergroup for this to work and must have the can_restrict_members administrator rights. Returns True on success.*/
func (cm *ChatManager) SetGeneralPermissions(canSendMessages, canSendMediaMessages, canSendPolls, canSendOtherMessages, canAddWebPagePreviews, canChangeInfo, canInviteUsers, canPinMessages bool) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.SetChatPermissions(
		cm.chatIdInt, cm.chatIdString, cm.fixThePerms(
			canSendMessages, canSendMediaMessages, canSendPolls, canSendOtherMessages, canAddWebPagePreviews, canChangeInfo, canInviteUsers, canPinMessages,
		),
	)
}

/*ExportInviteLink generates a new primary invite link for a chat; any previously generated primary link is revoked. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns the new invite link as String on success.

Note: Each administrator in a chat generates their own invite links. Bots can't use invite links generated by other administrators. If you want your bot to work with invite links, it will need to generate its own link using this method or by calling the getChat method. If your bot needs to generate a new primary invite link replacing its previous one, use this method again.*/
func (cm *ChatManager) ExportInviteLink() (*objs.StringResult, error) {
	return cm.bot.apiInterface.ExportChatInviteLink(
		cm.chatIdInt, cm.chatIdString,
	)
}

/*CreateInviteLink creates an additional invite link for a chat. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. The link can be revoked using the method RevokeInviteLink. Returns the new invite link as ChatInviteLink object.*/
func (cm *ChatManager) CreateInviteLink(name string, expireDate, memberLimit int, createsJoinRequest bool) (*objs.ChatInviteLinkResult, error) {
	return cm.bot.apiInterface.CreateChatInviteLink(
		cm.chatIdInt, cm.chatIdString, name, expireDate, memberLimit, createsJoinRequest,
	)
}

/*EditInviteLink edits a non-primary invite link created by the bot. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns the edited invite link as a ChatInviteLink object.*/
func (cm *ChatManager) EditInviteLink(inviteLink, name string, expireDate, memberLimit int, createsJoinRequest bool) (*objs.ChatInviteLinkResult, error) {
	return cm.bot.apiInterface.EditChatInviteLink(
		cm.chatIdInt, cm.chatIdString, inviteLink, name, expireDate, memberLimit, createsJoinRequest,
	)
}

/*RevokeInviteLink revokes an invite link created by the bot. If the primary link is revoked, a new link is automatically generated. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns the revoked invite link as ChatInviteLink object.*/
func (cm *ChatManager) RevokeInviteLink(inviteLink string) (*objs.ChatInviteLinkResult, error) {
	return cm.bot.apiInterface.RevokeChatInviteLink(
		cm.chatIdInt, cm.chatIdString, inviteLink,
	)
}

/*ApproveJoinRequest approves a chat join request. The bot must be an administrator in the chat for this to work and must have the can_invite_users administrator right. Returns True on success.*/
func (cm *ChatManager) ApproveJoinRequest(userId int) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.ApproveChatJoinRequest(
		cm.chatIdInt, cm.chatIdString, userId,
	)
}

/*DeclineJoinRequest can be used to decline a chat join request. The bot must be an administrator in the chat for this to work and must have the can_invite_users administrator right. Returns True on success.*/
func (cm *ChatManager) DeclineJoinRequest(userId int) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.DeclineChatJoinRequest(
		cm.chatIdInt, cm.chatIdString, userId,
	)
}

/*SetPhoto can be used to set a new profile photo for the chat. Photos can't be changed for private chats. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns True on success.*/
func (cm *ChatManager) SetPhoto(photoFile *os.File) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.SetChatPhoto(
		cm.chatIdInt, cm.chatIdString, photoFile,
	)
}

/*DeletePhoto can be used to delete a chat photo. Photos can't be changed for private chats. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns True on success.*/
func (cm *ChatManager) DeletePhoto() (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.DeleteChatPhoto(
		cm.chatIdInt, cm.chatIdString,
	)
}

/*SetTitle changes the title of a chat. Titles can't be changed for private chats. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns True on success.*/
func (cm *ChatManager) SetTitle(title string) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.SetChatTitle(
		cm.chatIdInt, cm.chatIdString, title,
	)
}

/*SetDescription changes the description of a group, a supergroup or a channel. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Returns True on success.*/
func (cm *ChatManager) SetDescription(description string) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.SetChatDescription(
		cm.chatIdInt, cm.chatIdString, description,
	)
}

/*PinMessage adds a message to the list of pinned messages in a chat. If the chat is not a private chat, the bot must be an administrator in the chat for this to work and must have the 'can_pin_messages' administrator right in a supergroup or 'can_edit_messages' administrator right in a channel. Returns True on success.*/
func (cm *ChatManager) PinMessage(messageId int, disableNotif bool) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.PinChatMessage(
		cm.chatIdInt, cm.chatIdString, messageId, disableNotif,
	)
}

/*UnpinMessage removes a message from the list of pinned messages in a chat. If the chat is not a private chat, the bot must be an administrator in the chat for this to work and must have the 'can_pin_messages' administrator right in a supergroup or 'can_edit_messages' administrator right in a channel. Returns True on success.*/
func (cm *ChatManager) UnpinMessage(messageId int) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.UnpinChatMessage(
		cm.chatIdInt, cm.chatIdString, messageId,
	)
}

/*UnpinAllMessages clears the list of pinned messages in a chat. If the chat is not a private chat, the bot must be an administrator in the chat for this to work and must have the 'can_pin_messages' administrator right in a supergroup or 'can_edit_messages' administrator right in a channel. Returns True on success.*/
func (cm *ChatManager) UnpinAllMessages() (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.UnpinAllChatMessages(
		cm.chatIdInt, cm.chatIdString,
	)
}

/*Leave can be used for your bot to leave a group, supergroup or channel. Returns True on success.*/
func (cm *ChatManager) Leave() (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.LeaveChat(
		cm.chatIdInt, cm.chatIdString,
	)
}

/*GetChatInfo gets up to date information about the chat (current name of the user for one-on-one conversations, current username of a user, group or channel, etc.). Returns a Chat object on success.*/
func (cm *ChatManager) GetChatInfo() (*objs.ChatResult, error) {
	return cm.bot.apiInterface.GetChat(
		cm.chatIdInt, cm.chatIdString,
	)
}

/*GetAdmins gets a list of administrators in a chat. On success, returns an Array of ChatMember objects that contains information about all chat administrators except other bots. If the chat is a group or a supergroup and no administrators were appointed, only the creator will be returned.*/
func (cm *ChatManager) GetAdmins() (*objs.ChatAdministratorsResult, error) {
	return cm.bot.apiInterface.GetChatAdministrators(
		cm.chatIdInt, cm.chatIdString,
	)
}

/*GetMembersCount gets the number of members in a chat. Returns Int on success.*/
func (cm *ChatManager) GetMembersCount() (*objs.IntResult, error) {
	return cm.bot.apiInterface.GetChatMemberCount(
		cm.chatIdInt,
		cm.chatIdString,
	)
}

/*GetMember gets information about a member of a chat. Returns a json serialized object of the member in string form on success.*/
func (cm *ChatManager) GetMember(userid int) (string, error) {
	res, err := cm.bot.apiInterface.GetChatMember(
		cm.chatIdInt, cm.chatIdString, userid,
	)
	if err != nil {
		return "", err
	}
	return string(res.Result), nil
}

/*SetStickerSet sets a new group sticker set for a supergroup. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Use the field can_set_sticker_set optionally returned in "GetChatInfo" to check if the bot can use this method. Returns True on success.*/
func (cm *ChatManager) SetStickerSet(stickerSetName string) (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.SetChatStickerSet(
		cm.chatIdInt, cm.chatIdString, stickerSetName,
	)
}

/*DeleteStickerSet deletes a group sticker set from a supergroup. The bot must be an administrator in the chat for this to work and must have the appropriate administrator rights. Use the field can_set_sticker_set optionally returned in "GetChatInfo" to check if the bot can use this method. Returns True on success.*/
func (cm *ChatManager) DeleteStickerSet() (*objs.LogicalResult, error) {
	return cm.bot.apiInterface.DeleteChatStickerSet(
		cm.chatIdInt, cm.chatIdString,
	)
}
