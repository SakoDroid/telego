package parser

import (
	"regexp"
	"strings"

	objs "github.com/SakoDroid/telego/objects"
)

var handlers = HandlerTree{}
var callbackHandlers = make(map[string]*callbackHandler)

type handler struct {
	regex    *regexp.Regexp      //The compiled regex.
	chatType string              //The ChatType this handler will act on
	function *func(*objs.Update) //The function to be executed
}

type callbackHandler struct {
	callbackData string
	function     *func(*objs.Update)
}

func AddHandler(patern string, handlerFunc func(*objs.Update), chatType ...string) error {
	hl := handler{chatType: strings.Join(chatType, ","), function: &handlerFunc}
	rgxp, err := regexp.Compile(patern)
	if err != nil {
		return err
	}
	hl.regex = rgxp
	handlers.AddHandler(&hl)
	return nil
}

func AddCallbackHandler(data string, handlerFun func(*objs.Update)) {
	hl := callbackHandler{callbackData: data, function: &handlerFun}
	callbackHandlers[data] = &hl
}

func checkHandlers(up *objs.Update) bool {
	if up.CallbackQuery != nil {
		return checkCallbackHanlders(up)
	} else {
		return checkTextMsgHandlers(up)
	}
}

func checkCallbackHanlders(up *objs.Update) bool {
	hdl := callbackHandlers[up.CallbackQuery.Data]
	if hdl != nil {
		go (*hdl.function)(up)
		return true
	}
	return false
}

func checkTextMsgHandlers(up *objs.Update) bool {
	if up.Message != nil && up.Message.Text != "" {
		hndl := handlers.GetHandler(up.Message)
		if hndl != nil {
			go (*hndl.function)(up)
			return true
		}
	}
	return false
}
