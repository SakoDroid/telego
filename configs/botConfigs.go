package configs

import "time"

const DefaultBotAPI = "https://api.telegram.org/bot"
const DefaultLogFile = "./bot-logs.log"

type BotConfigs struct {
	/*This is the bot api server. If you dont have a local bot api server, user "configs.DefaultBotAPI" for this field.*/
	BotAPI string
	/*The API key for your bot. You can get the api key (token) from botfather*/
	APIKey string
	/*The settings related to getting updates from the api server. This field shoud only be populated when Webhook field is false, otherwise it is ignored.*/
	UpdateConfigs *UpdateConfigs
	/*This field idicates if webhook should be used for receiving updates or not.
	Recommend : false*/
	Webhook bool
	/*All the logs related to bot will be written in this file. You can use configs.DefaultLogFile for default value*/
	LogFileAddress string
}

type UpdateConfigs struct {
	/*Limits the number of updates to be retrieved. Values between 1-100 are accepted. Defaults to 100.*/
	Limit int
	/*Timeout in seconds for long polling. Defaults to 0, i.e. usual short polling. Should be positive, short polling should be used for testing purposes only.*/
	Timeout int
	/*List of the update types you want your bot to receive. For example, specify [“message”, “edited_channel_post”, “callback_query”] to only receive updates of these types. See Update for a complete list of available update types. Specify an empty list to receive all update types except chat_member (default). If not specified, the previous setting will be used.
	Please note that this parameter doesnt affect updates created before the call to the getUpdates, so unwanted updates may be received for a short period of time.*/
	AllowedUpdates []string
	/*This field indicates the frequency to call getUpdates method. Default is one second*/
	UpdateFrequency time.Duration
}

func DefaultUpdateConfigs() *UpdateConfigs {
	return &UpdateConfigs{Limit: 100, Timeout: 0, UpdateFrequency: time.Duration(300 * time.Millisecond), AllowedUpdates: nil}
}
