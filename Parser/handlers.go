package Parser

import (
	"regexp"

	objs "github.com/SakoDroid/telego/objects"
)

var handlers = make([]*handler, 0)

type handler struct {
	regex    *regexp.Regexp      //The compiled regex.
	chatType string              //The ChatType this handler will act on
	function *func(*objs.Update) //The function to be executed
}

func AddHandler(patern, chatType string, handlerFunc func(*objs.Update)) error {
	hl := handler{chatType: chatType, function: &handlerFunc}
	rgxp, err := regexp.Compile(patern)
	if err != nil {
		return err
	}
	hl.regex = rgxp
	handlers = append(handlers, &hl)
	return nil
}

func checkHandlers(up *objs.Update) bool {
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
		return hndl.chatType == "all" || hndl.chatType == msg.Chat.Type
	}
	return false
}
