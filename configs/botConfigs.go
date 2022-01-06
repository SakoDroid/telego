package configs

import (
	"strings"
	"time"
)

const DefaultBotAPI = "https://api.telegram.org/bot"
const DefaultLogFile = "./bot-logs.log"

type BotConfigs struct {
	/*This is the bot api server. If you dont have a local bot api server, use "configs.DefaultBotAPI" for this field.*/
	BotAPI string
	/*The API key for your bot. You can get the api key (token) from botfather*/
	APIKey string
	/*The settings related to getting updates from the api server. This field shoud only be populated when Webhook field is false, otherwise it is ignored.*/
	UpdateConfigs *UpdateConfigs
	/*This field idicates if webhook should be used for receiving updates or not.
	Recommend : false*/
	Webhook bool
	/*This field represents the configs related to web hook.*/
	WebHookConfigs *WebHookConfigs
	/*All the logs related to bot will be written in this file. You can use configs.DefaultLogFile for default value*/
	LogFileAddress string
}

func (bc *BotConfigs) Check() bool {
	if bc.BotAPI == "" {
		return false
	}
	if bc.APIKey == "" {
		return false
	}
	if bc.Webhook {
		if bc.WebHookConfigs != nil {
			return bc.WebHookConfigs.check(bc.APIKey)
		}
		return false
	} else {
		return bc.UpdateConfigs != nil
	}
}

type WebHookConfigs struct {
	/*The web hook url.*/
	URL string
	/*The port that webhook server will run on. Telegram api only suppotrs 80,443,88,8443. 443 is recommended. Pass 0 for default https port (443)*/
	Port int
	/*The address of the public key certificate file.*/
	KeyFile string
	/*The address of the certificate file.*/
	CertFile string
	/*The fixed IP address which will be used to send webhook requests instead of the IP address resolved through DNS*/
	IP string
	/*Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery, 1-100. Defaults to 40. Use lower values to limit the load on your bot's server, and higher values to increase your bot's throughput.*/
	MaxConnections int
	/*List of the update types you want your bot to receive. For example, specify [“message”, “edited_channel_post”, “callback_query”] to only receive updates of these types. See Update for a complete list of available update types. Specify an empty list to receive all update types except chat_member (default). If not specified, the previous setting will be used.
	Please note that this parameter doesnt affect updates created before the call to the getUpdates, so unwanted updates may be received for a short period of time.*/
	AllowedUpdates []string
	/*Pass True to drop all pending updates*/
	DropPendingUpdates bool
}

func (whc *WebHookConfigs) check(apiKey string) bool {
	if whc.URL == "" {
		return false
	}
	if whc.KeyFile == "" {
		return false
	}
	if whc.CertFile == "" {
		return false
	}
	if whc.Port == 0 {
		whc.Port = 443
	}
	if !strings.HasSuffix(whc.URL, apiKey) {
		if !strings.HasSuffix(whc.URL, "/") {
			whc.URL += "/"
		}
		whc.URL += apiKey
	}
	return true
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
