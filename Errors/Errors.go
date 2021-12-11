package Errors

import "strconv"

type MethodNotSentError struct {
	Method, Reason string
}

func (mnse *MethodNotSentError) Error() string {
	return "Unable to send " + mnse.Method + ". " + mnse.Reason
}

type BotInterfaceAlreadyCreated struct {
}

func (biac *BotInterfaceAlreadyCreated) Error() string {
	return "only one bot interface can be created"
}

type UpdateRoutineAlreadyStarted struct {
}

func (uras *UpdateRoutineAlreadyStarted) Error() string {
	return "StartUpdateRoutine has already been called."
}

type UpdateNotOk struct {
	Offset int
}

func (uno *UpdateNotOk) Error() string {
	return "getUpdates returned \"ok\" : false. Offset used in request : " + strconv.Itoa(uno.Offset)
}
