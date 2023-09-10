package objects

import "encoding/json"

/*This object represents a bot command.*/
type BotCommand struct {
	/*Text of the command, 1-32 characters. Can contain only lowercase English letters, digits and underscores.*/
	Command string `json:"command"`
	/*Description of the command, 3-256 characters.*/
	Description string `json:"description"`
}

type BotCommandScope interface {
	FixTheType()
}

/*Represents the default scope of bot commands. Default commands are used if no commands with a narrower scope are specified for the user.*/
type BotCommandScopeDefault struct {
	/*Scope type*/
	Type string `json:"type"`
}

func (bc *BotCommandScopeDefault) FixTheType() {
	bc.Type = "default"
}

/*Represents the scope of bot commands, covering all private chats.*/
type BotCommandScopeAllPrivateChats struct {
	BotCommandScopeDefault
}

func (bc *BotCommandScopeAllPrivateChats) FixTheType() {
	bc.Type = "all_private_chats"
}

/*Represents the scope of bot commands, covering all group and supergroup chats.*/
type BotCommandScopeAllGroupChats struct {
	BotCommandScopeDefault
}

func (bc *BotCommandScopeAllGroupChats) FixTheType() {
	bc.Type = "all_group_chats"
}

/*Represents the scope of bot commands, covering all group and supergroup chat administrators.*/
type BotCommandScopeAllChatAdministrators struct {
	BotCommandScopeDefault
}

func (bc *BotCommandScopeAllChatAdministrators) FixTheType() {
	bc.Type = "all_chat_administrators"
}

/*Represents the scope of bot commands, covering a specific chat.*/
type BotCommandScopeChat struct {
	BotCommandScopeDefault
	/*Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)*/
	ChatId json.RawMessage `json:"chat_id"`
}

func (bc *BotCommandScopeChat) FixTheType() {
	bc.Type = "chat"
}

/*Represents the scope of bot commands, covering all administrators of a specific group or supergroup chat.*/
type BotCommandScopeChatAdministrators struct {
	BotCommandScopeChat
}

func (bc *BotCommandScopeChatAdministrators) FixTheType() {
	bc.Type = "chat_administrators"
}

/*Represents the scope of bot commands, covering a specific member of a group or supergroup chat.*/
type BotCommandScopeChatMember struct {
	BotCommandScopeChat
	/*Unique identifier of the target user*/
	UserId int `json:"user_id"`
}

func (bc *BotCommandScopeChatMember) FixTheType() {
	bc.Type = "chat_member"
}
