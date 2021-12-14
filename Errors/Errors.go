package Errors

import (
	"strconv"

	objs "github.com/SakoDroid/telebot/objects"
)

type MethodNotSentError struct {
	Method, Reason string
	FailureResult  *objs.FailureResult
}

func (mnse *MethodNotSentError) Error() string {
	out := "Unable to send " + mnse.Method + ". " + mnse.Reason
	if mnse.FailureResult != nil {
		out += "\nError code : " + strconv.Itoa(mnse.FailureResult.ErrorCode) + ", Description : " + mnse.FailureResult.Description
	}
	return out
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

type RequiredArgumentError struct {
	ArgName, MethodName string
}

func (ram *RequiredArgumentError) Error() string {
	return "Required argument \"" + ram.ArgName + " missing or has wrong value in " + ram.MethodName + " method."
}

type ChatIdProblem struct {
}

func (cip *ChatIdProblem) Error() string {
	return "Cannot have both chatIdInt and chatIdString at the same time. Only one of them is allowed."
}

type MediaGroupFullError struct {
}

func (mgfe *MediaGroupFullError) Error() string {
	return "the media group is full. (10 files)"
}
