package Parser

import (
	"regexp"

	objs "github.com/SakoDroid/telego/objects"
)

var handlers = make([]*handler, 0)
var callbackHandlers = make([]*callbackHandler, 0)

type handler struct {
	regex    *regexp.Regexp      //The compiled regex.
	chatType []string            //The ChatType this handler will act on
	function *func(*objs.Update) //The function to be executed
}

type callbackHandler struct {
	callbackData string
	function     *func(*objs.Update)
}

func AddHandler(patern string, handlerFunc func(*objs.Update), chatType ...string) error {
	hl := handler{chatType: chatType, function: &handlerFunc}
	rgxp, err := regexp.Compile(patern)
	if err != nil {
		return err
	}
	hl.regex = rgxp
	handlers = append(handlers, &hl)
	return nil
}

func AddCallbackHandler(data string, handlerFun func(*objs.Update)) {
	hl := callbackHandler{callbackData: data, function: &handlerFun}
	callbackHandlers = append(callbackHandlers, &hl)
}

func checkHandlers(up *objs.Update) bool {
	if up.CallbackQuery != nil {
		return checkCallbackHanlders(up)
	} else {
		return checkTextMsgHandlers(up)
	}
}

func checkCallbackHanlders(up *objs.Update) bool {
	for _, val := range callbackHandlers {
		if val.callbackData == up.CallbackQuery.Data {
			go (*val.function)(up)
			return true
		}
	}
	return false
}

func checkTextMsgHandlers(up *objs.Update) bool {
	if up.Message != nil && up.Message.Text != "" {
		for _, hndl := range handlers {
			if checkHandler(up.Message, hndl) {
				go (*hndl.function)(up)
				return true
			}
		}
	}
	return false
}

func checkHandler(msg *objs.Message, hndl *handler) bool {
	if hndl.regex.Match([]byte(msg.Text)) {
		for _, val := range hndl.chatType {
			if val == "all" || val == msg.Chat.Type {
				return true
			}
		}
		return false
	}
	return false
}
