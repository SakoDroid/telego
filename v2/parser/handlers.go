package parser

import (
	"regexp"
	"strings"

	objs "github.com/SakoDroid/telego/v2/objects"
)

var handlers = HandlerTree{}
var callbackHandlers = threadSafeMap[string, *callbackHandler]{internal: make(map[string]*callbackHandler)}
var userSharedHandlers = threadSafeMap[int, *chatRequestHandler]{internal: make(map[int]*chatRequestHandler)}
var chatSharedHandlers = threadSafeMap[int, *chatRequestHandler]{internal: make(map[int]*chatRequestHandler)}

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
	callbackHandlers.Add(data, &hl)
}

func AddUserSharedHandler(requestId int, handler func(*objs.Update)) {
	userSharedHandlers.Add(
		requestId,
		&chatRequestHandler{
			requestId: requestId,
			function:  &handler,
		},
	)
}

func AddChatSharedHandler(requestId int, handler func(*objs.Update)) {
	chatSharedHandlers.Add(
		requestId,
		&chatRequestHandler{
			requestId: requestId,
			function:  &handler,
		},
	)
}

func checkHandlers(up *objs.Update) bool {
	if up.CallbackQuery != nil {
		return checkCallbackHanlders(up)
	}

	if up.Message!=nil && up.Message.UserShared != nil {
		return checkUserSharedHandlers(up)
	}

	if up.Message!=nil && up.Message.ChatShared != nil {
		return checkChatSharedHandlers(up)
	}

	return checkTextMsgHandlers(up)
}

func checkCallbackHanlders(up *objs.Update) bool {
	hdl, ok := callbackHandlers.Load(up.CallbackQuery.Data)
	if ok && hdl != nil {
		go (*hdl.function)(up)
		return true
	}
	return false
}

func checkUserSharedHandlers(up *objs.Update) bool {
	hdl, ok := userSharedHandlers.LoadAndDelete(up.Message.UserShared.RequestId)
	if ok && hdl != nil && hdl.function != nil {
		go (*hdl.function)(up)
		return true
	}
	return false
}

func checkChatSharedHandlers(up *objs.Update) bool {
	hdl, ok := chatSharedHandlers.Load(up.Message.ChatShared.RequestId)
	if ok && hdl != nil && hdl.function != nil {
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
