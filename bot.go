package telebot

import (
	tba "github.com/SakoDroid/telebot/TBA"
	cfg "github.com/SakoDroid/telebot/configs"
	logger "github.com/SakoDroid/telebot/logger"
	objs "github.com/SakoDroid/telebot/objects"
)

type Bot struct {
	botCfg        *cfg.BotConfigs
	apiInterface  *tba.BotAPIInterface
	updateChannel *chan *objs.Update
}

/*Starts the bot. If the bot has already been started it returns an error.*/
func (bot *Bot) Run() error {
	logger.InitTheLogger(bot.botCfg)
	return bot.apiInterface.StartUpdateRoutine()
}

/*Returns the channel which new updates received from api server are pushed into.*/
func (bot *Bot) GetUpdateChannel() *chan *objs.Update {
	return bot.updateChannel
}

/*Send a text message to a chat (not channel, use SendMessageToChannel method for sending messages to channles) and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.*/
func (bot *Bot) SendMessage(chatId int, text, parseMode string, replyTo int, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendMessage(chatId, "", text, parseMode, nil, false, silent, false, replyTo, nil)
}

/*Send a text message to a channel and returns the sent message on success
If you want to ignore "parseMode" pass empty string. To ignore replyTo pass 0.*/
func (bot *Bot) SendMesssageToChannel(chatId, text, parseMode string, replyTo int, silent bool) (*objs.SendMethodsResult, error) {
	return bot.apiInterface.SendMessage(0, chatId, text, parseMode, nil, false, silent, false, replyTo, nil)
}

/*Returns a MessageForwarder which has several methods for forwarding a message*/
func (bot *Bot) ForwardMessage(messageId int, disableNotif bool) *MessageForwarder {
	return &MessageForwarder{bot: bot, messageId: messageId, disableNotif: disableNotif}
}

/*Returns a MessageCopier which has several methods for copying a message*/
func (bot *Bot) CopyMessage(messageId int, disableNotif bool) *MessageCopier {
	return &MessageCopier{bot: bot, messageId: messageId, disableNotif: disableNotif}
}

/*Returns a PhotoSender which has several methods for sending a photo. This method is only used for sending a photo to all types of chat except channels. To send a photo to a channel use "SendPhotoToChannel" method.*/
func (bot *Bot) SendPhoto(chatId, replyTo int, caption, parseMode string) *PhotoSender {
	return &PhotoSender{bot: bot, chatIdInt: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
}

func (bot *Bot) SendPhotoToChannel(chatId string, replyTo int, caption, parseMode string) *PhotoSender {
	return &PhotoSender{bot: bot, chatIdInt: 0, chatidString: chatId, replyTo: replyTo, caption: caption, parseMode: parseMode}
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
