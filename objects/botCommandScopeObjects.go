package objects

/*This object represents a bot command.*/
type BotCommand struct {
	/*Text of the command, 1-32 characters. Can contain only lowercase English letters, digits and underscores.*/
	Command string `json:"command"`
	/*Description of the command, 3-256 characters.*/
	Description string `json:"description"`
}

/*This object represents the scope to which bot commands are applied*/
type BotCommandScope struct {
	/*Scope type*/
	Type string `json:"type"`
}

/*Represents the default scope of bot commands. Default commands are used if no commands with a narrower scope are specified for the user.*/
type BotCommandScopeDefault struct {
	BotCommandScope
}

/*Represents the scope of bot commands, covering all private chats.*/
type BotCommandScopeAllPrivateChats struct {
	BotCommandScope
}

/*Represents the scope of bot commands, covering all group and supergroup chats.*/
type BotCommandScopeAllGroupChats struct {
	BotCommandScope
}

/*Represents the scope of bot commands, covering all group and supergroup chat administrators.*/
type BotCommandScopeAllChatAdministrators struct {
	BotCommandScope
}

/*Represents the scope of bot commands, covering a specific chat.*/
type BotCommandScopeChat struct {
	BotCommandScope
	/*Unique identifier for the target chat or username of the target supergroup (in the format @supergroupusername)*/
	ChatId string `json:"chat_id"`
}

/*Represents the scope of bot commands, covering all administrators of a specific group or supergroup chat.*/
type BotCommandScopeChatAdministrators struct {
	BotCommandScopeChat
}

/*Represents the scope of bot commands, covering a specific member of a group or supergroup chat.*/
type BotCommandScopeChatMember struct {
	BotCommandScopeChat
	/*Unique identifier of the target user*/
	UserId int `json:"user_id"`
}
