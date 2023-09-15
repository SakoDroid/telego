package telego

import (
	objs "github.com/SakoDroid/telego/v2/objects"
)

// BotManager is a tool for manaing personal info of the bot.
type BotManager struct {
	bot *Bot
}

/*
GetMe returns the received informations about the bot from api server.

---------------------

Official telegarm doc :

A simple method for testing your bot's authentication token. Requires no parameters. Returns basic information about the bot in form of a User object.
*/
func (bm *BotManager) GetMe() (*objs.Result[*objs.User], error) {
	return bm.bot.apiInterface.GetMe()
}

/*
SetDescription sets the description of the bot for the specified langauge. Description is shown in the chat with the bot if the chat is empty.

Arguments :

1. description : New bot description; 0-512 characters. Pass an empty string to remove the dedicated description for the given language.

2. languageCode : A two-letter ISO 639-1 language code. If empty, the description will be applied to all users for whose language there is no dedicated description.
*/
func (bm *BotManager) SetDescription(description, languageCode string) (*objs.Result[bool], error) {
	return bm.bot.apiInterface.SetMyDescription(description, languageCode)
}

/*
SetShortDescription sets the short description of the bot for the specified langauge. Short description is shown on the bot's profile page and is sent together with the link when users share the bot.

Arguments :

1. shortDescription : New short description for the bot; 0-120 characters. Pass an empty string to remove the dedicated short description for the given language.

2. languageCode : A two-letter ISO 639-1 language code. If empty, the description will be applied to all users for whose language there is no dedicated description.
*/
func (bm *BotManager) SetShortDescription(shortDescription, languageCode string) (*objs.Result[bool], error) {
	return bm.bot.apiInterface.SetMyShortDescription(shortDescription, languageCode)
}

/*
GetDescription returns description of the bot based on the specified language.
*/
func (bm *BotManager) GetDescription(languageCode string) (*objs.Result[*objs.BotDescription], error) {
	return bm.bot.apiInterface.GetMyDescription(languageCode)
}

/*
GetShortDescription returns short description of the bot based on the specified language.
*/
func (bm *BotManager) GetShortDescription(languageCode string) (*objs.Result[*objs.BotDescription], error) {
	return bm.bot.apiInterface.GetMyDescription(languageCode)
}

// GetName returns the bot's name according to the language code
func (bm *BotManager) GetName(languageCode string) (*objs.Result[*objs.BotName], error) {
	return bm.bot.apiInterface.GetMyName(languageCode)
}

/*
SetName sets bot's name in the specified language

Arguments :

1. name : bot's new name

2. languageCode : A two-letter ISO 639-1 language code. If empty, the name will be shown to all users for whose language there is no dedicated name.
*/
func (bm *BotManager) SetName(name, languageCode string) (*objs.Result[bool], error) {
	return bm.bot.apiInterface.SetMyName(name, languageCode)
}
