package parser

import (
	"regexp"
	"testing"

	"github.com/SakoDroid/telego/objects"
)

var tree = &HandlerTree{}
var testTable []handlerTest

type handlerTest struct {
	msg           *objects.Message
	expectedRegex string
}

func TestTree(t *testing.T) {
	initTheHandlers()
	initTheTable()
	for _, test := range testTable {
		hndl := tree.GetHandler(test.msg)
		if hndl == nil {
			t.Error("nil handler, msg details : ", test.msg.Text, ",", test.msg.Chat.Type)
			continue
		}
		if hndl.regex == nil {
			t.Error("nil regex, msg details : ", test.msg.Text, ",", test.msg.Chat.Type)
			continue
		}
		if hndl.regex.String() != test.expectedRegex {
			t.Error("Wrong regex, msg details : ", test.msg.Text, ",", test.msg.Chat.Type, ". Regex :", hndl.regex.String())
		}
	}
	nilTest1 := handlerTest{msg: &objects.Message{Text: "nope", Chat: &objects.Chat{Type: "private"}}}
	hndl := tree.GetHandler(nilTest1.msg)
	if hndl != nil {
		t.Error("not nil handler, msg details : ", nilTest1.msg.Text, ",", nilTest1.msg.Chat.Type)
	}

	t.Log("done")
}

func initTheHandlers() {
	handler1 := &handler{regex: regexp.MustCompile("hi"), chatType: "all"}
	handler2 := &handler{regex: regexp.MustCompile("hi guys"), chatType: "private"}
	handler3 := &handler{regex: regexp.MustCompile("start"), chatType: "all"}
	handler4 := &handler{regex: regexp.MustCompile("start again"), chatType: "private"}
	handler5 := &handler{regex: regexp.MustCompile("start bot"), chatType: "all"}
	handler6 := &handler{regex: regexp.MustCompile("hi everyone"), chatType: "private,group"}
	tree.AddHandler(handler1)
	tree.AddHandler(handler2)
	tree.AddHandler(handler3)
	tree.AddHandler(handler4)
	tree.AddHandler(handler5)
	tree.AddHandler(handler6)
}

func initTheTable() {
	test1 := handlerTest{msg: &objects.Message{Text: "hi guys", Chat: &objects.Chat{Type: "private"}}, expectedRegex: "hi guys"}
	test2 := handlerTest{msg: &objects.Message{Text: "hi guys", Chat: &objects.Chat{Type: "group"}}, expectedRegex: "hi"}
	test3 := handlerTest{msg: &objects.Message{Text: "start", Chat: &objects.Chat{Type: "group"}}, expectedRegex: "start"}
	test4 := handlerTest{msg: &objects.Message{Text: "start bot", Chat: &objects.Chat{Type: "supergroup"}}, expectedRegex: "start bot"}
	test5 := handlerTest{msg: &objects.Message{Text: "hi everyone", Chat: &objects.Chat{Type: "group"}}, expectedRegex: "hi everyone"}
	test6 := handlerTest{msg: &objects.Message{Text: "hi everyone", Chat: &objects.Chat{Type: "channel"}}, expectedRegex: "hi"}
	test7 := handlerTest{msg: &objects.Message{Text: "start again", Chat: &objects.Chat{Type: "group"}}, expectedRegex: "start"}
	test8 := handlerTest{msg: &objects.Message{Text: "start again", Chat: &objects.Chat{Type: "private"}}, expectedRegex: "start again"}
	testTable = []handlerTest{test1, test2, test3, test4, test5, test6, test7, test8}
}
