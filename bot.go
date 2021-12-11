package telebot

import (
	tba "github.com/SakoDroid/telebot/TBA"
	cfg "github.com/SakoDroid/telebot/configs"
	objs "github.com/SakoDroid/telebot/objects"
)

type Bot struct {
	botCfg        *cfg.BotConfigs
	apiInterface  *tba.BotAPIInterface
	updateChannel *chan *objs.Update
}

/*Send a text message to a chat (not channel, use SendMessageToChannel method for sending messages to channles) and returns the sent message on success*/
func (bot *Bot) SendMessage(chatId int, text, parseMode string, replyTo int, silent bool) (*objs.Message, error) {
	return bot.apiInterface.SendMessage(chatId, "", text, parseMode, nil, false, silent, false, replyTo, nil)
}

/*Starts the bot. If the bot has already been started it returns an error.*/
func (bot *Bot) Run() error {
	return bot.apiInterface.StartUpdateRoutine()
}

/*Stops the bot*/
func (bot *Bot) Stop() {
	bot.apiInterface.StopUpdateRoutine()
}

/*Returns and advanced version which gives more customized functions to iteract with the bot*/
func (bot *Bot) AdvancedMode() *AdvancedBot {
	return &AdvancedBot{Bot: bot}
}

/*Return a new bot instance with the specified configs*/
func NewBot(cfg *cfg.BotConfigs) (*Bot, error) {
	api, err := tba.CreateInterface(cfg)
	if err != nil {
		return nil, err
	}
	return &Bot{botCfg: cfg, apiInterface: api, updateChannel: api.GetUpdateChannel()}, nil
}
