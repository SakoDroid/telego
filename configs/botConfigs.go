package configs

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

// BlockedUser is a struct used for storing a blocked user informations.
type BlockedUser struct {
	UserID   int    `json:"user_id"`
	UserName string `json:"username"`
}

// DefaultBotAPI is the telegrams default bot api server.
const DefaultBotAPI = "https://api.telegram.org/bot"

// DefaultLogFile is a default file for saving the bot logs in it.
const DefaultLogFile = "STDOUT"

// BotConfigs is a struct holding the bots configs.
type BotConfigs struct {
	/*This is the bot api server. If you dont have a local bot api server, use "configs.DefaultBotAPI" for this field.*/
	BotAPI string `json:"bot_api"`
	/*The API key for your bot. You can get the api key (token) from botfather*/
	APIKey string `json:"api_key"`
	/*The settings related to getting updates from the api server. This field shoud only be populated when Webhook field is false, otherwise it is ignored.*/
	UpdateConfigs *UpdateConfigs `json:"update_configs,omitempty"`
	/*This field idicates if webhook should be used for receiving updates or not.
	Recommend : false*/
	Webhook bool `json:"webhook"`
	/*This field represents the configs related to web hook.*/
	WebHookConfigs *WebHookConfigs `json:"webhook_configs,omitempty"`
	/*All the logs related to bot will be written in this file. You can use configs.DefaultLogFile for default value*/
	LogFileAddress string `json:"log_file"`
	//BlockedUsers is a list of blocked users.
	BlockedUsers []BlockedUser `json:"blocked_users"`

	ConfigName string `json:"config_name"`
}

// Check checks the bot configs for any problem.
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

// StartCfgUpdateRoutine starts a routine which updates the configs every second.
func (bc *BotConfigs) StartCfgUpdateRoutine() {
	for {
		err := LoadInto(bc)
		if err != nil {
			println("Error in \"StartCfgUpdateRoutine\" function.", err.Error())
			break
		}
		time.Sleep(time.Second)
	}
}

// Load loads the configs from the config file (configs.json) and returns the BotConfigs pointer.
func Load(configName string) (*BotConfigs, error) {
	fl, err := os.Open(configName)
	defer fl.Close()
	if err != nil {
		return nil, err
	}
	st, err := fl.Stat()
	if err != nil {
		return nil, err
	}
	data := make([]byte, st.Size())
	_, err = fl.Read(data)
	if err != nil {
		return nil, err
	}
	bc := &BotConfigs{}
	err = json.Unmarshal(data, bc)
	return bc, err
}

// LoadInto works the same way as "Load" but it won't return the config, instead it loads the config into the given object.
func LoadInto(bc *BotConfigs) error {
	fl, err := os.Open(bc.ConfigName)
	defer fl.Close()
	if err != nil {
		return err
	}
	st, err := fl.Stat()
	if err != nil {
		return err
	}
	data := make([]byte, st.Size())
	_, err = fl.Read(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, bc)
	return err
}

// Dump saves the given BotConfigs struct in a json format in the config file (configs.json).
func Dump(bc *BotConfigs) error {
	fl, err := os.OpenFile(bc.ConfigName, os.O_CREATE|os.O_WRONLY, 0666)
	defer fl.Close()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(bc, "", " ")
	if err != nil {
		return err
	}
	_, err = fl.Write(data)
	return err
}

// WebHookConfigs contains the configs necessary for webhook.
type WebHookConfigs struct {
	/*The web hook url.*/
	URL string `json:"url"`
	/*The port that webhook server will run on. Telegram api only suppotrs 80,443,88,8443. 8443 is recommended. Pass 0 for default https port (443)*/
	Port int `json:"port"`
	/*The address of the public key certificate file.*/
	KeyFile string `json:"keyfile"`
	/*The address of the certificate file.*/
	CertFile string `json:"certfile"`
	/*Is your certificate self signed?*/
	SelfSigned bool
	/*The fixed IP address which will be used to send webhook requests instead of the IP address resolved through DNS*/
	IP string `json:"ip,omitempty"`
	/*Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery, 1-100. Defaults to 40. Use lower values to limit the load on your bot's server, and higher values to increase your bot's throughput.*/
	MaxConnections int `json:"max_connections"`
	/*List of the update types you want your bot to receive. For example, specify [“message”, “edited_channel_post”, “callback_query”] to only receive updates of these types. See Update for a complete list of available update types. Specify an empty list to receive all update types except chat_member (default). If not specified, the previous setting will be used.
	Please note that this parameter doesnt affect updates created before the call to the getUpdates, so unwanted updates may be received for a short period of time.*/
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
	/*Pass True to drop all pending updates*/
	DropPendingUpdates bool `json:"drop_pending_reqs"`
	/*A secret token to be sent in a header “X-Telegram-Bot-Api-Secret-Token” in every webhook request, 1-256 characters. Only characters A-Z, a-z, 0-9, _ and - are allowed. The header is useful to ensure that the request comes from a webhook set by you.*/
	SecretToken string `json:"secret_token,omitempty"`
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

// UpdateConfigs contains the necessary configs for receiving updates.
type UpdateConfigs struct {
	/*Limits the number of updates to be retrieved. Values between 1-100 are accepted. Defaults to 100.*/
	Limit int `json:"limit"`
	/*Timeout in seconds for long polling. Defaults to 0, i.e. usual short polling. Should be positive, short polling should be used for testing purposes only.*/
	Timeout int `json:"timeout"`
	/*List of the update types you want your bot to receive. For example, specify [“message”, “edited_channel_post”, “callback_query”] to only receive updates of these types. See Update for a complete list of available update types. Specify an empty list to receive all update types except chat_member (default). If not specified, the previous setting will be used.
	Please note that this parameter doesnt affect updates created before the call to the getUpdates, so unwanted updates may be received for a short period of time.*/
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
	/*This field indicates the frequency to call getUpdates method. Default is one second*/
	UpdateFrequency time.Duration `json:"update_freq"`
}

// DefaultUpdateConfigs returns a default update configs.
func DefaultUpdateConfigs() *UpdateConfigs {
	return &UpdateConfigs{Limit: 100, Timeout: 0, UpdateFrequency: time.Duration(300 * time.Millisecond), AllowedUpdates: nil}
}

// Default returns default setting for the bot.
func Default(apiKey string) *BotConfigs {
	return &BotConfigs{
		BotAPI:         DefaultBotAPI,
		APIKey:         apiKey,
		UpdateConfigs:  DefaultUpdateConfigs(),
		Webhook:        false,
		LogFileAddress: DefaultLogFile,
		ConfigName:     "configs.json",
	}
}

// Default returns default setting for the bot.
func DefaultConfigName(apiKey, configName string) *BotConfigs {
	return &BotConfigs{
		BotAPI:         DefaultBotAPI,
		APIKey:         apiKey,
		UpdateConfigs:  DefaultUpdateConfigs(),
		Webhook:        false,
		LogFileAddress: DefaultLogFile,
		ConfigName:     configName,
	}
}
