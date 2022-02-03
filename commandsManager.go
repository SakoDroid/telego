package telego

import (
	"errors"

	objs "github.com/SakoDroid/telego/objects"
)

//CommandsManager is a tool for managing bot commands
type CommandsManager struct {
	bot      *Bot
	commands []objs.BotCommand
	scope    objs.BotCommandScope
}

/*AddCommand adds a new command to the commands list.

At most 100 commands can be added.*/
func (cm *CommandsManager) AddCommand(command, description string) error {
	if len(cm.commands) >= 100 {
		return errors.New("command list is full")
	}
	cm.commands = append(cm.commands, objs.BotCommand{Command: command, Description: description})
	return nil
}

/*SetScope sets the scope of this command manager.

Scope can have these values : "defaut","all_group_chats","all_private_chats","all_chat_administrators","chat","chat_administrator","chat_member". If scope is not valid error is returned. */
func (cm *CommandsManager) SetScope(scope string, chatId []byte, userId int) error {
	switch scope {
	case "default":
		cm.scope = &objs.BotCommandScopeDefault{}
	case "all_group_chats":
		cm.scope = &objs.BotCommandScopeAllGroupChats{}
	case "all_private_chats":
		cm.scope = &objs.BotCommandScopeAllPrivateChats{}
	case "all_chat_administrators":
		cm.scope = &objs.BotCommandScopeAllChatAdministrators{}
	case "chat":
		cm.scope = &objs.BotCommandScopeChat{ChatId: chatId}
	case "chat_member":
		cm.scope = &objs.BotCommandScopeChatMember{BotCommandScopeChat: objs.BotCommandScopeChat{ChatId: chatId}, UserId: userId}
	case "chat_administrator":
		cm.scope = &objs.BotCommandScopeChatAdministrators{BotCommandScopeChat: objs.BotCommandScopeChat{ChatId: chatId}}
	default:
		return errors.New(scope + " value is note allowed.")
	}
	return nil
}

/*SetCommands calls the realted method on the api server and sets the added commansd with their specified scopes.

"languageCode" is a two-letter ISO 639-1 language code. If empty, commands will be applied to all users from the given scope, for whose language there are no dedicated commands

-------------------------

Official telegram doc :


Use this method to change the list of the bot's commands. See https://core.telegram.org/bots#commands for more details about bot commands. Returns True on success.*/
func (cm *CommandsManager) SetCommands(languageCode string) (*objs.LogicalResult, error) {
	if cm.scope == nil {
		return nil, errors.New("scope is not set. Use `SetScope` method")
	}
	return cm.bot.apiInterface.SetMyCommands(cm.commands, cm.scope, languageCode)
}

/*DeleteCommands can be used to delete the list of the bot's commands for the given scope and user language. After deletion, higher level commands will be shown to affected users. Returns True on success.*/
func (cm *CommandsManager) DeleteCommands(languageCode string) (*objs.LogicalResult, error) {
	if cm.scope == nil {
		return nil, errors.New("scope is not set. Use `SetScope` method")
	}
	return cm.bot.apiInterface.DeleteMyCommands(cm.scope, languageCode)
}

//GetCommands returns the commands of this bot.
func (cm *CommandsManager) GetCommands(languageCode string) ([]objs.BotCommand, error) {
	if cm.scope == nil {
		return nil, errors.New("scope is not set. Use `SetScope` method")
	}
	res, err := cm.bot.apiInterface.GetMyCommands(cm.scope, languageCode)
	if err != nil {
		return nil, err
	}
	cm.commands = res.Result
	return cm.commands, nil
}
