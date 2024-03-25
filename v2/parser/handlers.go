package parser

import (
	"regexp"
	"strings"

	objs "github.com/SakoDroid/telego/v2/objects"
)

// var handlers = handlerTree{}
// var callbackHandlers = threadSafeMap[string, *callbackHandler]{internal: make(map[string]*callbackHandler)}
// var userSharedHandlers = threadSafeMap[int, *chatRequestHandler]{internal: make(map[int]*chatRequestHandler)}
// var chatSharedHandlers = threadSafeMap[int, *chatRequestHandler]{internal: make(map[int]*chatRequestHandler)}

type handler struct {
	regex    *regexp.Regexp      //The compiled regex.
	chatType string              //The ChatType this handler will act on
	function *func(*objs.Update) //The function to be executed
}

type callbackHandler struct {
	callbackData string
	function     *func(*objs.Update)
}

func (up *UpdateParser) AddHandler(patern string, handlerFunc func(*objs.Update), chatType ...string) error {
	hl := handler{chatType: strings.Join(chatType, ","), function: &handlerFunc}
	rgxp, err := regexp.Compile(patern)
	if err != nil {
		return err
	}
	hl.regex = rgxp
	up.handlers.AddHandler(&hl)
	return nil
}

func (up *UpdateParser) AddCallbackHandler(data string, handlerFun func(*objs.Update)) {
	hl := callbackHandler{callbackData: data, function: &handlerFun}
	up.callbackHandlers.Add(data, &hl)
}

func (up *UpdateParser) AddUserSharedHandler(requestId int, handler func(*objs.Update)) {
	up.userSharedHandlers.Add(
		requestId,
		&chatRequestHandler{
			requestId: requestId,
			function:  &handler,
		},
	)
}

func (up *UpdateParser) AddChatSharedHandler(requestId int, handler func(*objs.Update)) {
	up.chatSharedHandlers.Add(
		requestId,
		&chatRequestHandler{
			requestId: requestId,
			function:  &handler,
		},
	)
}

func (up *UpdateParser) checkHandlers(update *objs.Update) bool {
	if update.CallbackQuery != nil {
		return up.checkCallbackHanlders(update)
	}

	if update.Message.UserShared != nil {
		return up.checkUserSharedHandlers(update)
	}

	if update.Message.ChatShared != nil {
		return up.checkChatSharedHandlers(update)
	}

	return up.checkTextMsgHandlers(update)
}

func (up *UpdateParser) checkCallbackHanlders(update *objs.Update) bool {
	hdl, ok := up.callbackHandlers.Load(update.CallbackQuery.Data)
	if ok && hdl != nil {
		go (*hdl.function)(update)
		return true
	}
	return false
}

func (up *UpdateParser) checkUserSharedHandlers(update *objs.Update) bool {
	hdl, ok := up.userSharedHandlers.LoadAndDelete(update.Message.UserShared.RequestId)
	if ok && hdl != nil && hdl.function != nil {
		go (*hdl.function)(update)
		return true
	}
	return false
}

func (up *UpdateParser) checkChatSharedHandlers(update *objs.Update) bool {
	hdl, ok := up.chatSharedHandlers.Load(update.Message.ChatShared.RequestId)
	if ok && hdl != nil && hdl.function != nil {
		go (*hdl.function)(update)
		return true
	}
	return false
}

func (up *UpdateParser) checkTextMsgHandlers(update *objs.Update) bool {
	if update.Message != nil && (update.Message.Text != "" || update.Message.Caption != "") {
		hndl := up.handlers.GetHandler(update.Message)
		if hndl != nil {
			go (*hndl.function)(update)
			return true
		}
	}
	return false
}
