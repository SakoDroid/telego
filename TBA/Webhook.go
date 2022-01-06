package TBA

import (
	cfg "github.com/SakoDroid/telego/configs"
	objs "github.com/SakoDroid/telego/objects"
)

var webhookCfg *cfg.WebHookConfigs
var interfaceUpdateChannel *chan *objs.Update
var chatUpdateChannel *chan *objs.ChatUpdate

func StartWebHook(cfg *cfg.WebHookConfigs, iuc *chan *objs.Update, cuc *chan *objs.ChatUpdate) error {
	webhookCfg = cfg
	interfaceUpdateChannel = iuc
	chatUpdateChannel = cuc
	return startTheServer()
}

func startTheServer() error {

}
