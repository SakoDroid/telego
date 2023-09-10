package configs

import (
	"os"
	"reflect"
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

func TestLoadAndDump(t *testing.T) {
	bc1 := Default("123hUHASDa66aDTDAFshdASDKabda6dg982edua")
	err := Dump(bc1)
	if err != nil {
		t.Error(err)
	}
	bc2, err := Load("configs.json")
	if err != nil {
		t.Error(err)
	}
	if reflect.DeepEqual(bc1, bc2) {
		t.Log("bc1 and bc2 are equal")
	} else {
		t.Error("bc1 and bc2 are not equal")
	}
	_ = os.Remove("configs.json")
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
