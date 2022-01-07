package configs

import (
	"testing"
)

type cfgTest struct {
	cfg            *BotConfigs
	expectedResult bool
}

var cfgs []cfgTest

func TestConfigCheck(t *testing.T) {
	initTheCfgs()
	for _, test := range cfgs {
		if test.cfg.Check() != test.expectedResult {
			t.Fail()
		}
	}
}

func initTheCfgs() {
	cfg1 := cfgTest{&BotConfigs{}, false}
	cfg2 := cfgTest{&BotConfigs{BotAPI: DefaultBotAPI}, false}
	cfg3 := cfgTest{&BotConfigs{BotAPI: DefaultBotAPI, APIKey: "sisduifhdsfsdf", Webhook: true}, false}
	cfg4 := cfgTest{&BotConfigs{BotAPI: DefaultBotAPI, APIKey: "sisduifhdsfsdf", Webhook: true, WebHookConfigs: &WebHookConfigs{}}, false}
	cfg5 := cfgTest{&BotConfigs{BotAPI: DefaultBotAPI, APIKey: "sisduifhdsfsdf", Webhook: false}, false}
	cfg6 := cfgTest{&BotConfigs{BotAPI: DefaultBotAPI, APIKey: "sisduifhdsfsdf", Webhook: false, UpdateConfigs: DefaultUpdateConfigs()}, true}
	cfgs = []cfgTest{cfg1, cfg2, cfg3, cfg4, cfg5, cfg6}
}
